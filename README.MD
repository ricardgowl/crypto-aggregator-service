# Crypto Aggregator Service

Microservicio en Go que agrega precios de criptomonedas desde múltiples proveedores de forma concurrente y los expone a través de una API REST.

## Tabla de contenidos

- [Arquitectura](#arquitectura)
- [Decisiones de diseño](#decisiones-de-diseño)
- [Estructura del proyecto](#estructura-del-proyecto)
- [Requisitos](#requisitos)
- [Ejecución](#ejecución)
- [Endpoints](#endpoints)
- [Configuración](#configuración)
- [Testing](#testing)

## Arquitectura

El servicio sigue una arquitectura hexagonal (ports & adapters) con separación clara entre capas:

```
┌─────────────────────────────────────────────────┐
│                   cmd/main.go                   │
│              (Bootstrap & Wiring)               │
├─────────────────────────────────────────────────┤
│              Adapters (HTTP API)                │
│   health_controller · poller_controller         │
│         router · server_response                │
├─────────────────────────────────────────────────┤
│              Services (Lógica)                  │
│           Poller (polling concurrente)          │
├─────────────────────────────────────────────────┤
│             Repositories (Datos)                │
│    LayoutStore (estado in-memory)               │
│    CryptoClient (interfaz de proveedores)       │
│    Bitso · Mock (implementaciones)              │
├─────────────────────────────────────────────────┤
│              Models (Dominio)                   │
│        Component · Model · Money · Ticker       │
└─────────────────────────────────────────────────┘
```

### Flujo de datos

1. El servicio arranca cargando un layout estático con los componentes (BTC, ETH, XRP).
2. El `Poller` inicia un loop en background que cada 5 segundos:
   - Lee el layout actual del `LayoutStore`.
   - Lanza goroutines concurrentes (una por componente) para consultar al proveedor asignado.
   - Actualiza el store con los precios obtenidos.
3. El endpoint `GET /fetch` devuelve el estado actual del layout con los precios más recientes.

## Decisiones de diseño

### Interfaz `CryptoClient`

Se definió la interfaz `CryptoClient` con dos métodos (`GetPrice`, `Name`) para desacoplar la lógica de negocio de los proveedores concretos. Esto permite:
- Agregar nuevos proveedores (Coinbase, CoinMarketCap) sin modificar el servicio.
- Sustituir proveedores reales por mocks en los tests.
- Configurar qué proveedor usa cada componente vía `config.yaml`.

### `LayoutStore` con `sync.RWMutex`

El estado se mantiene en memoria con un `LayoutStore` protegido por `RWMutex`:
- Lecturas concurrentes (`RLock`) desde el endpoint HTTP.
- Escrituras exclusivas (`Lock`) desde las goroutines del poller.
- `GetLayout()` devuelve una copia defensiva para evitar data races.

### Polling concurrente con `WaitGroup`

El `Poller.refresh()` lanza una goroutine por componente y sincroniza con `sync.WaitGroup`. Si un proveedor falla, se loguea el error y el modelo se actualiza con precio cero, sin afectar a los demás componentes.

### Fallback de proveedor

Si el vendor configurado para un componente no existe en el mapa de clientes, el servicio cae automáticamente al cliente `mock`, garantizando que siempre haya una respuesta.

### Fallback USD en Bitso

Bitso no ofrece pares USD para todas las criptomonedas. Cuando la consulta del book USD falla, se estima el precio USD dividiendo el precio MXN entre 20 como aproximación.

### Inyección de dependencias

Todas las dependencias se inyectan vía constructores (`NewPoller`, `NewHTTPServer`, `NewLayoutStore`), facilitando el testing y evitando estado global.

## Estructura del proyecto

```
crypto-aggregator-service/
├── cmd/
│   └── main.go                          # Punto de entrada
├── config/
│   ├── config.go                        # Carga de configuración (YAML + env vars)
│   ├── layout_loader.go                 # Layout estático predefinido
│   ├── config_test.go                   # Tests de configuración
│   └── logger.go                        # Logger estructurado (Zap + ECS)
├── internal/
│   ├── adapters/
│   │   ├── mock.go                      # Cliente mock con precios simulados
│   │   ├── mock_test.go                 # Tests del mock client
│   │   ├── httpapi/
│   │   │   ├── router.go               # Chi router con middleware
│   │   │   ├── router_test.go           # Tests del router
│   │   │   ├── health_controller.go     # Endpoints /health/live, /health/ready, /metrics
│   │   │   ├── health_controller_test.go
│   │   │   ├── poller_controller.go     # Endpoint /fetch
│   │   │   ├── poller_controller_test.go
│   │   │   ├── server_response.go       # Helpers para respuestas JSON
│   │   │   └── server_response_test.go
│   │   └── webclients/
│   │       ├── http_client.go           # Factory de http.Client con connection pooling
│   │       └── http_client_test.go
│   ├── models/
│   │   ├── crypto.go                    # Model, Money, Ticker
│   │   ├── crypto_test.go
│   │   ├── errors.go                    # ProvidersError, ErrNoProviders
│   │   ├── errors_test.go
│   │   ├── layout.go                    # Component, Layout
│   │   └── layout_test.go
│   ├── repositories/
│   │   ├── client.go                    # Interfaz CryptoClient
│   │   ├── bitso.go                     # Implementación Bitso API
│   │   ├── bitso_test.go               # Tests con httptest + redirectTransport
│   │   ├── layout_store.go             # Store in-memory thread-safe
│   │   └── layout_store_test.go        # Tests incluyendo concurrencia
│   └── services/
│       ├── poller.go                    # Servicio de polling concurrente
│       └── poller_test.go              # Tests con stub clients
├── resources/
│   └── config.yaml                      # Configuración de la aplicación
├── Dockerfile                           # Build multi-stage (Go builder + Alpine)
├── go.mod
└── go.sum
```

## Requisitos

- Go 1.25+
- Docker (opcional, para ejecución en contenedor)

## Ejecución

### Local

```bash
go run ./cmd/main.go
```

El servidor escucha en `http://localhost:3000`.

### Docker

```bash
docker build -t crypto-aggregator .
docker run --rm -p 3000:3000 crypto-aggregator
```

### Variables de entorno

La configuración se puede sobreescribir con variables de entorno. El formato convierte guiones bajos en niveles de anidación:

```bash
SERVER_PORT=8080 go run ./cmd/main.go
```

## Endpoints

| Método | Ruta             | Descripción                              |
|--------|------------------|------------------------------------------|
| GET    | `/fetch`         | Layout actualizado con precios actuales  |
| GET    | `/health/live`   | Liveness probe (Kubernetes)              |
| GET    | `/health/ready`  | Readiness probe (Kubernetes)             |
| GET    | `/metrics`       | Métricas Prometheus                      |

### Ejemplo de respuesta `GET /fetch`

```json
[
  {
    "id": 1,
    "component": "crypto_btc",
    "model": {
      "date": "2025-02-26T17:00:00Z",
      "name": "BTC",
      "ticker_symbol": "BTC",
      "price": {
        "usd": 50000.00,
        "mxn": 850000.00
      }
    }
  },
  {
    "id": 2,
    "component": "crypto_eth",
    "model": {
      "date": "2025-02-26T17:00:00Z",
      "name": "ETH",
      "ticker_symbol": "ETH",
      "price": {
        "usd": 3000.00,
        "mxn": 51000.00
      }
    }
  }
]
```

## Configuración

El archivo `resources/config.yaml` define el layout y el proveedor para cada componente:

```yaml
server:
  port: 3000
  refresh_interval: 10

app:
  layout:
    - id: 1
      component: crypto_btc
      vendor: bitso
    - id: 2
      component: crypto_eth
      vendor: bitso
    - id: 3
      component: crypto_xrp
      vendor: bitso
```

Proveedores disponibles: `bitso` (API real), `mock` (precios simulados).

## Testing

### Ejecutar todos los tests

```bash
go test ./... -v
```

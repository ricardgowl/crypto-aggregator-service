## Microservicio Agregador de Criptomonedas
### Objetivo
Desarrollar un microservicio que exponga un endpoint REST. Este endpoint deberá:
- Cargar una configuración (layout) que contiene una lista de componentes.
- Consultar de forma concurrente múltiples proveedores de datos para obtener los valores actuales de diversas criptomonedas.
- Integrar estos valores en el modelo correspondiente de cada componente.
- Responder con un JSON actualizado y estructurado.
---
### Layout de Configuración
El microservicio iniciará con un layout predefinido que indica qué componentes se deben actualizar. Por ejemplo:
``` json
[
{
"id": 1,
"component": "crypto_btc",
"model": {}
},
{
"id": 2,
"component": "crypto_eth",
"model": {}
},
{
"id": 3,
"component": "crypto_xrp",
"model": {}
}
]

```
---
### Modelo de Criptomoneda

Cada componente deberá llenar un modelo con la siguiente estructura:
```json
{
"date": "2025-02-26T17:00:00",
"name": "Bitcoin",
"ticker_symbol": "BTC",
"price": {
"usd": 123.456,
"mxn": 123.456
}
}

```

> *Nota:* Aunque el ejemplo muestra datos de Bitcoin, se espera que cada componente corresponda a la criptomoneda indicada en su nombre (por ejemplo, BTC, ETH, XRP).
---
### Requerimientos del Proyecto
#### 1. API REST
- *Diseño del Endpoint:*
    - Define un endpoint (por ejemplo, /aggregate o /cryptos) que retorne el layout actualizado con los datos obtenidos.
    - La respuesta debe seguir las convenciones RESTful y estar estructurada en JSON.
- *Simulación del Layout:*
    - Simula la obtención del JSON de configuración (layout) de los componentes a actualizar.
#### 2. Proveedores de Datos
- *Implementación de Proveedores:*
    - Implementa proveedores para obtener los valores de las criptomonedas. Se pueden utilizar las siguientes fuentes (elige las que sean gratuitas y estén disponibles):
        - *Bitso:* [Documentación Bitso API](https://docs.bitso.com/bitso-api/docs/ticker)
        - *CoinMarketCap*
        - *Coinbase*
- *Integración:*
    - Cada proveedor debe ser llamado de forma independiente y los resultados deben ser integrados en el modelo de cada componente.
#### 3. Concurrencia
- *Obtención Concurrente:*
    - Utiliza goroutines y channels (o mecanismos equivalentes) para llamar a los proveedores en paralelo.
    - Coordina los resultados de forma segura, evitando condiciones de carrera.
- *Gestión de Contexto:*
    - Emplea el package context para manejar cancelaciones y timeouts en las solicitudes a los proveedores.
#### 4. Buenas Prácticas y Arquitectura
- *Estructuración del Código:*
    - Aplica principios de inyección de dependencias, interfaces y composición para mantener el código desacoplado y fácilmente testeable.
- *Manejo de Errores:*
    - Implementa un manejo robusto de errores para cubrir fallos en la comunicación con los proveedores o en la integración de datos.
- *Testing y Documentación:*
    - Escribe pruebas unitarias que validen la funcionalidad crítica, especialmente el manejo concurrente.
    - Documenta el proyecto y las decisiones de diseño en un README, incluyendo instrucciones para ejecutar y probar el servicio.
      
---

### Resumen de la Tarea

1. *Diseñar el endpoint REST* que devuelva la configuración (layout) con datos actualizados.
2. *Simular la carga del layout* y la integración con los proveedores de datos.
3. *Implementar la obtención de valores de criptomonedas* de manera concurrente desde al menos uno de los proveedores (Bitso, CoinMarketCap, Coinbase).
4. *Integrar los datos* en el modelo de cada componente y responder en formato JSON.
5. *Aplicar buenas prácticas de Go* en cuanto a estructuración, manejo de errores, uso de contextos y pruebas unitarias.
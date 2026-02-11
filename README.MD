# Crypto Aggregator SVC

## Run locally
```bash
go run ./cmd/api
# GET http://localhost:3000/aggregate


```
> OR
```
docker build -t crypto-aggregator .
docker run --rm -p 3000:3000 crypto-aggregator

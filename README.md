### Run app
- Rename `.env.sample` file to `.env` and set the environment variables
- `go run ./cmd/ecoorders`

### Endpoint
- Payload Example
````
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"items":[{"product_id":10,"quantity":2},{"product_id":20,"quantity":1}]}'
````

- Response
```
{"order_id":1,"status":"processed","total":300}
```

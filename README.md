# Go JWT

## Setup
```
docker-compose up --build
```
Server should now be started on localhost:8080

## Testing
Note: This project was built locally with go version go1.21.4 darwin/arm64

There is only the handler unit tests which you can run with `go test ./...`

## Example requests
POST /signup
```
curl --location 'localhost:8080/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstName": "Anhkhoi",
    "lastName": "Vu-Nguyen",
    "email": "anhkhoi@email.com",
    "password": "password"
}'
```

POST /login
```
curl --location 'localhost:8080/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "anhkhoi@email.com",
    "password": "password"
}'
```

GET /users
```
curl --location 'localhost:8080/users' \
--header 'X-Authentication-Token: $TOKEN'
```

PUT /users
```
curl --location --request PUT 'localhost:8080/users' \
--header 'X-Authentication-Token: $TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "firstName": "Dapper",
    "lastName": "Labs"
}'
```
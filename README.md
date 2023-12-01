# Go JWT

## Setup
1. This project was built locally with go version 1.21 
2. Setup DB container
```
docker build -t dapper-psql .
docker run -d -p 5432:5432 dapper-psql
```
3. Setup API server
```
go build main.go
./main
```
4. Server should now be started on localhost:8080

## Testing
There is only handler unit tests you can run with `go test ./...`

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
--header 'X-Authentication-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE0ODc2NTYsInVzZXJfZW1haWwiOiJhbmhraG9pQGVtYWlsLmNvbSJ9.LbpPsTRGiAnG4Zy92tvmTiLQ81ewu3ow8D2R8XBHJII' \
--header 'Content-Type: application/json' \
--data '{
    "firstName": "Dapper",
    "lastName": "Labs"
}'
```
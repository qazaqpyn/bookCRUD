# Simple CRUD bookCRUD app

## Technologies 
- JWT Authentication 
- Loggin logrus
- Swagger
- Cookies Authentication
- Audit Log service gRPC Client or Message Queue

[Authentication processes use message queue and Book processes use gRPC for audit log]

## Start RabbitMQ with Docker
```
docker run -it --rm --name rabbitmq -p 5672:5672 rabbitmq
```

## Build & Run
```
go build -o app cmd/main.go && ./app
```



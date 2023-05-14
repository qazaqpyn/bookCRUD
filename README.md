# Simple CRUD bookCRUD app

## Technologies 
- JWT Authentication 
- Loggin logrus
- Swagger
- Cookies Authentication
- Audit Log service gRPC Client with Message Queue

## Start RabbitMQ with Docker
```
docker run -it --rm --name rabbitmq -p 5672:5672 rabbitmq
```

## Build & Run
```
go build -o app cmd/main.go && ./app
```



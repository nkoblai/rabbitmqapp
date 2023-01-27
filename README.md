# rabbitmqapp
**Prerequisites**:  go v1.19.2 or higher, docker

**Client argumens**:
```Usage of /tmp/go-build2796997556/b001/exe/main:
  -c string
        clientID as UUID
  -rmqh string
        RabbitMQ host
  -rmqn string
        RabbitMQ queue name
  -rt string
        request type, possible values: AddItem, RemoveItem, GetItem, GetAllItems
  -v string
        item value
```
**Server arguments**:
```Usage of /tmp/go-build2880048882/b001/exe/main:
-rmqh string
    RabbitMQ host
-rmqn string
    RabbitMQ queue name
```
**Run**:
1. Start RabbitMQ locally using: `docker run -d --name rabbitmql -p 5672:5672 -p 15672:15672 bbitmq:3.11-management`.
2. Enter server folder and start the server: `go run cmd/server/main.go` with argumens provided. Server will create a queue inside RabbitMQ on startup.
3. Enter client folder and start the client app by running: `go run cmd/client/main.go` with argumens provided.
4. During request processing server will log all actions into `server_logs.txt`.

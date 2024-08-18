# go-micro

A simple demo of microservices implemented in Go. <br/>
Tech stack: Go, PostgreSQL, MongoDB, RabbitMQ, RESTful API, RPC/gRPC, Docker, kubernetes.
#
### Broker Service
A broker service as the entry point for the go-micro project.

### Authentication Service
A authentication service implemented using JWT access tokens and refresh tokens.

### Logger Service
A logger service to save logs to MongoDB.

### Mail service
A mail service to send a html/plain email.

### Listener service
A listener service to listen RabbitMQ message.

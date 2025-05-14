module bff-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/KusakinDev/Catering-Menu-Service v0.0.0-20230830123456-c1a1b1b1b1b1
	github.com/lyfoore/catering_order_microservice v0.0.0-20230830123456-c2a2b2b2b2b2
	google.golang.org/grpc v1.58.2
	google.golang.org/protobuf v1.31.0
)

replace (
	github.com/KusakinDev/Catering-Menu-Service => ../Catering-Menu-Service
	github.com/lyfoore/catering_order_microservice => ../catering_order_microservice
)
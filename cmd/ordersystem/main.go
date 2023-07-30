package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"github.com/twsm000/goexp-cleanarch/configs"
	"github.com/twsm000/goexp-cleanarch/internal/event/handler"
	"github.com/twsm000/goexp-cleanarch/internal/infra/graph"
	"github.com/twsm000/goexp-cleanarch/internal/infra/grpc/pb"
	"github.com/twsm000/goexp-cleanarch/internal/infra/grpc/service"
	"github.com/twsm000/goexp-cleanarch/internal/infra/web/webserver"
	"github.com/twsm000/goexp-cleanarch/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	ws := webserver.NewWebServer(configs.WebServerPort)
	wsOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	ws.AddHandler(webserver.NewHandler("/orders", webserver.Post, wsOrderHandler.Create))
	ws.AddHandler(webserver.NewHandler("/orders", webserver.Get, wsOrderHandler.ListAll))
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go ws.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

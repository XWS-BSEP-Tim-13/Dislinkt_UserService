package startup

import (
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/api"
	user "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/persistence"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/startup/config"
	otgo "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

var grpcGatewayTag = otgo.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	userStore := server.initUserStore(mongoClient)
	connectionsStore := server.initConnectionStore(mongoClient)
	userService := server.initUserService(userStore, connectionsStore)

	userHandler := server.initUserHandler(userService)

	server.startGrpcServer(userHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.UserDBHost, server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initConnectionStore(client *mongo.Client) domain.ConnectionRequestStore {
	store := persistence.NewConnectionsMongoDBStore(client)
	store.DeleteAll()

	for _, connection := range connections {
		err := store.Insert(connection)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initUserStore(client *mongo.Client) domain.UserStore {
	store := persistence.NewUserMongoDBStore(client)
	store.DeleteAll()
	for _, user := range users {
		err := store.Insert(user)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initUserService(store domain.UserStore, conStore domain.ConnectionRequestStore) *application.UserService {
	return application.NewUserService(store, conStore)
}

func (server *Server) initUserHandler(service *application.UserService) *api.UserHandler {
	return api.NewUserHandler(service)
}

func (server *Server) startGrpcServer(userHandler *api.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

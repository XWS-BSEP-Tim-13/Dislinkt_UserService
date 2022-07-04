package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/api"
	user "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/persistence"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/startup/config"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/util"
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

const (
	serverCertFile = "cert/cert.pem"
	serverKeyFile  = "cert/key.pem"
	clientCertFile = "cert/client-cert.pem"
)

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	logger := logger.InitLogger("user-service", context.TODO())

	mongoClient := server.initMongoClient()
	userStore := server.initUserStore(mongoClient)
	connectionsStore := server.initConnectionStore(mongoClient)
	userService := server.initUserService(userStore, connectionsStore, logger)
	goValidator := server.initGoValidator()
	userHandler := server.initUserHandler(userService, goValidator, logger)

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

func (server *Server) initUserService(store domain.UserStore, conStore domain.ConnectionRequestStore, logger *logger.Logger) *application.UserService {
	return application.NewUserService(store, conStore, logger)
}

func (server *Server) initGoValidator() *util.GoValidator {
	return util.NewGoValidator()
}

func (server *Server) initUserHandler(service *application.UserService, goValidator *util.GoValidator, logger *logger.Logger) *api.UserHandler {
	return api.NewUserHandler(service, goValidator, logger)
}

func (server *Server) startGrpcServer(userHandler *api.UserHandler) {
	/*cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	pemClientCA, err := ioutil.ReadFile(clientCertFile)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequestClientCert,
		ClientCAs:    certPool,
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(config)),
	}*/

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

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
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/tracer"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/util"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	otgo "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type Server struct {
	config *config.Config
	tracer otgo.Tracer
	closer io.Closer
}

const (
	serverCertFile = "cert/cert.pem"
	serverKeyFile  = "cert/key.pem"
	clientCertFile = "cert/client-cert.pem"
)

func NewServer(config *config.Config) *Server {
	tracer, closer := tracer.Init()
	otgo.SetGlobalTracer(tracer)

	return &Server{
		config: config,
		tracer: tracer,
		closer: closer,
	}
}

func (server *Server) Start() {
	logger := logger.InitLogger("user-service", context.TODO())

	mongoClient := server.initMongoClient()
	userStore := server.initUserStore(mongoClient)
	connectionsStore := server.initConnectionStore(mongoClient)
	notificationStore := server.initNotificationStore(mongoClient)
	userService := server.initUserService(userStore, connectionsStore, logger, notificationStore)
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
	store.DeleteAll(context.TODO())

	for _, connection := range connections {
		err := store.Insert(context.TODO(), connection)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initNotificationStore(client *mongo.Client) domain.NotificationStore {
	store := persistence.NewNotificationMongoDBStore(client)
	store.DeleteAll()

	for _, connection := range notifications {
		err := store.Insert(connection)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initUserStore(client *mongo.Client) domain.UserStore {
	store := persistence.NewUserMongoDBStore(client)
	store.DeleteAll(context.TODO())
	for _, user := range users {
		err := store.Insert(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initUserService(store domain.UserStore, conStore domain.ConnectionRequestStore, logger *logger.Logger, notificationStore domain.NotificationStore) *application.UserService {
	return application.NewUserService(store, conStore, logger, notificationStore)
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
	}*/

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(server.tracer)),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(server.tracer)),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

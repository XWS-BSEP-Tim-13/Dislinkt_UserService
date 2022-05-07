package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/application"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	userPb := mapUserToPB(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) FindByFilter(ctx context.Context, request *pb.UserFilter) (*pb.GetAllResponse, error) {
	filter := request.Filter
	users, err := handler.service.FindByFilter(filter)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUserToPB(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) GetRequestsForUser(ctx context.Context, request *pb.GetRequest) (*pb.UserRequests, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	requests, err := handler.service.GetRequestsForUser(id)
	response := &pb.UserRequests{
		Requests: []*pb.ConnectionRequest{},
	}
	for _, request := range requests {
		fmt.Printf("Request: %s, id to: %s\n", request.To.FirstName, request.To.LastName)
		current := mapConnectionRequest(request)
		response.Requests = append(response.Requests, current)
	}
	return response, nil
}

func (handler *UserHandler) AcceptConnectionRequest(ctx context.Context, request *pb.GetRequest) (*pb.ConnectionResponse, error) {
	connectionId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	handler.service.AcceptConnection(connectionId)
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) RequestConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	idFrom, err := primitive.ObjectIDFromHex(request.Connection.IdFrom)
	idTo, err1 := primitive.ObjectIDFromHex(request.Connection.IdTo)
	fmt.Printf("Id from: %s, id to: %s\n", idFrom, idTo)
	if err != nil || err1 != nil {
		return nil, err
	}
	handler.service.RequestConnection(idFrom, idTo)
	fmt.Printf("Returning to func")
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) CheckIfUserCanReadPosts(ctx context.Context, request *pb.ConnectionBody) (*pb.ReadPostsResponse, error) {
	idFrom, err := primitive.ObjectIDFromHex(request.Connection.IdFrom)
	idTo, err1 := primitive.ObjectIDFromHex(request.Connection.IdTo)
	fmt.Printf("Id from: %s, id to: %s\n", idFrom, idTo)
	if err != nil || err1 != nil {
		return nil, err
	}
	isReadable, err1 := handler.service.CheckIfUserCanReadPosts(idFrom, idTo)

	if err1 != nil {
		return nil, err
	}
	response := &pb.ReadPostsResponse{
		IsReadable: isReadable,
	}
	return response, nil
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	//ctx = tracer.ContextWithSpan(context.Background(), span)
	for _, user := range users {
		current := mapUserToPB(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

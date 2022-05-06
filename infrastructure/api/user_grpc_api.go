package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/api/dto"
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
	userPb := mapUser(user)
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
		current := mapUser(user)
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
		current := mapConnectionRequest(request)
		response.Requests = append(response.Requests, current)
	}
	return response, nil
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
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) UpdatePersonalInfo(ctx context.Context, request *pb.UserInfoUpdate) (*pb.UserInfoUpdateResponse, error) {
	id, _ := primitive.ObjectIDFromHex(request.UserInfo.Id)
	userInfo := dto.NewUserInfo(id, request.UserInfo.FirstName, request.UserInfo.LastName, enum.Gender(request.UserInfo.Gender), request.UserInfo.DateOfBirth.AsTime(),
		request.UserInfo.Email, request.UserInfo.PhoneNumber, request.UserInfo.Biography)
	user := dto.NewUserFromUserInfo(*userInfo)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.UserInfo.Id

	if _, err := handler.service.UpdatePersonalInfo(user); err != nil {
		return nil, err
	}

	return response, nil
}

func (handler *UserHandler) AddExperience(ctx context.Context, request *pb.ExperienceUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.ExperienceUpdate.UserId
	exp := mapExperience(request.ExperienceUpdate.Experience)
	expId, _ := primitive.ObjectIDFromHex(request.ExperienceUpdate.UserId)

	if err := handler.service.AddExperience(exp, expId); err != nil {
		return nil, err
	}

	return response, nil
}

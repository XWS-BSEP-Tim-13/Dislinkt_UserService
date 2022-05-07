package api

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/api/dto"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto/user_info/user_info"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfoHandler struct {
	pb.UnimplementedUserInfoServiceServer
	service *application.UserInfoService
}

func NewUserInfoHandler(service *application.UserInfoService) *UserInfoHandler {
	return &UserInfoHandler{
		service: service,
	}
}

func (handler *UserInfoHandler) UpdatePersonalInfo(ctx context.Context, request *pb.UserInfoUpdate) (*pb.UserInfoUpdateResponse, error) {
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

func (handler *UserInfoHandler) AddExperience(ctx context.Context, request *pb.ExperienceUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.ExperienceUpdate.UserId
	exp := mapExperience(request.ExperienceUpdate.Experience)
	expId, _ := primitive.ObjectIDFromHex(request.ExperienceUpdate.UserId)

	if err := handler.service.AddExperience(exp, expId); err != nil {
		return nil, err
	}

	return response, nil
}

func (handler *UserInfoHandler) AddEducation(ctx context.Context, request *pb.EducationUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.EducationUpdate.UserId
	education := mapEducation(request.EducationUpdate.Education)
	expId, _ := primitive.ObjectIDFromHex(request.EducationUpdate.UserId)

	if err := handler.service.AddEducation(education, expId); err != nil {
		return nil, err
	}

	return response, nil
}

func (handler *UserInfoHandler) AddSkill(ctx context.Context, request *pb.SkillsUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.Skill.Skill
	userId, _ := primitive.ObjectIDFromHex(request.Skill.UserId)

	if err := handler.service.AddSkill(request.Skill.Skill, userId); err != nil {
		return nil, err
	}

	return response, nil
}

func (handler *UserInfoHandler) AddInterest(ctx context.Context, request *pb.InterestsUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.Interest.CompanyId
	userId, _ := primitive.ObjectIDFromHex(request.Interest.UserId)
	companyId, _ := primitive.ObjectIDFromHex(request.Interest.CompanyId)
	if err := handler.service.AddInterest(companyId, userId); err != nil {
		return nil, err
	}

	return response, nil
}

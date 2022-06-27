package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/jwt"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/logging"

	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/api/dto"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service     *application.UserService
	goValidator *util.GoValidator
	logger      *logger.Logger
}

func NewUserHandler(service *application.UserService, goValidator *util.GoValidator, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service:     service,
		goValidator: goValidator,
		logger:      logger,
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
		handler.logger.ErrorMessage("Action: u/:id")
		return nil, err
	}
	userPb := mapUserToPB(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	handler.logger.InfoMessage("Action: u/:id")
	return response, nil
}

func (handler *UserHandler) FindByFilter(ctx context.Context, request *pb.UserFilter) (*pb.GetAllResponse, error) {
	filter := request.Filter
	users, err := handler.service.FindByFilter(filter)
	if err != nil {
		handler.logger.ErrorMessage("Action: FU")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUserToPB(user)
		response.Users = append(response.Users, current)
	}

	handler.logger.InfoMessage("Action: FU")
	return response, nil
}

func (handler *UserHandler) GetRequestsForUser(ctx context.Context, request *pb.GetRequest) (*pb.UserRequests, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
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
	handler.logger.InfoMessage("User: " + username + " | Action: GCR")
	return response, nil
}

func (handler *UserHandler) AcceptConnectionRequest(ctx context.Context, request *pb.GetRequest) (*pb.ConnectionResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	connectionId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	handler.service.AcceptConnection(connectionId)
	handler.logger.InfoMessage("User: " + username + " | Action: ACR")
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) DeleteConnectionRequest(ctx context.Context, request *pb.GetRequest) (*pb.ConnectionResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	connectionId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	handler.service.DeleteConnectionRequest(connectionId)
	handler.logger.InfoMessage("User: " + username + " | Action: DCR")
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) DeleteConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	fmt.Printf("Request: %s, id to: %s\n", request.Connection.IdFrom, request.Connection.IdTo)
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	idFrom, err := primitive.ObjectIDFromHex(request.Connection.IdFrom)
	idTo, err1 := primitive.ObjectIDFromHex(request.Connection.IdTo)
	fmt.Printf("Id from: %s, id to: %s\n", idFrom, idTo)
	if err != nil || err1 != nil {
		return nil, err
	}
	err = handler.service.DeleteConnection(idFrom, idTo)
	if err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: DC/: " + idTo.String())
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: DC/: " + idTo.String())
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) RequestConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	idFrom, err := primitive.ObjectIDFromHex(request.Connection.IdFrom)
	idTo, err1 := primitive.ObjectIDFromHex(request.Connection.IdTo)
	fmt.Printf("Id from: %s, id to: %s\n", idFrom, idTo)
	if err != nil || err1 != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: CCR/: " + idTo.String())
		return nil, err
	}
	handler.service.RequestConnection(idFrom, idTo)
	handler.logger.InfoMessage("User: " + username + " | Action: CCR/: " + idTo.String())
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) GetUsernames(ctx context.Context, request *pb.ConnectionResponse) (*pb.UserConnectionUsernames, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		var connUsernames []string
		response := &pb.UserConnectionUsernames{
			Usernames: connUsernames,
		}
		return response, nil
	}
	fmt.Println("Conn username", username)
	connUsernames, err := handler.service.GetConnectionUsernamesForUser(username)
	response := &pb.UserConnectionUsernames{
		Usernames: connUsernames,
	}
	if err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: GCU/" + username)
		return nil, err
	}
	handler.logger.InfoMessage("User: " + username + " | Action: GCU/" + username)
	return response, nil
}

func (handler *UserHandler) ChangeAccountPrivacy(ctx context.Context, requst *pb.ReadPostsResponse) (*pb.ConnectionResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		return nil, err
	}
	err = handler.service.ChangeAccountPrivacy(username, requst.IsReadable)
	if err != nil {
		return nil, err
	}
	response := &pb.ConnectionResponse{}
	return response, nil
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
		handler.logger.ErrorMessage("Action: GU")
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

	handler.logger.InfoMessage("Action: GU")
	return response, nil
}

func (handler *UserHandler) UpdatePersonalInfo(ctx context.Context, request *pb.UserInfoUpdate) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	id, _ := primitive.ObjectIDFromHex(request.UserInfo.Id)
	userInfo := dto.NewUserInfo(id, request.UserInfo.FirstName, request.UserInfo.LastName, enum.Gender(request.UserInfo.Gender), request.UserInfo.DateOfBirth.AsTime(),
		request.UserInfo.Email, request.UserInfo.PhoneNumber, request.UserInfo.Biography)
	user := dto.NewUserFromUserInfo(*userInfo)

	validationErr := handler.goValidator.Validator.Struct(user)
	if validationErr != nil {
		handler.goValidator.PrintValidationErrors(validationErr)
		handler.logger.WarningMessage("User: " + username + " | Action: UP")
		return nil, status.Error(500, validationErr.Error())
	}

	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.UserInfo.Id

	if _, err := handler.service.UpdatePersonalInfo(user); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: UP")
		return nil, err
	}
	handler.logger.InfoMessage("User: " + username + " | Action: UP")
	return response, nil
}

func (handler *UserHandler) AddExperience(ctx context.Context, request *pb.ExperienceUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.ExperienceUpdate.UserId
	exp := mapExperience(request.ExperienceUpdate.Experience)

	validationErr := handler.goValidator.Validator.Struct(exp)
	if validationErr != nil {
		handler.logger.WarningMessage("User: " + username + " | Action: AExp")
		handler.goValidator.PrintValidationErrors(validationErr)
		return nil, status.Error(500, validationErr.Error())
	}

	expId, _ := primitive.ObjectIDFromHex(request.ExperienceUpdate.UserId)

	if err := handler.service.AddExperience(exp, expId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AExp")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: AExp")
	return response, nil
}

func (handler *UserHandler) AddEducation(ctx context.Context, request *pb.EducationUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.EducationUpdate.UserId
	education := mapEducation(request.EducationUpdate.Education)

	validationErr := handler.goValidator.Validator.Struct(education)
	if validationErr != nil {
		handler.goValidator.PrintValidationErrors(validationErr)
		handler.logger.WarningMessage("User: " + username + " | Action: AEdu")
		return nil, status.Error(500, validationErr.Error())
	}

	expId, _ := primitive.ObjectIDFromHex(request.EducationUpdate.UserId)

	if err := handler.service.AddEducation(education, expId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AEdu")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: AEdu")
	return response, nil
}

func (handler *UserHandler) CreateUser(ctx context.Context, request *pb.NewUser) (*pb.NewUser, error) {
	fmt.Println((*request).User)
	user := mapUserToDomain(request.User)

	err := handler.goValidator.Validator.Struct(user)
	if err != nil {
		handler.logger.WarningMessage("Action: CU")
		return nil, status.Error(500, err.Error())
	}

	newUser, err := handler.service.CreateNewUser(user)
	if err != nil {
		handler.logger.ErrorMessage("Action: CU")
		return nil, status.Error(400, err.Error())
	}

	response := &pb.NewUser{
		User: mapUserToPB(newUser),
	}

	handler.logger.InfoMessage("Action: CU")
	return response, nil
}

func (handler *UserHandler) ActivateAccount(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	email := request.Email

	resp, err := handler.service.ActivateAccount(email)
	if err != nil {
		handler.logger.ErrorMessage("Action: AA " + email)
		return nil, status.Error(500, err.Error())
	}

	response := &pb.ActivateAccountResponse{
		Message: resp,
	}

	handler.logger.InfoMessage("Action: AA " + email)
	return response, nil
}

func (handler *UserHandler) AddSkill(ctx context.Context, request *pb.SkillsUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.Skill.Skill
	userId, _ := primitive.ObjectIDFromHex(request.Skill.UserId)

	validationErr := handler.goValidator.ValidateSkill(request.Skill.Skill)
	if validationErr != nil {
		handler.goValidator.PrintValidationErrors(validationErr)
		handler.logger.WarningMessage("User: " + username + " | Action: AS " + request.Skill.Skill)
		return nil, status.Error(500, validationErr.Error())
	}

	if err := handler.service.AddSkill(request.Skill.Skill, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AS " + request.Skill.Skill)
		return nil, err
	}
	handler.logger.InfoMessage("User: " + username + " | Action: AS " + request.Skill.Skill)

	return response, nil
}

func (handler *UserHandler) RemoveSkill(ctx context.Context, request *pb.RemoveSkillRequest) (*pb.RemoveSkillResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	userId, err := primitive.ObjectIDFromHex(request.Skill.UserId)
	if err != nil {
		return nil, status.Error(500, "Error parsing id.")
	}

	if err := handler.service.RemoveSkill(request.Skill.Skill, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: RmS " + request.Skill.Skill)
		return nil, err
	}

	response := &pb.RemoveSkillResponse{
		Skill: request.Skill.Skill,
	}

	handler.logger.InfoMessage("User: " + username + " | Action: RmS " + request.Skill.Skill)
	return response, nil
}

func (handler *UserHandler) AddInterest(ctx context.Context, request *pb.InterestsUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.Interest.CompanyId
	userId, _ := primitive.ObjectIDFromHex(request.Interest.UserId)
	companyId, _ := primitive.ObjectIDFromHex(request.Interest.CompanyId)
	if err := handler.service.AddInterest(companyId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AInt ")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: AInt ")
	return response, nil
}

func (handler *UserHandler) DeleteExperience(ctx context.Context, request *pb.DeleteExperienceRequest) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.DeleteExperience.ExperienceId
	userId, _ := primitive.ObjectIDFromHex(request.DeleteExperience.UserId)
	experienceId, _ := primitive.ObjectIDFromHex(request.DeleteExperience.ExperienceId)
	if err := handler.service.DeleteExperience(experienceId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: RmExp ")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: RmExp ")
	return response, nil
}

func (handler *UserHandler) DeleteEducation(ctx context.Context, request *pb.DeleteEducationRequest) (*pb.UserInfoUpdateResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.DeleteEducation.EducationId
	userId, _ := primitive.ObjectIDFromHex(request.DeleteEducation.UserId)
	educationId, _ := primitive.ObjectIDFromHex(request.DeleteEducation.EducationId)
	if err := handler.service.DeleteEducation(educationId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: RmEdu ")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: RmEdu ")
	return response, nil
}

func (handler *UserHandler) RemoveInterest(ctx context.Context, request *pb.RemoveInterestRequest) (*pb.RemoveInterestResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	userId, err := primitive.ObjectIDFromHex(request.Interest.UserId)
	companyId, err := primitive.ObjectIDFromHex(request.Interest.CompanyId)
	if err != nil {
		return nil, status.Error(500, "Error parsing id.")
	}

	if err := handler.service.RemoveInterest(companyId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: RmInt ")
		return nil, err
	}

	response := &pb.RemoveInterestResponse{
		CompanyId: request.Interest.CompanyId,
	}

	handler.logger.InfoMessage("User: " + username + " | Action: RmInt ")
	return response, nil
}

func (handler *UserHandler) GetByUsername(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	username := request.Id
	user, err := handler.service.GetByUsername(username)
	if err != nil {
		handler.logger.ErrorMessage("Action: u/" + username)
		return nil, err
	}
	userPb := mapUserToPB(user)
	response := &pb.GetResponse{
		User: userPb,
	}

	handler.logger.InfoMessage("Action: u/" + username)
	return response, nil
}

func (handler *UserHandler) GetByEmail(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	username := request.Id
	user, err := handler.service.GetByEmail(username)
	if err != nil {
		handler.logger.ErrorMessage("Action: u/" + username)
		return nil, err
	}
	userPb := mapUserToPB(user)
	response := &pb.GetResponse{
		User: userPb,
	}

	handler.logger.InfoMessage("Action: u/: " + username)
	return response, nil
}

package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/jwt"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/tracer"

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
	span := tracer.StartSpanFromContext(ctx, "API Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(ctx, objectId)
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
	span := tracer.StartSpanFromContext(ctx, "API FindByFilter")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := request.Filter
	users, err := handler.service.FindByFilter(ctx, filter)
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
	span := tracer.StartSpanFromContext(ctx, "API GetRequestsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	requests, err := handler.service.GetRequestsForUser(ctx, id)
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

func (handler *UserHandler) DeleteConnectionRequest(ctx context.Context, request *pb.GetRequest) (*pb.ConnectionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API DeleteConnectionRequest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	connectionId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	handler.service.DeleteConnectionRequest(ctx, connectionId)
	handler.logger.InfoMessage("User: " + username + " | Action: DCR")
	return new(pb.ConnectionResponse), nil
}

func (handler *UserHandler) ChangeAccountPrivacy(ctx context.Context, request *pb.ReadPostsResponse) (*pb.ConnectionResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	span := tracer.StartSpanFromContext(ctx, "API ChangeAccountPrivacy")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err = handler.service.ChangeAccountPrivacy(ctx, username, request.IsReadable)
	if err != nil {
		return nil, err
	}
	response := &pb.ConnectionResponse{}
	return response, nil
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API GetAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	users, err := handler.service.GetAll(ctx)
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
	span := tracer.StartSpanFromContext(ctx, "API UpdatePersonalInfo")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

	if _, err := handler.service.UpdatePersonalInfo(ctx, user); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: UP")
		return nil, err
	}
	handler.logger.InfoMessage("User: " + username + " | Action: UP")
	return response, nil
}

func (handler *UserHandler) AddExperience(ctx context.Context, request *pb.ExperienceUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API AddExperience")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

	if err := handler.service.AddExperience(ctx, exp, expId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AExp")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: AExp")
	return response, nil
}

func (handler *UserHandler) AddEducation(ctx context.Context, request *pb.EducationUpdateRequest) (*pb.UserInfoUpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API AddEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

	if err := handler.service.AddEducation(ctx, education, expId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AEdu")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: AEdu")
	return response, nil
}

func (handler *UserHandler) GetNotificationsForUser(ctx context.Context, request *pb.ConnectionResponse) (*pb.NotificationResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		fmt.Println(err)
	}
	span := tracer.StartSpanFromContext(ctx, "API AddEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	notifications, err := handler.service.GetNotificationsForUser(ctx, username)
	if err != nil {
		return nil, err
	}
	response := &pb.NotificationResponse{
		Notification: []*pb.Notification{},
	}

	for _, notification := range notifications {
		current := mapNotificationToPB(notification)
		response.Notification = append(response.Notification, current)
	}
	return response, nil
}

func (handler *UserHandler) UpdateUserNotificationAlert(ctx context.Context, request *pb.ReadPostsResponse) (*pb.ConnectionResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	handler.service.UpdateNotificationAlert(username, request.IsReadable)
	response := &pb.ConnectionResponse{}
	return response, nil
}

func (handler *UserHandler) CreateNotification(ctx context.Context, request *pb.NotificationRequest) (*pb.ConnectionResponse, error) {
	notification := mapPbToNotificationDomain(request.Notification)
	fmt.Println(notification)
	handler.service.SaveNotification(notification)
	response := &pb.ConnectionResponse{}
	return response, nil
}

func (handler *UserHandler) MessageNotification(ctx context.Context, request *pb.Connection) (*pb.ConnectionResponse, error) {
	notification := notificationCreate(request.IdFrom, request.IdTo)
	fmt.Println(notification)
	handler.service.SaveNotification(notification)
	response := &pb.ConnectionResponse{}
	return response, nil
}

func (handler *UserHandler) CreateUser(ctx context.Context, request *pb.NewUser) (*pb.NewUser, error) {
	span := tracer.StartSpanFromContext(ctx, "API CreateUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println((*request).User)
	user := mapUserToDomain(request.User)

	err := handler.goValidator.Validator.Struct(user)
	if err != nil {
		handler.logger.WarningMessage("Action: CU")
		return nil, status.Error(500, err.Error())
	}

	newUser, err := handler.service.CreateNewUser(ctx, user)
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
	span := tracer.StartSpanFromContext(ctx, "API ActivateAccount")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	email := request.Email

	resp, err := handler.service.ActivateAccount(ctx, email)
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
	span := tracer.StartSpanFromContext(ctx, "API AddSkill")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

	if err := handler.service.AddSkill(ctx, request.Skill.Skill, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AS " + request.Skill.Skill)
		return nil, err
	}
	handler.logger.InfoMessage("User: " + username + " | Action: AS " + request.Skill.Skill)

	return response, nil
}

func (handler *UserHandler) RemoveSkill(ctx context.Context, request *pb.RemoveSkillRequest) (*pb.RemoveSkillResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API RemoveSkill")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	userId, err := primitive.ObjectIDFromHex(request.Skill.UserId)
	if err != nil {
		return nil, status.Error(500, "Error parsing id.")
	}

	if err := handler.service.RemoveSkill(ctx, request.Skill.Skill, userId); err != nil {
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
	span := tracer.StartSpanFromContext(ctx, "API AddInterest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.Interest.CompanyId
	userId, _ := primitive.ObjectIDFromHex(request.Interest.UserId)
	companyId, _ := primitive.ObjectIDFromHex(request.Interest.CompanyId)
	if err := handler.service.AddInterest(ctx, companyId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: AInt ")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: AInt ")
	return response, nil
}

func (handler *UserHandler) DeleteExperience(ctx context.Context, request *pb.DeleteExperienceRequest) (*pb.UserInfoUpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API DeleteExperience")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.DeleteExperience.ExperienceId
	userId, _ := primitive.ObjectIDFromHex(request.DeleteExperience.UserId)
	experienceId, _ := primitive.ObjectIDFromHex(request.DeleteExperience.ExperienceId)
	if err := handler.service.DeleteExperience(ctx, experienceId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: RmExp ")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: RmExp ")
	return response, nil
}

func (handler *UserHandler) DeleteEducation(ctx context.Context, request *pb.DeleteEducationRequest) (*pb.UserInfoUpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API DeleteEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	response := new(pb.UserInfoUpdateResponse)
	response.Id = request.DeleteEducation.EducationId
	userId, _ := primitive.ObjectIDFromHex(request.DeleteEducation.UserId)
	educationId, _ := primitive.ObjectIDFromHex(request.DeleteEducation.EducationId)
	if err := handler.service.DeleteEducation(ctx, educationId, userId); err != nil {
		handler.logger.ErrorMessage("User: " + username + " | Action: RmEdu ")
		return nil, err
	}

	handler.logger.InfoMessage("User: " + username + " | Action: RmEdu ")
	return response, nil
}

func (handler *UserHandler) RemoveInterest(ctx context.Context, request *pb.RemoveInterestRequest) (*pb.RemoveInterestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "API RemoveInterest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username, _ := jwt.ExtractUsernameFromToken(ctx)
	userId, err := primitive.ObjectIDFromHex(request.Interest.UserId)
	companyId, err := primitive.ObjectIDFromHex(request.Interest.CompanyId)
	if err != nil {
		return nil, status.Error(500, "Error parsing id.")
	}

	if err := handler.service.RemoveInterest(ctx, companyId, userId); err != nil {
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
	span := tracer.StartSpanFromContext(ctx, "API GetByUsername")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username := request.Id
	user, err := handler.service.GetByUsername(ctx, username)
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
	span := tracer.StartSpanFromContext(ctx, "API GetByEmail")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username := request.Id
	user, err := handler.service.GetByEmail(ctx, username)
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

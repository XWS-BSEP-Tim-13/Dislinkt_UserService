package api

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func mapUserToPB(user *domain.RegisteredUser) *pb.User {
	userPb := &pb.User{
		Id:          user.Id.Hex(),
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      pb.User_Gender(user.Gender),
		DateOfBirth: timestamppb.New(user.DateOfBirth),
		Biography:   user.Biography,
		IsPrivate:   user.IsPrivate,
	}

	for _, experience := range user.Experiences {
		userPb.Experiences = append(userPb.Experiences, &pb.Experience{
			Id:                 experience.Id.Hex(),
			Title:              experience.Title,
			EmploymentType:     pb.Experience_EmploymentType(experience.EmploymentType),
			CompanyName:        experience.CompanyName,
			Location:           experience.Location,
			IsCurrentlyWorking: experience.IsCurrentlyWorking,
			StartDate:          timestamppb.New(experience.StartDate),
			EndDate:            timestamppb.New(experience.EndDate),
			Industry:           experience.Industry,
			Description:        experience.Description,
		})
	}

	for _, education := range user.Educations {
		userPb.Educations = append(userPb.Educations, &pb.Education{
			Id:           education.Id.Hex(),
			School:       education.School,
			Degree:       pb.Education_Degree(education.Degree),
			FieldOfStudy: education.FieldOfStudy,
			StartDate:    timestamppb.New(education.StartDate),
			EndDate:      timestamppb.New(education.EndDate),
			Description:  education.Description,
		})
	}

	for _, skill := range user.Skills {
		userPb.Skills = append(userPb.Skills, skill)
	}

	for _, interest := range user.Interests {
		userPb.Interests = append(userPb.Interests, interest.Hex())
	}

	for _, connection := range user.Connections {
		userPb.Connections = append(userPb.Connections, connection)
	}

	return userPb
}

func mapUserToDomain(userPb *pb.User) *domain.RegisteredUser {
	//id, err := primitive.ObjectIDFromHex((*userPb).Id)
	//if err != nil {
	//	return &domain.RegisteredUser{}
	//}

	user := &domain.RegisteredUser{
		Username:    (*userPb).Username,
		FirstName:   (*userPb).FirstName,
		LastName:    (*userPb).LastName,
		Email:       (*userPb).Email,
		PhoneNumber: (*userPb).PhoneNumber,
		Gender:      enum.Gender((*userPb).Gender),
		DateOfBirth: (*((*userPb).DateOfBirth)).AsTime(),
		Biography:   (*userPb).Biography,
		IsPrivate:   (*userPb).IsPrivate,
	}

	user.Experiences = []domain.Experience{}
	for _, experience := range (*userPb).Experiences {
		id, err := primitive.ObjectIDFromHex(experience.Id)
		if err != nil {
			continue
		}

		user.Experiences = append(user.Experiences, domain.Experience{
			Id:                 id,
			Title:              experience.Title,
			EmploymentType:     enum.EmploymentType(experience.EmploymentType),
			CompanyName:        experience.CompanyName,
			Location:           experience.Location,
			IsCurrentlyWorking: experience.IsCurrentlyWorking,
			StartDate:          experience.StartDate.AsTime(),
			EndDate:            experience.EndDate.AsTime(),
			Industry:           experience.Industry,
			Description:        experience.Description,
		})
	}

	user.Educations = []domain.Education{}
	for _, education := range (*userPb).Educations {
		id, err := primitive.ObjectIDFromHex(education.Id)
		if err != nil {
			continue
		}

		user.Educations = append(user.Educations, domain.Education{
			Id:           id,
			School:       education.School,
			Degree:       enum.Degree(education.Degree),
			FieldOfStudy: education.FieldOfStudy,
			StartDate:    education.StartDate.AsTime(),
			EndDate:      education.EndDate.AsTime(),
			Description:  education.Description,
		})
	}

	user.Skills = []string{}
	for _, skill := range (*userPb).Skills {
		user.Skills = append(user.Skills, skill)
	}

	user.Interests = []primitive.ObjectID{}
	for _, interest := range (*userPb).Interests {
		interestId, err := primitive.ObjectIDFromHex(interest)
		if err != nil {
			continue
		}

		user.Interests = append(user.Interests, interestId)
	}

	user.Connections = []string{}
	for _, connection := range (*userPb).Connections {
		user.Connections = append(user.Connections, connection)
	}

	return user
}

func mapNotificationToPB(request *domain.Notification) *pb.Notification {
	notificationPb := &pb.Notification{
		Id:       request.Id.Hex(),
		Username: request.Username,
		Created:  timestamppb.New(request.Created),
		Type:     pb.NotificationType(request.NotificationType),
	}
	return notificationPb
}

func mapPbToNotificationDomain(request *pb.NotificationDto) *domain.Notification {
	notification := &domain.Notification{
		Id:               primitive.NewObjectID(),
		Created:          time.Now(),
		Username:         request.Username,
		NotificationType: enum.NotificationType(request.Type),
	}
	return notification
}

func mapConnectionRequest(request *domain.ConnectionRequest) *pb.ConnectionRequest {
	connectionPb := &pb.ConnectionRequest{
		Id:          request.Id.Hex(),
		From:        mapUserToPB(&request.From),
		To:          mapUserToPB(&request.To),
		RequestTime: timestamppb.New(request.RequestTime),
	}
	return connectionPb
}

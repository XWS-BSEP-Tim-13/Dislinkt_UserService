package api

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapUser(user *domain.RegisteredUser) *pb.User {
	userPb := &pb.User{
		Id:          user.Id.Hex(),
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
		userPb.Connections = append(userPb.Connections, connection.Hex())
	}
	
	return userPb
}

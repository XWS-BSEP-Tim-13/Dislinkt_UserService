package api

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mapExperience(request *pb.Experience) *domain.Experience {
	expId, _ := primitive.ObjectIDFromHex(request.Id)
	experience := &domain.Experience{
		Id:                 expId,
		Title:              request.Title,
		EmploymentType:     enum.EmploymentType(request.EmploymentType),
		CompanyName:        request.CompanyName,
		Location:           request.Location,
		IsCurrentlyWorking: request.IsCurrentlyWorking,
		StartDate:          request.StartDate.AsTime(),
		EndDate:            request.EndDate.AsTime(),
		Industry:           request.Industry,
		Description:        request.Description,
	}

	return experience
}

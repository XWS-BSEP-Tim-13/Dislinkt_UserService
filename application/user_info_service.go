package application

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfoService struct {
	store domain.UserStore
}

func NewUserInfoService(store domain.UserStore) *UserInfoService {
	return &UserInfoService{
		store: store,
	}
}

func (service *UserInfoService) AddExperience(experience *domain.Experience, userId primitive.ObjectID) error {
	experience.Id = primitive.NewObjectID()
	return service.store.AddExperience(experience, userId)
}

func (service *UserInfoService) AddEducation(education *domain.Education, userId primitive.ObjectID) error {
	education.Id = primitive.NewObjectID()
	return service.store.AddEducation(education, userId)
}

func (service *UserInfoService) AddSkill(skill string, userId primitive.ObjectID) error {
	return service.store.AddSkill(skill, userId)
}

func (service *UserInfoService) AddInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error {
	return service.store.AddInterest(companyId, userId)
}

func (service *UserInfoService) UpdatePersonalInfo(user *domain.RegisteredUser) (primitive.ObjectID, error) {
	return service.store.UpdatePersonalInfo(user)
}

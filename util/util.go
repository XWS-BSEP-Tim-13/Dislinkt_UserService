package util

import (
	"errors"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteExperience(experiences []domain.Experience, exp primitive.ObjectID) ([]domain.Experience, error) {
	for i, e := range experiences {
		if e.Id == exp {
			withoutElem := append(experiences[:i], experiences[i+1:]...)
			return withoutElem, nil
		}
	}
	err := errors.New("experience not in slice")
	return nil, err
}

func DeleteEducation(educations []domain.Education, edu primitive.ObjectID) ([]domain.Education, error) {
	for i, e := range educations {
		if e.Id == edu {
			withoutElem := append(educations[:i], educations[i+1:]...)
			return withoutElem, nil
		}
	}
	err := errors.New("education not in slice")
	return nil, err
}

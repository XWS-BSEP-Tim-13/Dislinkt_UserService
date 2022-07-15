package util

import "go.mongodb.org/mongo-driver/bson/primitive"

func ContainsId(s []primitive.ObjectID, e primitive.ObjectID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveElement(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveIdElement(s []primitive.ObjectID, i int) []primitive.ObjectID {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

package helpers

import (
	"anubis/app/DAL/entitiesDB"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserRole(project entitiesDB.MdProject, userID primitive.ObjectID) (string, error) {
	for _, member := range project.Members {
		if member.UserID == userID {
			return member.Role, nil
		}
	}
	return "", fmt.Errorf("user %s not found in project", userID.Hex())
}

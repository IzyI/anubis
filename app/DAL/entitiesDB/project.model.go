package entitiesDB

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MdProjectMember struct {
	UserID   primitive.ObjectID `bson:"userId"`
	Role     string             `bson:"role"` // admin, editor, viewer
	JoinedAt time.Time          `bson:"joinedAt"`
}
type MdProject struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Domain  string             `bson:"domain"`
	Members []MdProjectMember  `bson:"members"`
}

func (p *MdProject) GetUserRole(userID primitive.ObjectID) (string, error) {
	for _, member := range p.Members {
		if member.UserID == userID {
			return member.Role, nil
		}
	}
	return "", fmt.Errorf("user with ID %s not found", userID.Hex())
}

package entitiesDB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MdProjectMember struct {
	UserID   string    `bson:"userId"`
	Role     string    `bson:"role"` // admin, editor, viewer
	JoinedAt time.Time `bson:"joinedAt"`
}
type MdProject struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID string             `bson:"projectId"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
	Members   []MdProjectMember  `bson:"members"`
}

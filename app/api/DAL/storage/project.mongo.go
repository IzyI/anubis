package storage

import (
	"anubis/app/api/DAL/entitiesDB"
	"anubis/app/core/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RepositoryMongoProjects struct {
	db *mongo.Client
}

func NewRepositoryMongoProject(db *mongo.Client) *RepositoryMongoProjects {
	return &RepositoryMongoProjects{db: db}
}

func (r *RepositoryMongoProjects) CreateProject(service string, project *entitiesDB.MdProject) (primitive.ObjectID, error) {
	collection := r.db.Database(service).Collection("projects")

	project.CreatedAt = time.Now() // Устанавливаем дату создания
	_, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return project.ID, nil
}

func (r *RepositoryMongoProjects) AddMemberToProject(service string, projectID primitive.ObjectID, userID string, role string) error {
	collection := r.db.Database(service).Collection("projects")

	update := bson.M{
		"$addToSet": bson.M{
			"members": entitiesDB.MdProjectMember{
				UserID:   userID,
				Role:     role,
				JoinedAt: time.Now(),
			},
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": projectID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryMongoProjects) GetProjectsByUser(service string, userID string) ([]entitiesDB.MdProject, error) {
	collection := r.db.Database(service).Collection("projects")

	filter := bson.M{"members.userId": userID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			common.LogInfo("Error cursor 8h08:")
		}
	}(cursor, context.TODO())

	var projects []entitiesDB.MdProject
	for cursor.Next(context.TODO()) {
		var project entitiesDB.MdProject
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

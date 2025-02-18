package storage

import (
	"anubis/app/api/DAL/entitiesDB"
	"anubis/app/core/common"
	"anubis/tools/utils"
	"context"
	"fmt"
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

	_, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return project.ID, nil
}

func (r *RepositoryMongoProjects) AddMemberToProject(service string, projectID primitive.ObjectID, userID primitive.ObjectID, role string) error {
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

func (r *RepositoryMongoProjects) GetProjectsListByUser(service string, domain string, userID primitive.ObjectID) (map[string]string, error) {
	collection := r.db.Database(service).Collection("projects")

	// Получаем все записи Project для конкретного пользователя
	filter := bson.M{"members.userId": userID}
	var userProjects []entitiesDB.MdProject

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			// Логирование ошибки
			common.LogInfo("Error closing cursor 7gGhd")
		}
	}(cursor, context.TODO())

	// Читаем результаты из курсора
	if err = cursor.All(context.TODO(), &userProjects); err != nil {
		return nil, err
	}

	// Если нет ни одной записи, создаем дефолтный проект

	if len(userProjects) == 0 {
		var userProject entitiesDB.MdProject
		userProject.Name = "PRJ-" + service + "_" + utils.RandStringBytes(7)
		userProject.Domain = domain
		userProject.Members = []entitiesDB.MdProjectMember{
			{
				UserID:   userID,
				Role:     "O", // Дефолтная роль
				JoinedAt: time.Now(),
			},
		}

		res, err1 := collection.InsertOne(context.TODO(), userProject)
		if err1 != nil {
			return nil, err1
		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			userProject.ID = oid
		} else {
			return nil, fmt.Errorf("unexpected type of InsertedID in GetProjectsListByUser")
		}

		userProjects = append(userProjects, userProject) // Добавляем дефолтный проект в срез
	}

	// Формируем карту для возврата с именами проектов и их ID
	result := make(map[string]string)
	for _, project := range userProjects {
		result[project.Name] = project.ID.Hex()
	}

	return result, nil
}

func (r *RepositoryMongoProjects) GetProjectsByUser(service string, project *entitiesDB.MdProject, userId primitive.ObjectID) error {
	collection := r.db.Database(service).Collection("projects")
	filter := bson.M{"members.userId": userId, "_id": project.ID}
	err := collection.FindOne(context.Background(), filter).Decode(project)
	return err

}

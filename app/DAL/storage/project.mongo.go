package storage

import (
	"anubis/app/DAL/entitiesDB"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/tools/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RepositoryMongoProjects struct {
	db     *mongo.Client
	config *core.ServiceConfig
}

func NewRepositoryMongoProject(db *mongo.Client, config *core.ServiceConfig) *RepositoryMongoProjects {
	return &RepositoryMongoProjects{db: db, config: config}
}

func (r *RepositoryMongoProjects) CreateProject(service string, project *entitiesDB.MdProject) error {
	collection := r.db.Database(service).Collection("projects")

	res, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		project.ID = oid
	} else {
		return fmt.Errorf("unexpected type of InsertedID in CreateProject")
	}

	return err
}

func (r *RepositoryMongoProjects) UpdateProjectName(service string, domain string, project *entitiesDB.MdProject) error {
	collection := r.db.Database(service).Collection("projects")
	filter := bson.M{}
	filter["_id"] = project.ID
	memberFilter := bson.M{
		"userId": project.Members[0].UserID,
		"role":   r.config.ListServices[domain].Role["owner"],
	}
	filter["members"] = bson.M{"$elemMatch": memberFilter}
	update := bson.M{
		"$set": bson.M{
			"name": project.Name,
		},
	}
	opts := options.Update()
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)

	return err
}

func (r *RepositoryMongoProjects) DelProjectID(service string, domain string, project *entitiesDB.MdProject) error {
	collection := r.db.Database(service).Collection("projects")

	filter := bson.M{}
	filter["_id"] = project.ID
	filter["domain"] = project.Domain
	memberFilter := bson.M{
		"userId": project.Members[0].UserID,
		"role":   r.config.ListServices[domain].Role["owner"],
	}
	filter["members"] = bson.M{"$elemMatch": memberFilter}
	deleteResult, err := collection.DeleteOne(context.Background(), filter)

	if deleteResult.DeletedCount == 0 {
		return errors.New("Not found project")
	}
	return err
}

func (r *RepositoryMongoProjects) AddMemberToProject(
	service string,
	domain string,
	projectID primitive.ObjectID,
	ownerID primitive.ObjectID,
	userID primitive.ObjectID,
	role string,
) error {
	collection := r.db.Database(service).Collection("projects")

	ctx := context.TODO()

	// Фильтр: проверяем, что существует владелец с ролью "O" и что userID ещё нет в проекте
	filter := bson.M{
		"_id": projectID,
		"members": bson.M{
			"$elemMatch": bson.M{
				"userId": ownerID,
				"role":   r.config.ListServices[domain].Role["owner"],
			},
		},
		"members.userId": bson.M{
			"$ne": userID, // Проверяем, что пользователя ещё нет в проекте
		},
	}

	// Обновление: добавляем нового участника, если фильтр совпал
	update := bson.M{
		"$addToSet": bson.M{
			"members": entitiesDB.MdProjectMember{
				UserID:   userID,
				Role:     role,
				JoinedAt: time.Now(),
			},
		},
	}

	// Выполняем один запрос с проверками и обновлением
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	// Если документ не был найден, значит либо нет owner'а с ролью "O", либо пользователь уже в проекте
	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to add member: either owner not found or user already in project")
	}

	return nil
}

func (r *RepositoryMongoProjects) UpdateMemberRole(
	service string,
	domain string,
	projectID primitive.ObjectID,
	ownerID primitive.ObjectID,
	userID primitive.ObjectID,
	newRole string,
) error {
	collection := r.db.Database(service).Collection("projects")

	ctx := context.TODO()

	// Фильтр: проверяем, что существует владелец с ролью "O" и что userID уже есть в проекте
	filter := bson.M{
		"_id": projectID,
		"members": bson.M{
			"$elemMatch": bson.M{
				"userId": ownerID,
				"role":   r.config.ListServices[domain].Role["owner"],
			},
		},
		"members.userId": userID, // Проверяем, что пользователь уже в проекте
	}

	// Обновление: изменяем роль у участника
	update := bson.M{
		"$set": bson.M{
			"members.$[elem].role": newRole, // Обновляем только нужного участника
		},
	}

	// Опции: указываем, что обновлять нужно только конкретного участника
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.userId": userID}, // Фильтр для поиска участника в массиве
		},
	})

	// Выполняем запрос
	result, err := collection.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	// Если документ не был найден, значит либо owner не найден, либо userID нет в проекте
	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to update role: either owner not found or user not in project")
	}

	return nil
}

func (r *RepositoryMongoProjects) RemoveMemberFromProject(
	service string,
	domain string,
	projectID primitive.ObjectID,
	ownerID primitive.ObjectID,
	userID primitive.ObjectID,
) error {
	collection := r.db.Database(service).Collection("projects")

	ctx := context.TODO()

	// Фильтр: проверяем, что существует владелец с ролью "O" и что userID уже в проекте
	filter := bson.M{
		"_id": projectID,
		"members": bson.M{
			"$elemMatch": bson.M{
				"userId": ownerID,
				"role":   r.config.ListServices[domain].Role["owner"],
			},
		},
		"members.userId": userID, // Проверяем, что пользователь есть в проекте
	}

	// Обновление: удаляем участника
	update := bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"userId": userID, // Удаляем по userID
			},
		},
	}

	// Выполняем запрос
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	// Если документ не был найден, значит либо owner не найден, либо userID нет в проекте
	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to remove member: either owner not found or user not in project")
	}

	return nil
}

//func (r *RepositoryMongoProjects) GetProjectsByUserOwnerID(service string,domain string, userID primitive.ObjectID) ([]entitiesDB.MdProject, error) {
//	collection := r.db.Database(service).Collection("projects")
//
//	filter := bson.M{
//		"members.userId": userID,
//		"members.role":   r.config.ListServices[domain].Role["owner"],
//	}
//
//	cursor, err := collection.Find(context.TODO(), filter)
//	if err != nil {
//		return nil, fmt.Errorf("failed to find projects: %w", err)
//	}
//	defer cursor.Close(context.TODO())
//
//	var projects []entitiesDB.MdProject
//	if err = cursor.All(context.TODO(), &projects); err != nil {
//		return nil, fmt.Errorf("failed to decode projects: %w", err)
//	}
//
//	return projects, nil
//}

func (r *RepositoryMongoProjects) GetProjectByUserOwnerID(
	service string,
	domain string,
	project *entitiesDB.MdProject, userId primitive.ObjectID) error {
	collection := r.db.Database(service).Collection("projects")

	filter := bson.M{
		"_id":            project.ID,
		"members.userId": userId,
		"members.role":   r.config.ListServices[domain].Role["owner"],
	}

	err := collection.FindOne(context.Background(), filter).Decode(project)
	return err
}
func (r *RepositoryMongoProjects) GetProjectIDByUserID(service string, project *entitiesDB.MdProject, userId primitive.ObjectID) error {
	collection := r.db.Database(service).Collection("projects")

	filter := bson.M{
		"_id":            project.ID,
		"members.userId": userId,
	}

	err := collection.FindOne(context.Background(), filter).Decode(project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("project not found for user %s", userId)
		}
		return fmt.Errorf("failed to find project: %w", err)
	}
	return err
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
				Role:     r.config.ListServices[domain].Role["owner"],
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

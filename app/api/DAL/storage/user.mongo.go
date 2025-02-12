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
)

type RepositoryMongoUser struct {
	db *mongo.Client
}

func NewRepositoryMongoUser(db *mongo.Client) *RepositoryMongoUser {
	return &RepositoryMongoUser{db: db}
}

//func (r *RepositoryMongoUser) CreateUser() (*entitiesDB.MdUser, error) {
//	var user entitiesDB.MdUser
//	sql := `INSERT INTO users DEFAULT VALUES RETURNING uuid`
//	rows := r.db.QueryRow(context.Background(), sql)
//
//	err := rows.Scan(&user.Uuid)
//	if err != nil {
//		return &user, err
//	}
//
//	return &user, nil
//}

func (r *RepositoryMongoUser) CreateUser(service string, user *entitiesDB.MdUser) error {
	collection := r.db.Database(service).Collection("users")
	result, err := collection.InsertOne(context.TODO(), &user)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	} else {
		return fmt.Errorf("unexpected type of InsertedID")
	}

	return nil
}

func (r *RepositoryMongoUser) GetGroupUser(service string, domain string, userID primitive.ObjectID) (map[string]string, error) {
	collection := r.db.Database(service).Collection("users_group")

	// Получаем все записи MdUsersGroup для конкретного юзера
	filter := bson.M{"user_id": userID}
	var userGroups []entitiesDB.MdUsersGroup
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			common.LogInfo("Error cursor 8h78:")
		}
	}(cursor, context.TODO())

	// Читаем результаты из курсора
	if err = cursor.All(context.TODO(), &userGroups); err != nil {
		return nil, err
	}
	var userGroup entitiesDB.MdUsersGroup
	// Если нет ни одной записи, создаем дефолтную
	if len(userGroups) == 0 {
		userGroup.Domain = domain // Генерируем новый ID
		userGroup.UserID = userID
		userGroup.Owner = true
		userGroup.Role = domain + "BaseRole"
		userGroup.GroupName = domain + "_" + utils.RandStringBytes(7)
		res, err1 := collection.InsertOne(context.TODO(), userGroup)
		if err1 != nil {
			return nil, err1
		}
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			userGroup.ID = oid
		} else {
			return nil, fmt.Errorf("unexpected type of InsertedID")
		}

		userGroups = append(userGroups, userGroup) // Добавляем дефолтную группу в срез
	}

	// Формируем карту для возврата
	result := make(map[string]string)
	for _, group := range userGroups {
		result[group.GroupName] = group.ID.Hex()
	}

	return result, nil
}

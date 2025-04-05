package storage

import (
	"anubis/app/DAL/entitiesDB"
	"anubis/app/core/common"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"time"
)

type RepositoryMongoUser struct {
	db *mongo.Client
}

func NewRepositoryMongoUser(db *mongo.Client) *RepositoryMongoUser {
	return &RepositoryMongoUser{db: db}
}

func (r *RepositoryMongoUser) CreateUser(service string, user *entitiesDB.MdUser) error {
	collection := r.db.Database(service).Collection("users")
	result, err := collection.InsertOne(context.TODO(), &user)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	} else {
		return fmt.Errorf("unexpected type of InsertedID in CreateUser")
	}

	return nil
}

func (r *RepositoryMongoUser) GetUserByID(service string, user *entitiesDB.MdUser) error {
	collection := r.db.Database(service).Collection("users")
	filter := bson.M{"_id": user.ID}
	err := collection.FindOne(context.TODO(), filter).Decode(user)
	return err
}

func (r *RepositoryMongoUser) GetUsersByIDs(service string, userIDs []primitive.ObjectID) (map[primitive.ObjectID]entitiesDB.MdUser, error) {
	collection := r.db.Database(service).Collection("users")

	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			// Логирование ошибки
			common.LogInfo("Error closing cursor 7tg")
		}
	}(cursor, context.TODO())

	var users []entitiesDB.MdUser
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	usersMap := make(map[primitive.ObjectID]entitiesDB.MdUser, len(users))
	for _, user := range users {
		usersMap[user.ID] = user
	}

	return usersMap, nil
}

func (r *RepositoryMongoUser) InsertUsersSession(service string, userSession *entitiesDB.MdUsersSession) error {
	collection := r.db.Database(service).Collection("users_session")
	result, err := collection.InsertOne(context.TODO(), userSession)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		userSession.ID = oid

	} else {
		return fmt.Errorf("unexpected type of InsertedID in CreateSession")
	}
	return nil
}

func (r *RepositoryMongoUser) CheckOldAndUpdateSession(
	service string,
	userSession *entitiesDB.MdUsersSession,
) error {
	collection := r.db.Database(service).Collection("users_session")
	ctx := context.TODO()

	// 1. Находим ID самой старой сессии и список ID для деактивации
	filter := bson.M{
		"user_id":   userSession.UserID,
		"domain":    userSession.Domain,
		"is_active": true,
	}

	// Получаем все активные сессии
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			// Логирование ошибки
			common.LogInfo("Error closing cursor 237h")
		}
	}(cursor, context.TODO())

	var sessions []entitiesDB.MdUsersSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return err
	}

	// Сортируем сессии от старых к новым (по возрастанию CreatedAt)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.Before(sessions[j].CreatedAt)
	})

	// 3. Формируем список ID для деактивации
	var idsToDeactivate []primitive.ObjectID
	now := time.Now()

	// Добавляем все сессии, кроме последних 3
	if len(sessions) > 3 {
		// Берём все сессии, кроме последних 3
		sessionsToDeactivate := sessions[:len(sessions)-3]
		for _, s := range sessionsToDeactivate {
			idsToDeactivate = append(idsToDeactivate, s.ID)
		}
	}
	// Добавляем сессии с тем же device_id (предыдущая логика)
	for _, s := range sessions {
		if s.DeviceId == userSession.DeviceId {
			idsToDeactivate = append(idsToDeactivate, s.ID)
		}
	}

	// Удаляем дубликаты (если ID уже добавлены)
	uniqueIDs := make(map[primitive.ObjectID]struct{})
	for _, id := range idsToDeactivate {
		uniqueIDs[id] = struct{}{}
	}
	// 4. Выполняем единое обновление
	if len(uniqueIDs) > 0 {
		ids := make([]primitive.ObjectID, 0, len(uniqueIDs))
		for id := range uniqueIDs {
			ids = append(ids, id)
		}

		_, err = collection.UpdateMany(ctx,
			bson.M{"_id": bson.M{"$in": ids}},
			bson.M{"$set": bson.M{
				"is_active":  false,
				"expires_at": now,
			}},
		)
		if err != nil {
			return err
		}
	}

	if userSession.ID == primitive.NilObjectID {
		return r.InsertUsersSession(service, userSession)
	} else {
		return r.UpdateSessionsByID(service, userSession)

	}

}

func (r *RepositoryMongoUser) DeactivateUserSessionsByDomain(
	service string,
	userID primitive.ObjectID,
	domain string,
) error {
	collection := r.db.Database(service).Collection("users_session")
	// Определяем фильтр для поиска сессий
	filter := bson.M{
		"user_id": userID,
		"domain":  domain,
	}

	// Определяем обновление для установки is_active в true
	update := bson.M{
		"$set": bson.M{"is_active": false},
	}

	// Выполняем обновление для всех документов, которые соответствуют фильтру
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryMongoUser) DeactivateUserALLSessions(
	service string, userSession *entitiesDB.MdUsersSession, flagDomain bool) error {
	collection := r.db.Database(service).Collection("users_session")
	filter := bson.M{
		"user_id": userSession.UserID,
	}

	// Определяем обновление для установки is_active в true
	update := bson.M{
		"$set": bson.M{"is_active": false},
	}

	// Выполняем обновление для всех документов, которые соответствуют фильтру
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (r *RepositoryMongoUser) UpdateSessionsByID(service string, userSession *entitiesDB.MdUsersSession) error {
	collection := r.db.Database(service).Collection("users_session")

	// Создаем карту для оператора $set
	updateFields := bson.M{
		"user_id":     userSession.UserID,
		"domain":      userSession.Domain,
		"device_id":   userSession.DeviceId,
		"device_type": userSession.DeviceType,
		"created_at":  userSession.CreatedAt,
		"expires_at":  userSession.ExpiresAt,
		"ip":          userSession.IP,
		"is_active":   userSession.IsActive,
		"hash_token":  userSession.HashToken,
	}

	filter := bson.M{"_id": userSession.ID}
	update := bson.M{"$set": updateFields}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *RepositoryMongoUser) DeleteSessionsByID(service string, userSession *entitiesDB.MdUsersSession) error {
	collection := r.db.Database(service).Collection("users_session")

	filter := bson.M{"_id": userSession.ID}
	collection.FindOneAndDelete(context.TODO(), filter)
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}

func (r *RepositoryMongoUser) DeactivateOldAndCreateSession(
	service string,
	userSession *entitiesDB.MdUsersSession,
) error {
	err := r.DeactivateUserSessionsByDomain(service, userSession.UserID, userSession.Domain)
	if err != nil {
		return err
	}
	err = r.InsertUsersSession(service, userSession)

	return err
}

func (r *RepositoryMongoUser) GetUsersSessionByID(
	service string,
	id primitive.ObjectID,
	userSession *entitiesDB.MdUsersSession,
) error {

	collection := r.db.Database(service).Collection("users_session")

	// Create a filter to find the document by _id
	filter := bson.M{"_id": id}

	// Find the document
	err := collection.FindOne(context.TODO(), filter).Decode(userSession)

	return err
}

func (r *RepositoryMongoUser) DeactivateUserSessionsByTokenFamily(
	service string, idTokenFamily primitive.ObjectID) error {

	collection := r.db.Database(service).Collection("token_families")

	// Помечаем семейство как скомпрометированное
	_, err := collection.UpdateOne(context.TODO(),
		bson.M{"_id": idTokenFamily},
		bson.M{"$set": bson.M{
			"is_compromised": true,
			"compromised_at": time.Now(),
		}},
	)
	if err != nil {
		return err
	}
	collection = r.db.Database(service).Collection("users_session")
	_, err = collection.UpdateMany(context.TODO(),
		bson.M{"family_id": idTokenFamily},
		bson.M{"$set": bson.M{"is_active": false}},
	)
	return err
}

func (r *RepositoryMongoUser) UserSessionsSetActive(service string, userSessionID primitive.ObjectID, active bool) error {
	collection := r.db.Database(service).Collection("users_session")
	// Определяем фильтр для поиска сессий
	filter := bson.M{
		"_id": userSessionID,
	}

	// Определяем обновление для установки is_active в true
	update := bson.M{
		"$set": bson.M{"is_active": active},
	}

	// Выполняем обновление для всех документов, которые соответствуют фильтру
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

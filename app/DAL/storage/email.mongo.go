package storage

import (
	"anubis/app/DAL/entitiesDB"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoryMongoAuthEmail struct {
	db *mongo.Client
}

func NewRepositoryMongoAuthEmail(db *mongo.Client) *RepositoryMongoAuthEmail {
	return &RepositoryMongoAuthEmail{db: db}
}

func (r *RepositoryMongoAuthEmail) GetEmail(service string, email *entitiesDB.MdEmailAuth) error {
	err := r.db.Database(service).Collection("email_auth").FindOne(context.TODO(), bson.M{"email": email.Email}).Decode(&email)

	return err
}

func (r *RepositoryMongoAuthEmail) SaveEmail(service string, email *entitiesDB.MdEmailAuth) error {
	collection := r.db.Database(service).Collection("email_auth")

	filter := bson.M{"email": email.Email}
	update := bson.M{
		"$set": bson.M{
			"user_id": email.UserID,
		},
		"$setOnInsert": bson.M{
			"verification":  email.Verification,
			"password_hash": email.PasswordHash,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&email)
	return err
}

func (r *RepositoryMongoAuthEmail) SaveEmailAuth(service string, emailCode *entitiesDB.MdEmailCodeAuth) error {
	collection := r.db.Database(service).Collection("email_code_auth")
	result, err := collection.InsertOne(context.TODO(), emailCode)
	if err != nil {
		return err
	}

	// Извлекаем сгенерированный ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		emailCode.ID = oid
	} else {
		return fmt.Errorf("unexpected type of InsertedID in SaveSmsAuth")
	}
	return nil
}

func (r *RepositoryMongoAuthEmail) EmailCodeValidUser(service string, ID primitive.ObjectID, emailCode string) (string, error) {
	collection := r.db.Database(service).Collection("email_code_auth")

	filter := bson.M{
		"_id":        ID,
		"email_code": emailCode,
	}

	// Ищем документ
	var result entitiesDB.MdEmailCodeAuth
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		// Если документ не найден или произошла ошибка
		return "", err
	}

	// Возвращаем поле Phone из найденного документа
	return result.Email, nil
}

func (r *RepositoryMongoAuthEmail) SaveVerifyEmail(service string, verify bool, email string) error {
	collection := r.db.Database(service).Collection("email_auth")

	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"verification": verify,
		},
	}
	opts := options.Update()
	result, err := collection.UpdateOne(context.Background(), filter, update, opts)

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with email: %s", email)
	}
	return err
}

func (r *RepositoryMongoAuthEmail) GetEmailVerificationUserID(service string, email string, verification bool) (primitive.ObjectID, string, error) {
	// Определение коллекции
	collection := r.db.Database(service).Collection("email_auth")

	// Определение фильтра для поиска пользователя по номеру телефона
	filter := bson.M{"email": email, "verification": verification}

	// Переменная для хранения найденного документа
	var user entitiesDB.MdPhoneAuth

	// Попытка найти пользователя
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return primitive.NilObjectID, "", err
	}

	// Возвращаем UserID и passwordHash
	return user.UserID, user.PasswordHash, nil
}

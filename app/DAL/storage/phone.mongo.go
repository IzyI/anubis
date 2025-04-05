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

type RepositoryMongoAuthPhone struct {
	db *mongo.Client
}

func NewRepositoryMongoAuthPhone(db *mongo.Client) *RepositoryMongoAuthPhone {
	return &RepositoryMongoAuthPhone{db: db}
}

func (r *RepositoryMongoAuthPhone) GetPhone(service string, phone *entitiesDB.MdPhoneAuth) error {
	err := r.db.Database(service).Collection("phone_auth").FindOne(context.TODO(), bson.M{"phone": phone.Phone}).Decode(&phone)

	return err
}

func (r *RepositoryMongoAuthPhone) SavePhone(service string, phone *entitiesDB.MdPhoneAuth) error {
	collection := r.db.Database(service).Collection("phone_auth")

	filter := bson.M{"phone": phone.Phone}
	update := bson.M{
		"$set": bson.M{
			"user_id": phone.UserID,
		},
		"$setOnInsert": bson.M{
			"verification":  phone.Verification,
			"country_code":  phone.CountryCode,
			"password_hash": phone.PasswordHash,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&phone)
	return err
}

func (r *RepositoryMongoAuthPhone) SaveSmsAuth(service string, sms *entitiesDB.MdSmsAuth) error {
	collection := r.db.Database(service).Collection("sms_auth")
	result, err := collection.InsertOne(context.TODO(), sms)
	if err != nil {
		return err
	}

	// Извлекаем сгенерированный ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		sms.ID = oid
	} else {
		return fmt.Errorf("unexpected type of InsertedID in SaveSmsAuth")
	}
	return nil
}

func (r *RepositoryMongoAuthPhone) SmsValidUser(service string, ID primitive.ObjectID, smsCode string) (int64, error) {
	collection := r.db.Database(service).Collection("sms_auth")

	filter := bson.M{
		"_id":      ID,
		"sms_code": smsCode,
	}

	// Ищем документ
	var result entitiesDB.MdSmsAuth
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		// Если документ не найден или произошла ошибка
		return 0, err
	}

	// Возвращаем поле Phone из найденного документа
	return result.Phone, nil
}

func (r *RepositoryMongoAuthPhone) SaveVerifyPhone(service string, verify bool, phone int64) error {
	collection := r.db.Database(service).Collection("phone_auth")

	filter := bson.M{"phone": phone}
	update := bson.M{
		"$set": bson.M{
			"verification": verify,
		},
	}
	opts := options.Update()
	result, err := collection.UpdateOne(context.Background(), filter, update, opts)

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with phone: %d", phone)
	}
	return err

}

func (r *RepositoryMongoAuthPhone) GetPhoneVerificationUserID(service string, phone int64, verification bool) (primitive.ObjectID, string, error) {
	// Определение коллекции
	collection := r.db.Database(service).Collection("phone_auth")

	// Определение фильтра для поиска пользователя по номеру телефона
	filter := bson.M{"phone": phone, "verification": verification}

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

//func (r *RepositoryMongoAuthPhone) GetUuidUser(uuid uuid.UUID) error {
//	//var p string
//	//sql := `SELECT phone FROM users   WHERE uuid=($1) and verification=true`
//	//rows := r.db.QueryRow(context.Background(), sql, uuid)
//	//
//	//err := rows.Scan(&p)
//	//if err != nil {
//	//	return err
//	//}
//
//	return nil
//}

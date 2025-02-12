package dtb

import (
	"anubis/tools/dtb"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func NewClientMongo(
	ctx context.Context,
	maxAttempts int,
	maxDelay time.Duration,
	databaseURI string,
	maxpool int,
) (client *mongo.Client, err error) {

	clientOptions := options.Client().ApplyURI(databaseURI)
	clientOptions.SetMaxPoolSize(uint64(maxpool))
	// Настраиваем таймаут подключения
	clientOptions.SetConnectTimeout(60 * time.Second)

	var connectErr error
	client, connectErr = mongo.Connect(ctx, clientOptions)
	if connectErr != nil {
		log.Printf("failed to create mongo client: %w", connectErr)
		return nil, connectErr
	}

	// Функция для повторных попыток подключения
	err = dtb.DoWithAttempts(func() error {
		pingErr := client.Ping(ctx, readpref.Primary())
		if pingErr != nil {
			log.Printf("failed to ping mongo: %v", pingErr)
			return pingErr
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		log.Fatal("All attempts exceeded. Failed to connect to MongoDB")
	}

	return client, nil
}

func CloseMongoDBConnection(client *mongo.Client) {

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

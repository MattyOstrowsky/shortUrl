package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(host string, port int, username string, password string) (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB server: %v", err)
	}

	return client, nil
}

func GetValue(client *mongo.Client, collectionName string, database string, key string) (string, error) {
	collection := client.Database(database).Collection(collectionName)

	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{"url": key}).Decode(&result)
	if err != nil {
		return "", err
	}

	value := result["shortUrl"].(string)
	return value, nil
}

func Set(client *mongo.Client, collectionName string, database string, key string, value string) error {
	collection := client.Database(database).Collection(collectionName)

	_, err := collection.InsertOne(context.Background(), bson.M{"url": key, "shortUrl": value})
	if err != nil {
		return fmt.Errorf("failed to set value for key %s: %v", key, err)
	}

	return nil
}

func Update(client *mongo.Client, collectionName string, database string, key string, value string) error {
	collection := client.Database(database).Collection(collectionName)

	_, err := collection.UpdateOne(context.Background(), bson.M{"url": key}, bson.M{"$set": bson.M{"shortUrl": value}})
	if err != nil {
		return fmt.Errorf("failed to update value for key %s: %v", key, err)
	}

	return nil
}

func Delete(client *mongo.Client, collectionName string, database string, key string) error {
	collection := client.Database(database).Collection(collectionName)

	_, err := collection.DeleteOne(context.Background(), bson.M{"url": key})
	if err != nil {
		return fmt.Errorf("failed to delete entry for key %s: %v", key, err)
	}

	return nil
}

func GetAll(client *mongo.Client, collectionName string, database string) (map[string]string, error) {
	collection := client.Database(database).Collection(collectionName)

	var results []bson.M
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all entries: %v", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %v", err)
	}

	resultMap := make(map[string]string)
	for _, result := range results {
		key := result["url"].(string)
		value := result["shortUrl"].(string)
		resultMap[key] = value
	}

	return resultMap, nil
}

func GetKey(client *mongo.Client, collectionName string, database string, value string) (string, error) {
	collection := client.Database(database).Collection(collectionName)

	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{"shortUrl": value}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to get key for value %s: %v", value, err)
	}

	key := result["url"].(string)
	return key, nil
}

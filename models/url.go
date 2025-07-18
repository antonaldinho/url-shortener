package models

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"url-shortener/db"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type URL struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ShortURL    string        `json:"shortUrl" bson:"shortUrl"`
	OriginalURL string        `json:"originalUrl" bson:"originalUrl"`
}

func GetURL(id string) (*URL, error) {
	fmt.Println("id: " + id)
	coll := db.GetInstance().Client.Database("db").Collection("urls")
	var result *URL
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID")
	}

	filter := bson.D{{"_id", objectID}}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("product not found")
	}
	return result, nil
}

func CreateShortURL(url string) (*URL, error) {
	shortUrl, err := generateShortURL()
	if err != nil {
		return nil, err
	}
	coll := db.GetInstance().Client.Database("db").Collection("urls")
	result, err := coll.InsertOne(context.TODO(), URL{
		ShortURL:    shortUrl,
		OriginalURL: url,
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf(err.Error())
	}

	return &URL{ID: result.InsertedID.(bson.ObjectID), ShortURL: shortUrl, OriginalURL: url}, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func generateShortURL() (string, error) {
	// This function would contain logic to shorten the URL
	randomString, err := generateRandomString(6)
	if err != nil {
		return "", err
	}
	return ("http://short.url/" + randomString), nil
}

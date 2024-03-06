package database

import (
	"context"
	"encoding/base64"
	"fmt"
	"service/app/models"
	"service/app/queries"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type db struct {
	service  *mongo.Collection
	states *mongo.Collection
}

func (d *db) CreateToken(ctx context.Context, guid string) (string, error) {
	const msgLog = "platform.database.mongodb.createtoken"
	
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()+guid), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}

	request := models.ResponseDB{
		GUID:         guid,
		RefreshToken: refreshToken, 
	}

    _, err = d.service.InsertOne(ctx, request) 
    if err != nil {
        return "", fmt.Errorf("%s: %w", msgLog, err)
    }
	refreshTokenEncoded := base64.StdEncoding.EncodeToString(refreshToken)
    return refreshTokenEncoded, err
}

func(d *db) UpdateToken(ctx context.Context, guid string) (string, error) {
	const msgLog = "platform.database.mongodb.updatetoken"
	filter := bson.D{{"guid", guid}}
	// add check filter is null	
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()+guid), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}

	update := bson.M{
		"refresh": refreshToken,
	}
	_, err = d.service.ReplaceOne(ctx, filter, update)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}

	refreshTokenEncoded := base64.StdEncoding.EncodeToString(refreshToken)
	return  refreshTokenEncoded, nil
}
	
func (d *db) SearchTokenByRefresh(ctx context.Context, refreshToken string) (result string, err error) {
	const msgLog = "platform.database.mongodb.searchtokenbyrefresh"
	filter := bson.D{{"refresh", refreshToken}}
	// add check filter is null
	err = d.service.FindOne(ctx, filter).Decode(&result)
	if err != nil {
        return "", fmt.Errorf("%s: %w", msgLog, err)
	}
	return result, err
}

func (d *db) SearchTokenByGuid(ctx context.Context, guid string) (result string, err error) {
	const msgLog = "platform.database.mongodb.searchtokenbyguid"
	filter := bson.D{{"guid", guid}}
	var answer models.ResponseDB
	// add check filter is null	
	err = d.service.FindOne(ctx, filter).Decode(&answer)
	if err != nil {
        return "", fmt.Errorf("%s: %w", msgLog, err)
	}
	result = base64.StdEncoding.EncodeToString(answer.RefreshToken)
	return result, err
}

func (d *db) DeleteToken(ctx context.Context, access string) error {
	return nil
}




func NewStorage(database *mongo.Database, collectionName string) queries.Storage {
	return &db{
		service: database.Collection(collectionName),
	}
}

func strToObjectId(strIdsRoom []string) []primitive.ObjectID {
	var objectIds []primitive.ObjectID
	for _, strId := range strIdsRoom {
		objectId, _ := primitive.ObjectIDFromHex(strId)
		objectIds = append(objectIds, objectId)
	}
	return objectIds
}
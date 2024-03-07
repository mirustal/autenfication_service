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

func (d *db) CreateRefreshToken(ctx context.Context, guid string) (string, error) {
	const msgLog = "platform.database.mongodb.createtoken"
	randomToken := time.Now().UTC().GoString() + " " + guid
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(randomToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}
	

	request := models.ResponseDB{
		GUID:         guid,
		RefreshToken: string(refreshToken), 
	}
	_, err = d.service.InsertOne(ctx, request)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}

	refreshTokenEncoded := base64.StdEncoding.EncodeToString([]byte(randomToken))
	return refreshTokenEncoded, nil
}



func(d *db) UpdateRefreshToken(ctx context.Context, guid string) (string, error) {
	const msgLog = "platform.database.mongodb.updatetoken"
	filter := bson.D{{"guid", guid}}
	// add check filter is null	
	randomToken := time.Now().String() + guid
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(randomToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}

	update := bson.M{
		"guid": guid,
		"refreshtoken": string(refreshToken),
	}
	_, err = d.service.ReplaceOne(ctx, filter, update)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msgLog, err)
	}

	refreshTokenEncoded := base64.StdEncoding.EncodeToString([]byte(randomToken))
	return  refreshTokenEncoded, nil
}
	
func (d *db) SearchTokenByRefresh(ctx context.Context, refreshToken string) (result string, err error) {
	const msgLog = "platform.database.mongodb.searchtokenbyrefresh"
	filter := bson.D{{"refresh", refreshToken}}

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
	result = answer.RefreshToken
	return result, err
}

func (d *db) DeleteRefreshToken(ctx context.Context, access string) error {
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


func (d *db) ValidateRefreshToken(ctx context.Context, guid, token string) (bool, error) {

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, fmt.Errorf("error decoding token: %w", err)
	}


	hashedToken, err := d.SearchTokenByGuid(ctx, guid)
	if err != nil {
		return false, fmt.Errorf("error finding token by GUID: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), decodedToken)
	if err != nil {
		return false, nil
	}

	return true, nil
}
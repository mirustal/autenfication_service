package database

import (
	"context"
	"fmt"
	"service/app/models"
	"service/app/queries"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	service  *mongo.Collection
	states *mongo.Collection
}

func (d *db) CreateToken(ctx context.Context) (models.AccessResponse, error) {
	const msgLog = "platform.database.mongodb.createtoken"
    var newAccessToken, newRefreshToken string
     request :=  models.AccessResponse{
        Access: newAccessToken,    
        Refresh: newRefreshToken, 
    }

    _, err := d.service.InsertOne(ctx, request) 
    if err != nil {
        return models.AccessResponse{}, fmt.Errorf("%s: %w", msgLog, err)
    }
	
    return request, err
}

func(d *db) UpdateToken(ctx context.Context, refreshToken string) (models.AccessResponse, error) {
	const msgLog = "platform.database.mongodb.updatetoken"
	result, err := d.SearchTokenByRefresh(ctx, refreshToken)
	if err != nil {
        return models.AccessResponse{}, fmt.Errorf("%s: %w", msgLog, err)
	}

	err = d.DeleteToken(ctx, result.Access)
	if err != nil {
        return models.AccessResponse{}, fmt.Errorf("%s: %w", msgLog, err)
	}

	result, err = d.CreateToken(ctx)
	if err != nil {
        return models.AccessResponse{}, fmt.Errorf("%s: %w", msgLog, err)
	}

	return result, err
}
	
func (d *db) SearchTokenByRefresh(ctx context.Context, refreshToken string) (result models.AccessResponse, err error) {
	const msgLog = "platform.database.mongodb.findtoken"
	filter := bson.D{{"_refresh", refreshToken}}
	// add check filter is null
	err = d.service.FindOne(ctx, filter).Decode(&result)
	if err != nil {
        return models.AccessResponse{}, fmt.Errorf("%s: %w", msgLog, err)
	}
	return result, err
}

func (d *db) DeleteToken(ctx context.Context, access string) error {
	const msgLog = "platform.database.mongodb.deletetoken"
	filter := bson.D{{"_access", access}}
	_, err := d.service.DeleteOne(ctx, filter)
	if err != nil {
        return fmt.Errorf("%s: %w", msgLog, err)
	}
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
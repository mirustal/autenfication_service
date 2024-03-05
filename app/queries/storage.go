package queries

import (
	"service/app/models"
	"context"
)

type Storage interface {
	CreateToken(context.Context) (models.AccessResponse, error)
	UpdateToken(context.Context, string) (models.AccessResponse, error)
	SearchTokenByRefresh(context.Context, string) (models.AccessResponse, error)
	DeleteToken(context.Context, string) (error)
}
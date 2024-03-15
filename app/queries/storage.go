package queries

import (
	"context"
)

type Storage interface {
	CreateRefreshToken(context.Context, string) (string, error)
	UpdateRefreshToken(context.Context, string) (string, error)
	SearchTokenByRefresh(context.Context, string) (string, error)
	SearchTokenByGuid(context.Context, string) (string, error)
	DeleteRefreshToken(context.Context, string) (error)
	ValidateRefreshToken(context.Context, string, string) (bool, error)
}
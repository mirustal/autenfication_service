package queries

import (
	"context"
)

type Storage interface {
	CreateToken(context.Context, string) (string, error)
	UpdateToken(context.Context, string) (string, error)
	SearchTokenByRefresh(context.Context, string) (string, error)
	SearchTokenByGuid(context.Context, string) (string, error)
	DeleteToken(context.Context, string) (error)
}
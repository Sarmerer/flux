package database

import (
	"github.com/flow/internal/domain/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ types.Executor = (*pgxpool.Pool)(nil)
var _ types.Executor = (pgx.Tx)(nil)

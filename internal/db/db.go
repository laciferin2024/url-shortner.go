package db

import (
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
)

type db struct {
	*genesis.Service
}

type PostgresDB struct {
	db
}

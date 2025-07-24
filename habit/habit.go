package habit

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	DB *sqlx.DB
}

// CreateUser demonstrates NamedExec with sqlx
func (app *Application) CreateUser(ctx context.Context, name string, email string) error {
	return nil
}

// [TODO] debug
func Foo() error {
	return nil
}

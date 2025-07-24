package habit

import (
	"github.com/kahnwong/habit-tracker/db"
	_ "github.com/mattn/go-sqlite3"
)

// [TODO] undo activity (today only)

func Foo() error {
	_ = db.Foo()

	return nil
}

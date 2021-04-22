package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	tools "github.com/imega/stock-miner/sql"
	_ "github.com/mattn/go-sqlite3"
)

type CloseDB func() error

func CreateDB(txFunc tools.TxFunc) (*sql.DB, CloseDB, error) {
	file, err := ioutil.TempFile("", "stockminer")
	if err != nil {
		return nil, nil, err
	}

	filename := file.Name()
	if err := file.Close(); err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	wrapper := tools.TxWrapper{db}
	wrapper.Transaction(ctx, nil, txFunc)

	return db,
		CloseDB(func() error {
			errDB := db.Close()
			if err := os.Remove(filename); err != nil || errDB != nil {
				return fmt.Errorf("%s, %s", errDB, err)
			}

			return nil
		}),
		nil
}

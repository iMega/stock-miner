package storage

import (
	"context"
	"database/sql"
	"os"

	tools "github.com/imega/stock-miner/sql"
)

type tableInfo struct {
	Rows []row
}

type row struct {
	CID          int
	Name         string
	Type         string
	NotNull      int
	DefaultValue string
	PK           int
}

func MigrateDatabase(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return err
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return err
	}

	var getTableInfo = func(name string) (tableInfo, error) {
		var ti tableInfo

		rows, err := db.Query("PRAGMA table_info(?)", name)
		defer rows.Close()
		if err != nil {
			return ti, err
		}

		for rows.Next() {
			r := row{}
			err := rows.Scan(&r.CID, &r.Name, &r.Type, &r.NotNull, &r.DefaultValue, &r.PK)
			if err != nil {
				return ti, err
			}

			ti.Rows = append(ti.Rows, r)
		}

		return ti, nil
	}

	wrapper := tools.TxWrapper{db}
	wrapper.Transaction(context.Background(), nil, func(ctx context.Context, tx *sql.Tx) error {
		slotInfo, err := getTableInfo("slot")
		if err != nil {
			return err
		}
		if err := slotTableMigrate(ctx, tx, slotInfo); err != nil {
			return err
		}

		return nil
	})

	return db.Close()
}

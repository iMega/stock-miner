package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	tools "github.com/imega/stock-miner/sql"
)

type tableInfo struct {
	Columns []col
}

type col struct {
	CID          int
	Name         string
	Type         string
	NotNull      int
	DefaultValue string
	PK           int
}

func MigrateDatabase(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return fmt.Errorf("failed getting file info, %w", err)
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return fmt.Errorf("failed to open database, %w", err)
	}

	getInfo := func(name string) (tableInfo, error) {
		return getTableInfo(db, name)
	}

	tx := func(ctx context.Context, tx *sql.Tx) error {
		slotInfo, err := getInfo("slot")
		if err != nil {
			return err
		}

		if err := slotTableMigrate(ctx, tx, slotInfo); err != nil {
			return err
		}

		slotInfo, err = getInfo("stock_item_approved")
		if err != nil {
			return err
		}

		return stockItemApprovedTableMigrate(ctx, tx, slotInfo)
	}

	wrapper := tools.TxWrapper{db}
	if err := wrapper.Transaction(context.Background(), nil, tx); err != nil {
		return fmt.Errorf("failed to execute transaction, %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database, %w", err)
	}

	return nil
}

func getTableInfo(db *sql.DB, name string) (tableInfo, error) {
	var ti tableInfo

	rows, err := db.Query("PRAGMA table_info(" + name + ")")
	if err != nil {
		return ti, fmt.Errorf("failed to execute query, %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dfl sql.NullString

		r := col{}

		err := rows.Scan(&r.CID, &r.Name, &r.Type, &r.NotNull, &dfl, &r.PK)
		if err != nil {
			return ti, fmt.Errorf("failed to scan row, %w", err)
		}

		if dfl.Valid {
			r.DefaultValue = dfl.String
		}

		ti.Columns = append(ti.Columns, r)
	}

	if err := rows.Err(); err != nil {
		return ti, fmt.Errorf("failed getting row, %w", err)
	}

	return ti, nil
}

func hasColumn(ti tableInfo, c col) bool {
	for _, v := range ti.Columns {
		if c.Name == v.Name {
			return true
		}
	}

	return false
}

func equalColumn(ti tableInfo, c col) bool {
	for _, v := range ti.Columns {
		if c.Name != v.Name {
			continue
		}

		if c.CID != v.CID {
			return false
		}

		if c.NotNull != v.NotNull {
			return false
		}

		if c.PK != v.PK {
			return false
		}

		if c.DefaultValue != v.DefaultValue {
			return false
		}

		if c.Type != v.Type {
			return false
		}

		return true
	}

	return false
}

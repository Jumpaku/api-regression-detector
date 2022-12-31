package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Jumpaku/api-regression-detector/log"
	"go.uber.org/multierr"
)

type Exec interface {
	Write(ctx context.Context, stmt string, params []any) (err error)
	Read(ctx context.Context, stmt string, params []any) (rows Rows, err error)
}

type exec struct {
	tx *sql.Tx
}

func rollback(ctx context.Context, tx *sql.Tx, err error) error {
	return multierr.Combine(err, tx.Rollback())
}
func commit(ctx context.Context, tx *sql.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		return rollback(ctx, tx, err)
	}
	return nil
}
func Transaction(ctx context.Context, db *sql.DB, handler func(ctx context.Context, exec Exec) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { err = rollback(ctx, tx, err) }()
	e := exec{tx: tx}
	err = handler(ctx, &e)
	if err != nil {
		return err
	}
	return commit(ctx, tx)
}

func (e *exec) Write(ctx context.Context, stmt string, params []any) (err error) {
	log.Stderr("SQL\nstatement: %v\nparams   : %v", stmt, params)
	_, err = e.tx.Exec(stmt, params...)
	if err != nil {
		return err
	}
	return nil
}

func (e *exec) Read(ctx context.Context, stmt string, params []any) (rows Rows, err error) {
	log.Stderr("SQL\nstatement: %v\nparams   : %v", stmt, params)
	itr, err := e.tx.Query(stmt, params...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = multierr.Combine(err, itr.Close())
	}()
	for itr.Next() {
		var columns []string
		columns, err = itr.Columns()
		if err != nil {
			return nil, err
		}
		types, err := itr.ColumnTypes()
		if err != nil {
			return nil, err
		}
		var values []any
		for range columns {
			var val any
			values = append(values, &val)
		}
		itr.Scan(values...)
		row := Row{}
		for i, column := range columns {
			row[column] = values[i]
		}
		for i, v := range columns {
			fmt.Printf(`%v:%v=%v:%v`+"\n", v, types[i].ScanType(), values[i], reflect.TypeOf(values[i]))
		}
		rows = append(rows, row)
	}
	return rows, nil
}

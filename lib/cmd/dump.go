package cmd

import (
	"context"
	"time"

	libdb "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func Dump(
	ctx context.Context,
	db libdb.DB,
	tableNames []string,
	schemaGetter SchemaGetter,
	rowLister RowLister,
) (tables.DumpTables, error) {
	dumpTables := tables.DumpTables{}

	err := db.RunTransaction(ctx, func(ctx context.Context, tx libdb.Tx) error {
		var err error

		dbTables := libdb.Tables{}
		for _, tableName := range tableNames {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return errors.Wrap(err, "fail to get schema of table %s", tableName)
			}
			rows, err := rowLister.ListRows(ctx, tx, tableName, schema)
			if err != nil {
				return errors.Wrap(err, "fail to list rows of table %s", tableName)
			}
			dbTables[tableName] = libdb.Table{Name: tableName, Schema: schema, Rows: rows}
		}

		dumpTables, err = convertTablesDBToJson(dbTables)
		if err != nil {
			return errors.Wrap(err, "fail to convert tables to JSON")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction for Dump failed")
	}

	return dumpTables, nil
}

func convertTablesDBToJson(dbTables libdb.Tables) (jsonTables tables.DumpTables, err error) {
	jsonTables = tables.DumpTables{}
	for dbTableName, dbTable := range dbTables {
		jsonRows := []tables.Row{}
		for _, dbRow := range dbTable.Rows {
			jsonRow := tables.Row{}
			for dbColumnName, dbColumnValue := range dbRow {
				jsonRow[dbColumnName], err = convertDBColumnValueToJsonValue(dbColumnValue)
				if err != nil {
					return nil, errors.Wrap(err, "fail to convert DB value to JSON value")
				}
			}

			jsonRows = append(jsonRows, jsonRow)
		}

		jsonTables[dbTableName] = jsonRows
	}

	return jsonTables, nil
}

func convertDBColumnValueToJsonValue(dbVal *libdb.ColumnValue) (*wrap.JsonValue, error) {
	switch dbVal.Type {
	case libdb.ColumnTypeBoolean:
		v, err := dbVal.AsBool()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to parse DB column value as bool value")
		}

		if !v.Valid {
			return wrap.Null(), nil
		}

		return wrap.Boolean(v.Bool), nil
	case libdb.ColumnTypeInteger:
		v, err := dbVal.AsInteger()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to parse DB column value as integer value")
		}

		if !v.Valid {
			return wrap.Null(), nil
		}

		return wrap.Number(v.Int64), nil
	case libdb.ColumnTypeFloat:
		v, err := dbVal.AsFloat()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to parse DB column value as float value")
		}

		if !v.Valid {
			return wrap.Null(), nil
		}

		return wrap.Number(v.Float64), nil
	case libdb.ColumnTypeString:
		v, err := dbVal.AsString()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to parse DB column value as string value")
		}

		if !v.Valid {
			return wrap.Null(), nil
		}

		return wrap.String(v.String), nil
	case libdb.ColumnTypeTime:
		v, err := dbVal.AsTime()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to parse DB column value as time.Time value")
		}

		if !v.Valid {
			return wrap.Null(), nil
		}

		return wrap.String(v.Time.Format(time.RFC3339)), nil
	default:
		v, err := dbVal.AsBytes()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to parse DB column value as []byte value")
		}

		if !v.Valid {
			return wrap.Null(), nil
		}

		return wrap.String(string(v.Bytes)), nil
	}
}

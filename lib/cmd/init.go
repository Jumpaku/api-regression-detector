package cmd

import (
	"context"
	"sort"

	libdb "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/lib/topological_sort"
)

func Init(ctx context.Context,
	db libdb.DB,
	jsonTables tables.InitTables,
	schemaGetter SchemaGetter,
	rowClearer RowClearer,
	rowCreator RowCreator,
) error {
	err := db.RunTransaction(ctx, func(ctx context.Context, tx libdb.Tx) error {
		tableSchema := map[string]libdb.Schema{}
		for _, table := range jsonTables {
			if _, ok := tableSchema[table.Name]; ok {
				continue
			}

			schema, err := schemaGetter.GetSchema(ctx, tx, table.Name)
			if err != nil {
				return errors.Wrap(
					errors.DBFailure.Err(err),
					errors.Info{"tableName": table.Name}.AppendTo("fail to get schema of table"))
			}

			tableSchema[table.Name] = schema
		}

		tableNames := []string{}
		dependencies := topological_sort.NewGraph[string]()
		for tableName, schema := range tableSchema {
			tableNames = append(tableNames, tableName)
			for _, depended := range schema.Dependencies {
				dependencies.Arrow(tableName, depended)
			}
		}

		tableNamesOrder, ok := topological_sort.Perform(dependencies)
		if !ok {
			return errors.DBFailure.New(
				errors.Info{"dependencies": dependencies}.AppendTo("fail to clear rows in table due to cycle"))
		}

		sort.Slice(tableNames, func(i, j int) bool {
			return tableNamesOrder[tableNames[i]] < tableNamesOrder[tableNames[j]]
		})

		for _, tableName := range tableNames {
			err := rowClearer.ClearRows(ctx, tx, tableName)
			if err != nil {
				return errors.Wrap(
					errors.DBFailure.Err(err),
					errors.Info{"tableName": tableName}.AppendTo("fail to clear rows in table"))
			}
		}

		for _, table := range jsonTables {
			err := rowCreator.CreateRows(ctx, tx, table.Name, tableSchema[table.Name], table.Rows)
			if err != nil {
				return errors.Wrap(
					errors.DBFailure.Err(err),
					errors.Info{"tableName": table.Name}.AppendTo("fail to create rows in table"))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(errors.DBFailure.Err(err), "fail to run transaction for Init")
	}

	return nil
}

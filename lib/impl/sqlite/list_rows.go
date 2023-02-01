package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type selectOperation struct {
}

func ListRows() selectOperation {
	return selectOperation{}
}

var _ cmd.RowLister = selectOperation{}

func (o selectOperation) ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) (rows []db.Row, err error) {
	rows, err = tx.Read(ctx, fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, tableName, strings.Join(schema.PrimaryKeys, ", ")), nil)
	if err != nil {
		return nil, err
	}
	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range schema.ColumnTypes.GetColumnNames() {
			col, ok := row.GetColumnValue(columnName)
			if !ok {
				return nil, fmt.Errorf("column %s not found", columnName)
			}
			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, fmt.Errorf("column %s not found", columnName)
			}
			outRow.SetColumnValue(columnName, col, typ)
		}
		out.Rows = append(out.Rows, outRow)
	}
	return out.Rows, nil
}
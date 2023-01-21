package mysql

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type selectOperation struct {
}

func Select() selectOperation {
	return selectOperation{}
}

var _ cmd.Select = selectOperation{}

func (o selectOperation) Select(ctx context.Context, exec db.Exec, table string) (rows db.Rows, err error) {
	rows, err = exec.Read(ctx, fmt.Sprintf(`SELECT * FROM %s`, table), nil)
	if err != nil {
		return nil, err
	}
	for i, row := range rows {
		rows[i] = db.Row{}
		for col, val := range row {
			if val == nil {
				rows[i][col] = nil
			} else {
				rows[i][col] = string(val.([]byte))
			}
		}
	}
	return rows, nil
}

package sqlite

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type truncateOperation struct{}

func ClearRows() truncateOperation {
	return truncateOperation{}
}

var _ cmd.RowClearer = truncateOperation{}

func (o truncateOperation) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	errInfo := errors.Info{tableName: tableName}

	foreignKeys, err := getForeignKeyCheck(ctx, tx)
	if err != nil {
		return errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to get foreign_keys variable"))
	}

	if foreignKeys {
		err := tx.Write(ctx, `PRAGMA foreign_keys = 0`, nil)
		if err != nil {
			return errors.Wrap(
				errors.DBFailure.Err(err),
				errInfo.AppendTo("fail to disable foreign key check"))
		}
	}

	err = tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s`, tableName), nil)
	if err != nil {
		return errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to delete all rows in table"))
	}

	if foreignKeys {
		err := tx.Write(ctx, `PRAGMA foreign_keys = 1`, nil)
		if err != nil {
			return errors.Wrap(
				errors.DBFailure.Err(err),
				errInfo.AppendTo("fail to enable foreign key check"))
		}
	}

	autoInclement, err := getAutoInclementEnable(ctx, tx)
	if err != nil {
		return errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to check auto inclement is enabled"))
	}

	if autoInclement {
		err = tx.Write(ctx, `DELETE FROM sqlite_sequence WHERE name = ?`, []any{tableName})
		if err != nil {
			return errors.Wrap(
				errors.DBFailure.Err(err),
				errInfo.AppendTo("fail to reset auto inclement in table"))
		}
	}

	return nil
}

func getForeignKeyCheck(ctx context.Context, tx db.Tx) (bool, error) {
	rows, err := tx.Read(ctx, `SELECT foreign_keys AS foreign_keys FROM pragma_foreign_keys()`, nil)
	if err != nil {
		return false, errors.Wrap(
			errors.DBFailure.Err(err),
			"fail to get foreign_keys variable")
	}

	foreignKeys, ok := rows[0]["foreign_keys"]
	if !ok {
		return false, errors.BadKeyAccess.New("foreign_keys column not found")
	}

	foreignKeysInteger, err := foreignKeys.AsInteger()
	if err != nil {
		return false, errors.Wrap(
			errors.BadConversion.Err(err),
			errors.Info{"foreignKeys": foreignKeys}.AppendTo("fail to get foreign_keys variable"))
	}

	return foreignKeysInteger.Int64 != 0, nil
}

func getAutoInclementEnable(ctx context.Context, tx db.Tx) (bool, error) {
	rows, err := tx.Read(ctx, `SELECT name FROM sqlite_master WHERE name='sqlite_sequence'`, nil)
	if err != nil {
		return false, errors.Wrap(
			errors.DBFailure.Err(err),
			"fail to get foreign_keys variable")
	}

	return len(rows) > 0, nil
}

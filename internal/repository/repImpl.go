package repository

import (
	"context"
	errorslist "javacode/internal/errors"

	"github.com/pkg/errors"
)

func (r *repImpl) SumInc(ctx context.Context, sumInc int) (int, error) {

	repSum := 0

	tx, err := r.DB.Begin()
	if err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	defer tx.Rollback()

	if err := r.DB.QueryRowContext(ctx, "SELECT sum FROM javacode FOR UPDATE").Scan(&repSum); err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	if err := r.DB.QueryRowContext(ctx, "UPDATE javacode SET sum=sum+$1 RETURNING sum", sumInc).Scan(&repSum); err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	return repSum, nil
}

func (r *repImpl) SumDec(ctx context.Context, sumDec int) (int, error) {

	repSum := 0

	tx, err := r.DB.Begin()
	if err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}
	defer tx.Rollback()

	if err := r.DB.QueryRowContext(ctx, "SELECT sum FROM javacode FOR UPDATE").Scan(&repSum); err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	if repSum < sumDec {
		return 0, errorslist.ErrInsufficientFunds
	}

	if err := r.DB.QueryRowContext(ctx, "UPDATE javacode SET sum=sum-$1 RETURNING sum", sumDec).Scan(&repSum); err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Wrap(err, errorslist.ErrSqlFailed)
	}

	return repSum, nil
}

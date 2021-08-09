package sql

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var transactionManagerContextKey = "!!TXProvider"

func Transactional(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			transactionManager := &TransactionManager{db: db}
			tmContext := context.WithValue(ctx.Request().Context(), transactionManagerContextKey, transactionManager)
			ctx.SetRequest(ctx.Request().WithContext(tmContext))

			err := next(ctx)
			if err != nil {
				transactionManager.Rollback()
			} else {
				return transactionManager.Commit()
			}
			return err
		}
	}
}

type TransactionManager struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func (tm *TransactionManager) getTransaction() (*sqlx.Tx, error) {
	var err error
	if tm.tx == nil {
		tm.tx, err = tm.db.Beginx()
	}
	return tm.tx, err
}

func (tm *TransactionManager) Rollback() {
	if tm.tx != nil {
		if rollbackErr := tm.tx.Rollback(); rollbackErr != nil {
			logrus.Errorf("Error while rolling back transaction: %v", rollbackErr)
		}
		tm.tx = nil
	}
}

func (tm *TransactionManager) Commit() error {
	var err error
	if tm.tx != nil {
		if err = tm.tx.Commit(); err != nil {
			logrus.Errorf("Error while committing transaction: %v", err)
		}
		tm.tx = nil
	}
	return err
}

func ExecuteTransactional(db *sqlx.DB, acceptor func(ctx context.Context) error) error {
	tm := &TransactionManager{db: db}
	ctx := context.WithValue(context.Background(), transactionManagerContextKey, tm)
	if err := acceptor(ctx); err != nil {
		tm.Rollback()
		return err
	} else {
		return tm.Commit()
	}
}

func GetTransactionManager(ctx context.Context) (*TransactionManager, error) {
	tm, ok := ctx.Value(transactionManagerContextKey).(*TransactionManager)
	if !ok {
		return nil, errors.New("transaction manager not registered")
	}
	return tm, nil
}

func GetTransaction(ctx context.Context) (*sqlx.Tx, error) {
	tm, err := GetTransactionManager(ctx)
	if err != nil {
		return nil, err
	}
	return tm.getTransaction()
}

package inbox

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"time"
)

const notificationSchema = `
	CREATE TABLE IF NOT EXISTS notification (
		id char(36) NOT NULL,
		customer_id integer(11) NOT NULL,
		sender_did varchar(100) NOT NULL,
		date DATETIME NOT NULL,
		PRIMARY KEY (id)
	);
`

func NewRepository(db *sqlx.DB) Repository {
	tx, _ := db.Beginx()
	tx.MustExec(notificationSchema)
	if err := tx.Commit(); err != nil {
		panic(err)
	}
	return Repository{db: db}
}

type sqlNotification struct {
	ID         string    `db:"id"`
	CustomerID int        `db:"customer_id"`
	SenderDID  string    `db:"sender_did"`
	Date       time.Time `db:"date"`
}

type Repository struct {
	db *sqlx.DB
}

func (f Repository) getAll(ctx context.Context, customerID int) ([]sqlNotification, error) {
	const query = `SELECT * FROM notification WHERE customer_id = ?`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	var result []sqlNotification
	err = tx.SelectContext(ctx, &result, query, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return []sqlNotification{}, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

func (f Repository) registerNotification(ctx context.Context, customerID int, senderDID string) error {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return err
	}
	const query = `INSERT INTO notification 
		(id, customer_id, sender_did, date)
		values(:id, :customer_id, :sender_did, :date)`
	_, err = tx.NamedExecContext(ctx, query, sqlNotification{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		SenderDID:  senderDID,
		Date:       time.Now(),
	})
	return err
}

package transfer

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type sqlTransfer struct {
	ID          string       `db:"id"`
	CustomerID  string       `db:"customer_id"`
	Dob         sql.NullTime `db:"date"`
	status      string       `db:"status"`
	description string       `db:"description"`
}

const schema = `
	CREATE TABLE IF NOT EXISTS transfer (
		id char(36) NOT NULL,
		customer_id varchar(100) NOT NULL,
		date DATETIME DEFAULT NULL,
		status char(10) NOT NULL DEFAULT 'created',
		description varchar(200) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE(customer_id, id)
	);
`

type SQLiteTransferRepository struct {
	factory Factory
	db      *sqlx.DB
}

func NewSQLiteTransferRepository(factory Factory, db *sqlx.DB) *SQLiteTransferRepository {
	if db == nil {
		panic("missing db")
	}
	db.MustExec(schema)

	return &SQLiteTransferRepository{
		factory: factory,
		db:      db,
	}
}

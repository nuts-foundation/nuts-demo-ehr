package acl

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
)

type aclEntry struct {
	ID string `db:"id"`
	// TenantDID is the DID of the tenant that is granting access. It is the owner of the resource.
	TenantDID string `db:"tenant_did"`
	// AuthorizedDID is the DID of the entity that is being granted access
	AuthorizedDID string `db:"authorized_did"`
	Operation     string `db:"operation"`
	Resource      string `db:"resource"`
}

type AuthorizedResource struct {
	Resource  string
	Operation string
}

const schema = `
	CREATE TABLE IF NOT EXISTS acl (
		id char(36) NOT NULL,
		tenant_did char(240) NOT NULL,
		authorized_did char(240) NOT NULL,
	    operation char(50) NOT NULL,
		resource varchar(400) NOT NULL,
		PRIMARY KEY (id)
	);
`

func NewRepository(db *sqlx.DB) (*Repository, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec(schema)
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Repository{}, nil
}

type Repository struct {
}

func (r Repository) GrantAccess(ctx context.Context, tenantDID, authorizedDID, operation string, resource string) error {
	// First check whether the authorized does not already have access
	if hasAccess, err := r.HasAccess(ctx, tenantDID, authorizedDID, operation, resource); err != nil {
		return err
	} else if hasAccess {
		return nil
	}

	const query = `INSERT INTO acl
		(id, tenant_did, authorized_did, operation, resource)
		VALUES (:id, :tenant_did, :authorized_did, :operation, :resource)`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(query, aclEntry{
		ID:            uuid.NewString(),
		TenantDID:     tenantDID,
		AuthorizedDID: authorizedDID,
		Operation:     operation,
		Resource:      resource,
	})
	return err
}

func (r Repository) AuthorizedResources(ctx context.Context, tenantDID, authorizedDID string) ([]AuthorizedResource, error) {
	const query = `SELECT operation, resource FROM acl WHERE tenant_did = ? AND authorized_did = ?`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	var result []AuthorizedResource
	err = tx.Select(&result, query, tenantDID, authorizedDID)
	return result, err
}

// AuthorizedParties returns the DIDs of the parties that are authorized to perform the given operation on the given resource
func (r Repository) AuthorizedParties(ctx context.Context, tenantDID, resource, operation string) ([]string, error) {
	const query = `SELECT authorized_did FROM acl WHERE tenant_did = ? AND resource = ? AND operation = ?`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	var result []string
	err = tx.Select(&result, query, tenantDID, resource, operation)
	return result, err
}

func (r Repository) HasAccess(ctx context.Context, tenantDID, authorizedDID, operation, resource string) (bool, error) {
	const query = `SELECT COUNT(*) FROM acl WHERE tenant_did = ? AND authorized_did = ? AND operation = ? AND resource = ?`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return false, err
	}
	var count int
	err = tx.Get(&count, query, tenantDID, authorizedDID, operation, resource)
	return count > 0, err
}

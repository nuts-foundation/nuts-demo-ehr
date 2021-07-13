package transfer

import (
	"context"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(ctx context.Context, customerID, id string) (*domain.Transfer, error)
	FindByPatientID(ctx context.Context, customerID, patientID string) ([]domain.Transfer, error)
	Create(ctx context.Context, customerID, dossierID, description string, date time.Time) (*domain.Transfer, error)
	Update(ctx context.Context, customerID, transferID string, updateFn func(c domain.Transfer) (*domain.Transfer, error)) (*domain.Transfer, error)
	Cancel(ctx context.Context, customerID, id string)
	CreateNegotiation(ctx context.Context, customerID, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error)
	ListNegotiations(ctx context.Context, customerID, transferID string) ([]domain.TransferNegotiation, error)
}

type Factory struct{}

func (f Factory) NewTransfer(description string, date time.Time, dossierID domain.ObjectID) *domain.Transfer {
	return &domain.Transfer{
		Id:           domain.ObjectID(uuid.NewString()),
		DossierID:    dossierID,
		Description:  description,
		Status:       domain.TransferStatusCreated,
		TransferDate: openapi_types.Date{Time: date},
	}
}

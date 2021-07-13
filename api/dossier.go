package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/sirupsen/logrus"
	"net/http"
)

type GetDossierParams = domain.GetDossierParams
type CreateDossierRequest = domain.CreateDossierRequest

func (w Wrapper) GetDossier(ctx echo.Context, params GetDossierParams) error {
	panic("not implemented")
}

func (w Wrapper) CreateDossier(ctx echo.Context) error {
	request := domain.CreateDossierRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	logrus.Infof("Creating dossier (name=%s, patientID=%s)", request.Name, request.PatientID)
	return ctx.JSON(http.StatusOK, domain.Dossier{
		Id:        "712BC28B-C368-475F-9D98-1F80F61E0956",
		Name:      request.Name,
		PatientID: request.PatientID,
	})
}

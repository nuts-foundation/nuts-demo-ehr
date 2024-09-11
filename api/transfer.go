package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"net/http"
	"net/url"
	"strings"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type GetPatientTransfersParams = types.GetPatientTransfersParams

func (w Wrapper) CreateTransfer(ctx echo.Context) error {
	request := types.CreateTransferRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderService.CreateTransfer(ctx.Request().Context(), cid, request)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) GetPatientTransfers(ctx echo.Context, params GetPatientTransfersParams) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfers, err := w.TransferSenderRepo.FindByPatientID(ctx.Request().Context(), cid, params.PatientID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfers)
}

func (w Wrapper) GetTransfer(ctx echo.Context, transferID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	transfer, err := w.TransferSenderService.GetTransferByID(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) ChangeTransferRequestState(ctx echo.Context, requesterDID string, fhirTaskID string, params ChangeTransferRequestStateParams) error {
	updateRequest := &types.TransferNegotiationStatus{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	err = w.TransferReceiverService.UpdateTransferRequestState(ctx.Request().Context(), cid, requesterDID, fhirTaskID, string(updateRequest.Status))
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) UpdateTransfer(ctx echo.Context, transferID string) error {
	updateRequest := &types.TransferProperties{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	_, err = w.TransferSenderRepo.Update(ctx.Request().Context(), cid, transferID, func(t *types.Transfer) (*types.Transfer, error) {
		//t.Description = updateRequest.Description
		t.TransferDate = updateRequest.TransferDate
		return t, nil
	})
	if err != nil {
		return err
	}

	transfer, err := w.TransferSenderService.GetTransferByID(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) CancelTransfer(ctx echo.Context, transferID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderRepo.Cancel(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) StartTransferNegotiation(ctx echo.Context, transferID string) error {
	request := types.CreateTransferNegotiationRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiation, err := w.TransferSenderService.CreateNegotiation(ctx.Request().Context(), cid, transferID, request.OrganizationID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *negotiation)
}

func (w Wrapper) AssignTransferDirect(ctx echo.Context, transferID string) error {
	request := types.CreateTransferNegotiationRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}
	_, err = w.TransferSenderService.AssignTransfer(ctx.Request().Context(), *customer, transferID, request.OrganizationID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) ListTransferNegotiations(ctx echo.Context, transferID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiations, err := w.TransferSenderRepo.ListNegotiations(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}
	// Enrich with organization info
	for i, negotiation := range negotiations {
		organization, err := w.OrganizationRegistry.Get(ctx.Request().Context(), negotiation.OrganizationID)
		if err != nil {
			logrus.Warnf("Error while fetching organization info for negotiation (DID=%s): %v", negotiation.OrganizationID, err)
			continue
		}
		negotiations[i].Organization = types.FromNutsOrganization(*organization)
	}
	return ctx.JSON(http.StatusOK, negotiations)
}

func (w Wrapper) UpdateTransferNegotiationStatus(ctx echo.Context, transferID string, negotiationID string) error {
	request := types.TransferNegotiationStatus{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	newState := request.Status
	if newState == transfer.InProgressState {
		_, err = w.TransferSenderService.ConfirmNegotiation(ctx.Request().Context(), cid, transferID, negotiationID)
	} else if newState == transfer.CancelledState {
		_, err = w.TransferSenderService.CancelNegotiation(ctx.Request().Context(), cid, transferID, negotiationID)
	}
	if err != nil {
		return fmt.Errorf("unable to update transfer negotiation state: %w", err)
	}
	negotiation, err := w.TransferSenderRepo.UpdateNegotiationState(ctx.Request().Context(), cid, negotiationID, request.Status)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, negotiation)
}

func (w Wrapper) NotifyTransferUpdate(ctx echo.Context, taskID string) error {
	// This gets called by a transfer sending XIS to inform the local node there's FHIR tasks to be retrieved.
	// The PEP added introspection result to the X-Userinfo header
	b64IntrospectionResult := ctx.Request().Header.Get("X-Userinfo")
	//log.Errorf("X-Userinfo: %s", b64IntrospectionResult)
	if b64IntrospectionResult == "" {
		return errors.New("missing X-Userinfo header")
	}

	// b64 -> json string
	introspectionResult, err := base64.URLEncoding.DecodeString(b64IntrospectionResult)
	if err != nil {
		return fmt.Errorf("failed to base64 decode X-Userinfo header: %w", err)
	}

	// json string -> map
	target := map[string]interface{}{}
	err = json.Unmarshal(introspectionResult, &target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal X-Userinfo header: %w", err)
	}

	// client_id for senderDID and sub for customerDID
	_ = json.Unmarshal([]byte(introspectionResult), &target)
	issuerURLStr := target["iss"].(string)
	// get senderDID via custom policy param
	senderClientID := target["client_id"].(string)
	// we need the subjectID, which is at the end of the path, not panic safe
	issuerURL, _ := url.Parse(issuerURLStr)
	idx := strings.LastIndex(issuerURL.Path, "/")
	customerID := issuerURL.Path[idx+1:]

	codeError := datatypes.Code("error")
	codeInvalid := datatypes.Code("invalid")
	severityError := datatypes.Code("error")
	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {

		return ctx.JSON(http.StatusInternalServerError, &resources.OperationOutcome{
			Domain: resources.Domain{
				Text: &datatypes.Narrative{
					Div: fhir.ToStringPtr("an error occurred"),
				},
			},
			Issue: []resources.OperationOutcomeIssue{
				{
					Code:     &codeError,
					Severity: &severityError,
					Details: &datatypes.CodeableConcept{
						Text: fhir.ToStringPtr(err.Error()),
					},
				},
			},
		})
	}

	if customer == nil {
		logrus.Warnf("Received transfer notification for unknown customer: %s", customerID)

		return ctx.JSON(http.StatusNotFound, &resources.OperationOutcome{
			Domain: resources.Domain{
				Text: &datatypes.Narrative{
					Div: fhir.ToStringPtr("taskOwner unknown on this server"),
				},
			},
			Issue: []resources.OperationOutcomeIssue{
				{
					Code:     &codeInvalid,
					Severity: &codeError,
					Details: &datatypes.CodeableConcept{
						Text: fhir.ToStringPtr(fmt.Sprintf("received transfer notification for unknown taskOwner with ID: %s", senderClientID)),
					},
				},
			},
		})
	}

	if err := w.NotificationHandler.Handle(ctx.Request().Context(), notification.Notification{
		TaskID:     taskID,
		SenderID:   senderClientID,
		CustomerID: customerID,
	}); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusAccepted)
}

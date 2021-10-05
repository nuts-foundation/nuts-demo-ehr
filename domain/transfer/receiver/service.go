package receiver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type TransferService interface {
	// CreateOrUpdate creates or updates an incoming transfer record in the local storage
	CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error
	UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error
	GetTransferRequest(ctx context.Context, customerID int, requesterDID string, fhirTaskID string) (*types.TransferRequest, error)
}

type service struct {
	transferRepo           TransferRepository
	notifier               transfer.Notifier
	auth                   auth.Service
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo           customers.Repository
	registry               registry.OrganizationRegistry
	vcr                    registry.VerifiableCredentialRegistry
}

func NewTransferService(authService auth.Service, localFHIRClientFactory fhir.Factory, transferRepository TransferRepository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry) TransferService {
	return &service{
		auth:                   authService,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		registry:               organizationRegistry,
		vcr:                    vcr,
		notifier:               transfer.FireAndForgetNotifier{},
	}
}

func (s service) CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error {
	_, err := s.transferRepo.CreateOrUpdate(ctx, status, fhirTaskID, customerID, senderDID)
	if err != nil {
		return err
	}

	return nil
}

func (s service) UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return err
	}

	taskPath := fmt.Sprintf("/Task/%s", fhirTaskID)
	fhirClient, err := s.getRemoteFHIRClient(ctx, requesterDID, *customer.Did, taskPath)
	if err != nil {
		return err
	}

	remoteFHIRRepo := fhir.NewFHIRRepository(fhirClient)
	task, err := remoteFHIRRepo.GetTask(ctx, fhirTaskID)
	if err != nil {
		return err
	}

	// state machine
	if (*task.Status == transfer.InProgressState && newState == transfer.CompletedState) ||
		(*task.Status == transfer.RequestedState && newState == transfer.AcceptedState) {
		err = remoteFHIRRepo.UpdateTaskStatus(ctx, fhirTaskID, newState)
		if err != nil {
			return err
		}
		// update was a success. Get the remote task again and update the local transfer_request
		task, err = remoteFHIRRepo.GetTask(ctx, fhirTaskID)
		if err != nil {
			return err
		}
		_, err = s.transferRepo.CreateOrUpdate(ctx, fhir.FromCodePtr(task.Status), fhirTaskID, customerID, requesterDID)
		if err != nil {
			return fmt.Errorf("could update incomming transfers with new state")
		}
		return nil
	}

	return fmt.Errorf("invalid state change from %s to %s", *task.Status, newState)
}

func (s service) taskContainsCode(task resources.Task, code datatypes.Code) bool {
	for _, input := range task.Input {
		if fhir.FromCodePtr(input.Type.Coding[0].Code) == string(code) {
			return true
		}
	}

	return false
}

func (s service) GetTransferRequest(ctx context.Context, customerID int, requesterDID string, fhirTaskID string) (*types.TransferRequest, error) {
	const getTransferRequestErr = "unable to get transferRequest: %w"

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, fmt.Errorf("unable to find customer: %w", err)
	}

	taskPath := fmt.Sprintf("/Task/%s", fhirTaskID)
	fhirClient, err := s.getRemoteFHIRClient(ctx, requesterDID, *customer.Did, taskPath)
	if err != nil {
		return nil, err
	}

	fhirReceiverService := eoverdracht.NewReceiverFHIRTransferService(fhir.NewFHIRRepository(fhirClient))

	task, err := fhirReceiverService.GetTask(ctx, fhirTaskID)
	if err != nil {
		return nil, fmt.Errorf(getTransferRequestErr, err)
	}

	if task.NursingHandoffID == nil {
		return nil, fmt.Errorf(getTransferRequestErr, errors.New("invalid task, expected an advanceNotice composition"))
	}

	advanceNotice, err := s.getAdvanceNotice(ctx, fhirClient, *task.AdvanceNoticeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get advance notice: %w", err)
	}
	domainAdvanceNotice, err := eoverdracht.FHIRAdvanceNoticeToDomainTransfer(advanceNotice)
	if err != nil {
		return nil, err
	}

	organization, err := s.registry.Get(ctx, requesterDID)
	if err != nil {
		return nil, fmt.Errorf("unable to get organization from registry: %w", err)
	}

	// TODO: Do we need nil checks?
	transferRequest := types.TransferRequest{
		Sender:        *organization,
		AdvanceNotice: domainAdvanceNotice,
		Status:        task.Status,
	}

	// If the task input contains the nursing handoff, add that one too.
	if task.NursingHandoffID != nil {
		nursingHandoff, err := s.getNursingHandoff(ctx, fhirClient, *task.NursingHandoffID)
		if err != nil {
			return nil, fmt.Errorf("unable to get nursing handoff: %w", err)
		}
		domainTransfer, err := eoverdracht.FHIRNursingHandoffToDomainTransfer(nursingHandoff)
		if err != nil {
			return nil, fmt.Errorf("unable to convert fhir nursing handoff to domain transfer: %w", err)
		}
		transferRequest.NursingHandoff = &domainTransfer
	}

	return &transferRequest, nil
}

func (s service) getRemoteFHIRClient(ctx context.Context, authorizerDID string, localRequesterDID string, resource string) (fhir.Client, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, authorizerDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", authorizerDID, err)
	}
	credentials, err := s.vcr.FindAuthorizationCredentials(ctx, transfer.SenderServiceName, localRequesterDID, resource)

	var transformed = make([]vc.VerifiableCredential, len(credentials))
	for i, c := range credentials {
		bytes, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		tCred := vc.VerifiableCredential{}
		if err = json.Unmarshal(bytes, &tCred); err != nil {
			return nil, err
		}
		transformed[i] = tCred
	}

	accessToken, err := s.auth.RequestAccessToken(ctx, localRequesterDID, authorizerDID, transfer.SenderServiceName, transformed)
	if err != nil {
		return nil, err
	}

	return fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken))(), nil
}

// getAdvanceNotice fetches a complete nursing handoff from a FHIR server
func (s service) getNursingHandoff(ctx context.Context, fhirClient fhir.Client, fhirCompositionPath string) (eoverdracht.NursingHandoff, error) {
	nursingHandoff := eoverdracht.NursingHandoff{}

	// Fetch the composition
	err := fhirClient.ReadOne(ctx, "/"+fhirCompositionPath, &nursingHandoff.Composition)
	if err != nil {
		return eoverdracht.NursingHandoff{}, fmt.Errorf("error while fetching the advance notice composition(composition-id=%s): %w", fhirCompositionPath, err)
	}

	// Fetch the Patient
	err = fhirClient.ReadOne(ctx, "/"+fhir.FromStringPtr(nursingHandoff.Composition.Subject.Reference), &nursingHandoff.Patient)
	if err != nil {
		return eoverdracht.NursingHandoff{}, fmt.Errorf("error while fetching the transfer subject (patient): %w", err)
	}

	// Fetch the careplan
	careplan, err := eoverdracht.FilterCompositionSectionByType(nursingHandoff.Composition.Section, eoverdracht.CarePlanCode)
	if err != nil {
		return eoverdracht.NursingHandoff{}, err
	}

	// Fetch the nursing diagnosis
	nursingDiagnosis, err := eoverdracht.FilterCompositionSectionByType(careplan.Section, eoverdracht.NursingDiagnosisCode)
	if err != nil {
		return eoverdracht.NursingHandoff{}, err
	}

	// the nursing diagnosis contains both conditions and procedures
	for _, entry := range nursingDiagnosis.Entry {
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Condition") {
			conditionID := fhir.FromStringPtr(entry.Reference)
			condition := resources.Condition{}
			err = fhirClient.ReadOne(ctx, "/"+conditionID, &condition)
			if err != nil {
				return eoverdracht.NursingHandoff{}, fmt.Errorf("error while fetching a advance notice condition (condition-id=%s): %w", conditionID, err)
			}
			nursingHandoff.Problems = append(nursingHandoff.Problems, condition)
		}
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Procedure") {
			procedureID := fhir.FromStringPtr(entry.Reference)
			procedure := eoverdracht.Procedure{}
			err = fhirClient.ReadOne(ctx, "/"+procedureID, &procedure)
			if err != nil {
				return eoverdracht.NursingHandoff{}, fmt.Errorf("error while fetching a advance notice procedure (procedure-id=%s): %w", procedureID, err)
			}
			nursingHandoff.Interventions = append(nursingHandoff.Interventions, procedure)
		}
	}

	return nursingHandoff, nil
}

// getAdvanceNotice fetches a complete advance notice from a FHIR server
func (s service) getAdvanceNotice(ctx context.Context, fhirClient fhir.Client, fhirCompositionPath string) (eoverdracht.AdvanceNotice, error) {
	advanceNotice := eoverdracht.AdvanceNotice{}

	err := fhirClient.ReadOne(ctx, "/"+fhirCompositionPath, &advanceNotice.Composition)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching the advance notice composition(composition-id=%s): %w", fhirCompositionPath, err)
	}

	if advanceNotice.Composition.Subject.Reference != nil {
		err = fhirClient.ReadOne(ctx, "/"+fhir.FromStringPtr(advanceNotice.Composition.Subject.Reference), &advanceNotice.Patient)
		if err != nil {
			return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching the transfer subject (patient): %w", err)
		}
	}

	careplan, err := eoverdracht.FilterCompositionSectionByType(advanceNotice.Composition.Section, eoverdracht.CarePlanCode)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, err
	}

	nursingDiagnosis, err := eoverdracht.FilterCompositionSectionByType(careplan.Section, eoverdracht.NursingDiagnosisCode)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, err
	}

	// the nursing diagnosis contains both conditions and procedures
	for _, entry := range nursingDiagnosis.Entry {
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Condition") {
			conditionID := fhir.FromStringPtr(entry.Reference)
			condition := resources.Condition{}
			err = fhirClient.ReadOne(ctx, "/"+conditionID, &condition)
			if err != nil {
				return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching a advance notice condition (condition-id=%s): %w", conditionID, err)
			}
			advanceNotice.Problems = append(advanceNotice.Problems, condition)
		}
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Procedure") {
			procedureID := fhir.FromStringPtr(entry.Reference)
			procedure := eoverdracht.Procedure{}
			err = fhirClient.ReadOne(ctx, "/"+procedureID, &procedure)
			if err != nil {
				return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching a advance notice procedure (procedure-id=%s): %w", procedureID, err)
			}
			advanceNotice.Interventions = append(advanceNotice.Interventions, procedure)
		}
	}

	return advanceNotice, nil
}

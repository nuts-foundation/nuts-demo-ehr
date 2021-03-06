package patients

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(ctx context.Context, customerID, id string) (*domain.Patient, error)
	Update(ctx context.Context, customerID, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (*domain.Patient, error)
	NewPatient(ctx context.Context, customerID string, patient domain.PatientProperties) (*domain.Patient, error)
	All(ctx context.Context, customerID string) ([]domain.Patient, error)
}

type Factory struct{}

// NewUUIDPatient creates a new patient from a list of properties. It generates a new UUID for the patientID.
func (f Factory) NewUUIDPatient(patientProperties domain.PatientProperties) (*domain.Patient, error) {
	if patientProperties.Gender == "" {
		patientProperties.Gender = domain.PatientPropertiesGenderUnknown
	}
	return &domain.Patient{
		ObjectID:         domain.ObjectID(uuid.NewString()),
		PatientProperties: patientProperties,
	}, nil
}

type MemoryPatientRepository struct {
	// indices on customerID and patientID
	patients map[string]map[domain.ObjectID]domain.Patient
	lock     *sync.RWMutex
	factory  Factory
}

func NewMemoryPatientRepository(factory Factory) *MemoryPatientRepository {
	return &MemoryPatientRepository{
		patients: map[string]map[domain.ObjectID]domain.Patient{},
		lock:     &sync.RWMutex{},
		factory:  factory,
	}
}

func (r MemoryPatientRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Patient, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.getPatient(ctx, customerID, id)
}

func (r MemoryPatientRepository) Update(ctx context.Context, customerID, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (*domain.Patient, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	patient, err := r.getPatient(ctx, customerID, id)
	if err != nil {
		return nil, err
	}
	updatedPatient, err := updateFn(*patient)
	if err != nil {
		return nil, err
	}
	r.patients[customerID][domain.ObjectID(id)] = *updatedPatient
	return patient, nil
}

func (r MemoryPatientRepository) NewPatient(ctx context.Context, customerID string, patientProperties domain.PatientProperties) (*domain.Patient, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	customerPatients, ok := r.patients[customerID]
	if !ok {
		customerPatients = map[domain.ObjectID]domain.Patient{}
		r.patients[customerID] = customerPatients
	}
	patient, err := r.factory.NewUUIDPatient(patientProperties)
	if err != nil {
		return nil, err
	}
	customerPatients[patient.ObjectID] = *patient
	return patient, nil
}

func (r MemoryPatientRepository) All(ctx context.Context, customerID string) ([]domain.Patient, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	result := make([]domain.Patient, len(r.patients[customerID]))
	idx := 0
	for _, p := range r.patients[customerID] {
		result[idx] = p
		idx++
	}
	return result, nil
}

func (r MemoryPatientRepository) getPatient(ctx context.Context, customerID, patientID string) (*domain.Patient, error) {
	customerPatients, ok := r.patients[customerID]
	if !ok {
		return nil, nil
	}

	patient, ok := customerPatients[domain.ObjectID(patientID)]
	if !ok {
		return nil, nil
	}
	return &patient, nil
}

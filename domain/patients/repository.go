package patients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

const maxMinimumAgeForAvatar = 76

type Repository interface {
	FindByID(ctx context.Context, customerID int, id string) (*domain.Patient, error)
	Update(ctx context.Context, customerID int, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (*domain.Patient, error)
	NewPatient(ctx context.Context, customerID int, patient domain.PatientProperties) (*domain.Patient, error)
	All(ctx context.Context, customerID int, name *string) ([]domain.Patient, error)
}

type Factory struct{}

// NewUUIDPatient creates a new patient from a list of properties. It generates a new UUID for the patientID.
func (f Factory) NewUUIDPatient(patientProperties domain.PatientProperties) (*domain.Patient, error) {
	if patientProperties.Gender == "" {
		patientProperties.Gender = domain.PatientPropertiesGenderUnknown
	}
	return &domain.Patient{
		ObjectID:          domain.ObjectID(uuid.NewString()),
		PatientProperties: patientProperties,
	}, nil
}

func (f Factory) NewPatientWithAvatar(properties domain.PatientProperties) (*domain.Patient, error) {
	patient, err := f.NewUUIDPatient(properties)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		IdleConnTimeout: 5 * time.Second,
	}
	client := &http.Client{Transport: tr}
	const fakeFaceAPIURL = "https://fakeface.rest/face/json"
	url, err := url.Parse(fakeFaceAPIURL)
	if err != nil {
		return patient, err
	}

	var gender string
	switch patient.Gender {
	case domain.PatientPropertiesGenderMale:
		gender = "male"
	case domain.PatientPropertiesGenderFemale:
		gender = "female"
	default:
		// For "other" and "unknown" we take a random gender
		if rand.Intn(1000)%2 == 0 {
			gender = "male"
		} else {
			gender = "female"
		}
	}
	q := url.Query()
	q.Add("gender", gender)
	if !patient.Dob.IsZero() {
		age := math.Floor(time.Since(patient.Dob.Time).Hours() / 24 / 365)
		minimumAge := int(age - 3)
		// FakeFace API doesn't return results for age >= 77, for some reason so we cap it at 76
		if minimumAge > maxMinimumAgeForAvatar {
			minimumAge = maxMinimumAgeForAvatar
		}
		q.Set("minimum_age", fmt.Sprintf("%d", minimumAge))
		q.Set("maximum_age", fmt.Sprintf("%d", int(age+3)))
	}
	url.RawQuery = q.Encode()
	resp, err := client.Get(url.String())
	if err != nil {
		return patient, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return patient, err
	}
	apiResponse := map[string]interface{}{}
	json.Unmarshal(body, &apiResponse)

	if avatarURL, ok := apiResponse["image_url"]; ok {
		tmp := avatarURL.(string)
		patient.AvatarUrl = &tmp
	}

	return patient, nil
}

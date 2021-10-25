package fhir

import (
	"context"
	"encoding/json"

	"github.com/stretchr/testify/assert"
)

type mockClient struct{
	t assert.TestingT
	ExpectedCreateOrUpdate map[string]interface{}
}

func NewMockClientWithExpectedCreateOrUpdate(t assert.TestingT, expected map[string]interface{}) mockClient {
	return mockClient{t: t, ExpectedCreateOrUpdate: expected}
}

func (m mockClient) CreateOrUpdate(ctx context.Context, resource interface{}) error {
	resourceJSON, _ := json.Marshal(resource)
	resourceMap := map[string]interface{}{}
	_ = json.Unmarshal(resourceJSON, &resourceMap)

	assert.Equal(m.t, m.ExpectedCreateOrUpdate, resourceMap)
	return nil
}

func (mockClient) ReadMultiple(ctx context.Context, path string, params map[string]string, results interface{}) error {
	panic("implement me")
}

func (mockClient) ReadOne(ctx context.Context, path string, result interface{}) error {
	panic("implement me")
}



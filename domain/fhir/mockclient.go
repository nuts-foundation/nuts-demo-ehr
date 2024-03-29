package fhir

import (
	"context"
	"encoding/json"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	t                      assert.TestingT
	ExpectedCreateOrUpdate map[string]interface{}
	readMock               map[string]map[string]interface{}
}

func NewMockClientWithExpectedCreateOrUpdate(t assert.TestingT, expected map[string]interface{}) mockClient {
	return mockClient{t: t, ExpectedCreateOrUpdate: expected}
}

func NewMockClientWithReadMock(t assert.TestingT, expected map[string]map[string]interface{}) mockClient {
	return mockClient{t: t, readMock: expected}
}

func (m mockClient) CreateOrUpdate(ctx context.Context, resource interface{}) error {
	resourceJSON, _ := json.Marshal(resource)
	_ = json.Unmarshal(resourceJSON, &resource)
	return nil
}

func (mockClient) ReadMultiple(ctx context.Context, path string, params map[string]string, results interface{}) error {
	panic("implement me")
}

func (m mockClient) ReadOne(ctx context.Context, path string, result interface{}) error {
	if mockData, ok := m.readMock[path]; ok {
		resourceJSON, _ := json.Marshal(mockData)
		return json.Unmarshal(resourceJSON, &result)
	}
	m.t.Errorf("unexpected call to ReadOne with path %s", path)
	return nil
}

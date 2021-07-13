package fhir

type Gateway interface {
	CreateTask(interface{}) error
	UpdateTask(interface{}) error
}

type StubGateway struct {
}

func (s StubGateway) CreateTask(i interface{}) error {
	return nil
}

func (s StubGateway) UpdateTask(i interface{}) error {
	return nil
}

package fhir

type Gateway interface {
	CreateTask(interface{}) error
	UpdateTask(interface{}) error
}

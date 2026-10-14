package validator

type Validator interface {
	ValidateStruct(i interface{}) error
	Engine() interface{}
}

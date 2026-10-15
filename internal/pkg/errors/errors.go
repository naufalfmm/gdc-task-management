package errors

import (
	v10 "github.com/go-playground/validator/v10"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/validator"
)

func GenerateError(data any, err error) error {
	if _, ok := err.(v10.ValidationErrors); !ok {
		return err
	}

	return validator.NewFromValidationErrors(data, err)
}

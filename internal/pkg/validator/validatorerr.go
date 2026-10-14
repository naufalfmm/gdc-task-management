package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	validatorErr struct {
		namespace   string
		fieldname   string
		jsonTag     string
		validateTag string
	}

	validatorErrJson struct {
		Error string `json:"error"`
	}
)

func (ve validatorErr) Error() string {
	return fmt.Sprintf("Error %s for %s", ve.validateTag, ve.jsonTag)
}

func (ve validatorErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(validatorErrJson{
		Error: ve.Error(),
	})
}

func NewFromValidationError(data interface{}, valErr validator.FieldError) validatorErr {
	var (
		fieldName, jsonTag string

		vo = reflect.ValueOf(data)
		to = reflect.TypeOf(data)

		re = regexp.MustCompile(`(.*?)(\[.*?\])`)
	)

	for _, cn := range strings.Split(valErr.StructNamespace(), ".")[1:] {
		submatches := re.FindAllStringSubmatch(cn, -1)
		additional := ""
		if len(submatches) > 0 {
			cn = submatches[0][1]
			additional = submatches[0][2]
		}

		tc, ok := to.FieldByName(cn)
		if !ok {
			return validatorErr{
				namespace:   valErr.StructNamespace(),
				fieldname:   valErr.Field(),
				jsonTag:     "",
				validateTag: valErr.ActualTag(),
			}
		}
		vo = vo.FieldByName(cn)

		fieldName += cn

		if len(jsonTag) > 0 {
			jsonTag += "."
		}
		jsonTag += (tc.Tag.Get("json") + additional)

		if tc.Type.Kind() == reflect.Slice {
			to = reflect.TypeOf(vo.Interface()).Elem()
			vo = reflect.Indirect(reflect.ValueOf(to))
		}

		if tc.Type.Kind() == reflect.Struct {
			to = reflect.TypeOf(vo.Interface())
			vo = reflect.Indirect(reflect.ValueOf(to))
		}
	}

	return validatorErr{
		namespace:   valErr.StructNamespace(),
		fieldname:   fieldName,
		jsonTag:     jsonTag,
		validateTag: valErr.ActualTag(),
	}
}

type (
	ValidatorErrs []validatorErr

	validatorErrsJson struct {
		Error []validatorErrJson `json:"errors"`
	}
)

func (ves ValidatorErrs) Error() string {
	var errString string
	for i, ve := range ves {
		errString += ve.Error()

		if i < len(ves)-1 {
			errString += " & "
		}
	}

	return errString
}

func (ves ValidatorErrs) MarshalJSON() ([]byte, error) {
	jsonErrs := make([]validatorErrJson, len(ves))
	for i, ve := range ves {
		jsonErrs[i] = validatorErrJson{
			Error: ve.Error(),
		}
	}

	return json.Marshal(validatorErrsJson{
		Error: jsonErrs,
	})
}

func NewFromValidationErrors(data interface{}, err error) ValidatorErrs {
	valErrs := err.(validator.ValidationErrors)
	validationErrs := make(ValidatorErrs, len(valErrs))

	for i, valErr := range valErrs {
		validationErrs[i] = NewFromValidationError(data, valErr)
	}

	return validationErrs
}

package v10

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/validator"
)

type v10Validator struct {
	v *v10.Validate
}

func NewV10() (validator.Validator, error) {
	langEn := en.New()
	langId := id.New()
	uni := ut.New(langEn, langEn, langId)
	trans, _ := uni.GetTranslator("en")

	validate := v10.New()
	if err := enTrans.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	return &v10Validator{
		v: validate,
	}, nil
}

func (v *v10Validator) ValidateStruct(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return err
	}
	return nil
}

func (v *v10Validator) Engine() interface{} {
	return v.v
}

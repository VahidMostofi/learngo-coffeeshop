package data

import (
	"log"
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ValidationError wraps the validators FieldError
type ValidationError struct {
	validator.FieldError
	translator ut.Translator
}

func (v ValidationError) Error() string {
	return v.Translate(v.translator)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

type Validation struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewValidation() *Validation {
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")

	validate := validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		log.Fatalf("failed to register a default translation to validator: %w", err)
	}

	validate.RegisterValidation("sku", validateSKU)

	// validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
	// 	return ut.Add("required", "{0} is a required field.", true)
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("required", fe.Field())
	// 	return t
	// })

	// validate.RegisterTranslation("gt", trans, func(ut ut.Translator) error {
	// 	return ut.Add("gt", "{0} must be greater than {1}.", true)
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("gt", fe.Field(), fe.Param())
	// 	return t
	// })

	validate.RegisterTranslation("sku", translator, func(ut ut.Translator) error {
		return ut.Add("sku", "{0} must follow this regex format: {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("sku", fe.Field(), "'[a-z]+-[a-z]+-[a-z]+'")
		return t
	})
	return &Validation{validate, translator}
}

// Validate an object
func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs {
		ve := ValidationError{err.(validator.FieldError), v.translator}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

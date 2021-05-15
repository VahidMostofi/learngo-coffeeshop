package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ = uni.GetTranslator("en")

	validate = validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
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

	validate.RegisterTranslation("sku", trans, func(ut ut.Translator) error {
		return ut.Add("sku", "{0} must follow this regex format: {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("sku", fe.Field(), "'[a-z]+-[a-z]+-[a-z]+'")
		return t
	})
}

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id of this user
	//
	// required: true
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func TranslateError(e error) map[string]string {
	res := make(map[string]string)
	// return e.(validator.ValidationErrors).Translate(trans)
	for _, e := range e.(validator.ValidationErrors) {
		res[e.Field()] = e.Translate(trans)
	}
	return res
}

func (p *Product) Validate() error {
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found.")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

package data

import (
	"fmt"
)

// ErrProductNotFound is an error raised when a product can't be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found.")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id of this user
	//
	// required: true
	// min: 1
	ID int `json:"id"`

	// name of the product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description of the product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price of the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
}

// Products is a collection of Product
type Products []*Product

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// GetProductByID returns a single product that matches
// with the id filed.
// if no product with this id is found, it returnes an
// Product not found error.
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// UpdateProduct replaces the proeduct with the provided
// id (in product object) by the new provided product.
// if no product with this id is found, it returnes an
// Product not found error.
func UpdateProduct(p *Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = p

	return nil
}

// AddProduct adds a a new product in the database.
func AddProduct(p Product) {
	p.ID = getNextId()
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
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
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}

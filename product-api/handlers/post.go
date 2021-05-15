package handlers

import (
	"net/http"

	"github.com/vahidmostofi/coffeeshop/data"
)

// swagger:route POST /products products createProduct
// Returns new product
// responses:
//	200: productResponse

// AddProduct creates one product on the database
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

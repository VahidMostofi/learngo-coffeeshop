package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vahidmostofi/coffeeshop/data"
)

// swagger:route PUT /products/{id} products updateProduct
// Returns updated product
// responses:
//	200: productResponse

// UpdateProduct updates one product on the database
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid id is provided", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "product no found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}
}

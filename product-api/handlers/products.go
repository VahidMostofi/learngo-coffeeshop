package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/vahidmostofi/coffeeshop/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("Put", r.URL.Path)
		regx := regexp.MustCompile(`/([0-9]+)`)

		g := regx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 || len(g[0]) != 2 {
			p.l.Println("invalid URL", g)
			http.Error(rw, "invalid URL", http.StatusBadRequest)
			return
		}

		idStr := g[0][1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(rw, "invalid id", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, rw, r)
		p.l.Println("got id:", id)
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json.", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	data.UpdateProduct(id, prod)
}

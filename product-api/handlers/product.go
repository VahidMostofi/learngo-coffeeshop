package handlers

import (
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, h *http.Request) {

}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Loukis-13/currency-API-GO/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", routes.Home)
	router.HandleFunc("/moedas", routes.Moedas).Methods("GET")
	router.HandleFunc("/moedas/{moeda}", routes.Moeda).Methods("GET")
	router.HandleFunc("/troca", routes.Troca).Methods("GET")

	router.HandleFunc("/usuarios", routes.Usuarios).Methods("GET")
	router.HandleFunc("/usuario/{id}", routes.Usuario).Methods("GET")
	router.HandleFunc("/usuario/troca", routes.UsuarioTroca).Methods("POST")

	fmt.Println("Servindo na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

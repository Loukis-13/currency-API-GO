package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Loukis-13/currency-API-GO/auth"
	"github.com/Loukis-13/currency-API-GO/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", routes.Home)
	router.HandleFunc("/moedas", routes.Moedas).Methods("GET")
	router.HandleFunc("/moedas/{moeda}", routes.Moeda).Methods("GET")
	router.HandleFunc("/moedas/troca", routes.Troca).Methods("GET")

	router.HandleFunc("/usuarios", routes.Usuarios).Methods("GET")
	router.HandleFunc("/usuario/{id}", routes.Usuario).Methods("GET")
	router.HandleFunc("/usuario/troca", routes.UsuarioTroca).Methods("POST")
	router.HandleFunc("/usuario/troca", routes.ExcluirTroca).Methods("DELETE")

	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/registrar", auth.Registrar).Methods("POST")
	router.HandleFunc("/excluir", auth.Excluir).Methods("DELETE")

	fmt.Println("Servindo na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

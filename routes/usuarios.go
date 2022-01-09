package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Loukis-13/currency-API-GO/database"

	"github.com/gorilla/mux"
)

func Usuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := database.PegarTodosUsuarios()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func Usuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	usuario, err := database.PegarUsuarioPorId(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(usuario)
}

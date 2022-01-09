package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Loukis-13/currency-API-GO/auth"
	"github.com/Loukis-13/currency-API-GO/currencies"
	"github.com/Loukis-13/currency-API-GO/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// STRUCTS

type pedidoTroca struct {
	Moeda1      string  `json:"moeda1"`
	Moeda2      string  `json:"moeda2"`
	Quantidade1 float64 `json:"quantidade1"`
}

type respostaTroca struct {
	Moeda1      string  `json:"moeda1" bson:"moeda1"`
	Moeda2      string  `json:"moeda2" bson:"moeda2"`
	Quantidade1 float64 `json:"quantidade1" bson:"quantidade1"`
	Quantidade2 float64 `json:"quantidade2" bson:"quantidade2"`
	ValorMoeda2 float64 `json:"valorMoeda2" bson:"valorMoeda2"`
	FeitaEm     string  `json:"feitaEm" bson:"feitaEm"`
}

type dataTroca struct {
	FeitaEm string `json:"feitaEm"`
}

func NewRespostaTroca(body pedidoTroca) (respostaTroca, error) {
	valor, err := currencies.Trocar(body.Moeda1, body.Moeda2)
	if err != nil {
		return respostaTroca{}, err
	}

	r := respostaTroca{
		Moeda1:      body.Moeda1,
		Moeda2:      body.Moeda2,
		Quantidade1: body.Quantidade1,
		Quantidade2: valor * body.Quantidade1,
		ValorMoeda2: valor,
		FeitaEm:     time.Now().Format("2006/01/02 15:04:05"),
	}

	return r, nil
}

// ROTAS

func UsuarioTroca(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header não é 'application/json'", http.StatusUnsupportedMediaType)
		return
	}

	idUsuario, err := auth.PegarToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body := pedidoTroca{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Moeda1 == "" {
		http.Error(w, "Necessário passar 'moeda1'", http.StatusBadRequest)
		return
	}
	if body.Moeda2 == "" {
		http.Error(w, "Necessário passar 'moeda2'", http.StatusBadRequest)
		return
	}
	if body.Quantidade1 <= 0 {
		body.Quantidade1 = 1
	}

	troca, err := NewRespostaTroca(body)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	objectId, err := primitive.ObjectIDFromHex(idUsuario)
	if err != nil {
		json.NewEncoder(w).Encode("ID inválida: " + err.Error())
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$push": bson.M{"trocas": bson.M{"$each": []respostaTroca{troca}, "$position": 0}}}

	_, err = database.Usuarios.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Error ao salvar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(troca)
}

func ExcluirTroca(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header não é 'application/json'", http.StatusUnsupportedMediaType)
		return
	}

	body := dataTroca{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idUsuario, err := auth.PegarToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	objectId, err := primitive.ObjectIDFromHex(idUsuario)
	if err != nil {
		json.NewEncoder(w).Encode("ID inválida: " + err.Error())
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$pull": bson.M{"trocas": bson.M{"feitaEm": body.FeitaEm}}}

	_, err = database.Usuarios.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Error ao salvar", http.StatusInternalServerError)
		return
	}

	usuario, err := database.PegarUsuarioPorId(idUsuario)
	if err != nil {
		http.Error(w, "Item excluído, erro ao retornar usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

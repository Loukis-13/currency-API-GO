package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Loukis-13/currency-API-GO/internal/database"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Usuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var result []bson.M

	sensitiveInformation := options.Find().SetProjection(bson.M{"_id": 0, "_class": 0, "password": 0})
	cursor, err := database.Usuarios.Find(context.TODO(), bson.M{}, sensitiveInformation)
	if err == mongo.ErrNoDocuments {
		json.NewEncoder(w).Encode("Usuário não encontrado")
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = cursor.All(context.TODO(), &result)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(result)
}

func Usuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var result bson.M

	objectId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode("ID inválida")
		return
	}

	sensitiveInformation := options.FindOne().SetProjection(bson.M{"_id": 0, "_class": 0, "password": 0})
	err = database.Usuarios.FindOne(context.TODO(), bson.M{"_id": objectId}, sensitiveInformation).Decode(&result)
	if err == mongo.ErrNoDocuments {
		json.NewEncoder(w).Encode("Usuário não encontrado")
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(result)
}

func UsuarioTroca(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header não é 'application/json'", http.StatusUnsupportedMediaType)
		return
	}

	body := pedidoTroca{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
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
	if body.Quantidade1 == 0 {
		body.Quantidade1 = 1
	}

	troca, err := NewRespostaTroca(body)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	objectId, err := primitive.ObjectIDFromHex("61b35c55735aaf1544f41be0")
	if err != nil {
		json.NewEncoder(w).Encode("ID inválida")
		return
	}

	match := bson.M{"_id": objectId}
	add := bson.M{"$push": bson.M{"trocas": bson.M{"$each": []respostaTroca{troca}, "$position": 0}}}

	_, err = database.Usuarios.UpdateOne(context.TODO(), match, add)
	if err != nil {
		http.Error(w, "Error ao salvar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(troca)
}

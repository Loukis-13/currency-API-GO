package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Loukis-13/currency-API-GO/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// STRUCTS

type UsernamePassword struct {
	Id       string `bson:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Usuario struct {
	Username string   `json:"username"`
	Password string   `json:"-"`
	Trocas   []string `json:"trocas"`
}

// ROTAS

func Login(w http.ResponseWriter, r *http.Request) {
	var err error
	var body, user UsernamePassword

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.Usuarios.FindOne(context.TODO(), bson.M{"username": body.Username}).Decode(&user)
	if err != nil || !VerificarSenha(user.Password, body.Password) {
		http.Error(w, "usuario ou senha incorretos", http.StatusUnauthorized)
		return
	}

	validToken, err := GerarToken(map[string]string{"_id": user.Id})
	if err != nil {
		json.NewEncoder(w).Encode("Falha ao gerar token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"Bearer Token": validToken})
}

func Registrar(w http.ResponseWriter, r *http.Request) {
	var err error
	var body UsernamePassword

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	senhaEncriptada, err := GerarSenha(body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usuario := Usuario{Username: body.Username, Password: senhaEncriptada, Trocas: make([]string, 0)}
	_, err = database.Usuarios.InsertOne(context.TODO(), usuario)
	if mongo.IsDuplicateKeyError(err) {
		http.Error(w, "Nome de usuario já em uso", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func Excluir(w http.ResponseWriter, r *http.Request) {
	var err error
	var body, user UsernamePassword

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.Usuarios.FindOne(context.TODO(), bson.M{"username": body.Username}).Decode(&user)
	if err != nil || !VerificarSenha(user.Password, body.Password) {
		http.Error(w, "usuario ou senha incorretos", http.StatusUnauthorized)
		return
	}

	idUser, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = database.Usuarios.DeleteOne(context.TODO(), bson.M{"_id": idUser})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Usuário excluído")
}

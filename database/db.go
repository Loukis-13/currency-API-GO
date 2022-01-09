package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var Usuarios *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func init() {
	fmt.Println("Conectando banco de dados...")

	var err error

	Client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// defer Client.Disconnect(ctx)
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	Usuarios = Client.Database("myFirstDatabase").Collection("users")
}

//

func PegarUsuarioPorId(id string) (result bson.M, erro error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválida")
	}

	sensitiveInformation := options.FindOne().SetProjection(bson.M{"_id": 0, "_class": 0, "password": 0})

	err = Usuarios.FindOne(context.TODO(), bson.M{"_id": objectId}, sensitiveInformation).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("Usuário não encontrado")
	} else if err != nil {
		return nil, err
	}

	return
}

func PegarTodosUsuarios() (result []bson.M, erro error) {
	sensitiveInformation := options.Find().SetProjection(bson.M{"_id": 0, "_class": 0, "password": 0})

	cursor, err := Usuarios.Find(context.TODO(), bson.M{}, sensitiveInformation)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}

	return
}

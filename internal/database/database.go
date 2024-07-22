package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"pokemon-api/internal/models"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PokemonEntry struct {
	*models.PokemonEntry
}

type Service interface {
	Health() map[string]string
	RegisterPokemon(pokemon PokemonEntry) (string, error)
	GetAllPokemonEntries() ([]PokemonEntry, error)
}

type service struct {
	db *mongo.Client
}

var (
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	user = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_ROOT_PASSWORD")
)

func Connect() Service {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
    clientOptions.SetAuth(options.Credential{
        Username: user,
        Password: password,
    })

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) RegisterPokemon(pokemon PokemonEntry) (string, error) {
	collection := s.db.Database("pokedex").Collection("pokemons")

	result, err := collection.InsertOne(context.TODO(), pokemon)
	if err != nil {
		return "", err
	}

	log.Printf("Add doc to mongo, id: %v", result)

	return "Inserted", nil
}

func (s *service) GetAllPokemonEntries()([]PokemonEntry, error) {
	collection := s.db.Database("pokedex").Collection("pokemons")
	
	cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
		return nil, err
    }
	
    defer cursor.Close(context.Background())
	
	var entries []PokemonEntry

	if err := cursor.All(context.Background(), &entries); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
    
	log.Printf("Entries: %v", entries)
    return entries, nil
}

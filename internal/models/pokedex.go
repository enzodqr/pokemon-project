package models

type PokemonEntry struct {
	Name      string   `bson:"name"`
	Types     []string `bson:"types"`
	Abilities []string `bson:"abilities"`
	Evolution string   `bson:"evolution"`
}
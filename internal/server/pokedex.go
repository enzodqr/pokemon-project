// pokemon-api/internal/server/server.go
package server

import (
    "context"
    "fmt"
    "net/http"

    "pokemon-api/internal/handler"
    "pokemon-api/internal/database"
    "pokemon-api/internal/models"
)

type PokemonEntry struct {
	*models.PokemonEntry
}

func (pokemon PokemonEntry) Valid(ctx context.Context) map[string]string {
    problems := make(map[string]string)

    if pokemon.Name == "" {
        problems["Name"] = "Name is required"
    }
    if len(pokemon.Types) == 0 {
        problems["Types"] = "At least one type is required"
    }
    if len(pokemon.Abilities) == 0 {
        problems["Abilities"] = "At least one ability is required"
    }

    return problems
}

func (s *Server) RegisterPokemon(w http.ResponseWriter, r *http.Request) {
    pokemon, problems, err := handler.Decode[PokemonEntry](r)

    if err != nil {
        if len(problems) > 0 {
            http.Error(w, fmt.Sprintf("Validation errors: %v", problems), http.StatusBadRequest)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    result, err := s.db.RegisterPokemon(database.PokemonEntry(pokemon))

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Pokemon created: %s", result)
}

func (s *Server) GetAllPokemonEntries(w http.ResponseWriter, r *http.Request) {
    var entries []models.PokemonEntry
    result, err := s.db.GetAllPokemonEntries()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    for _, pokemon := range result {
        entries = append(entries, models.PokemonEntry{
            Name:      pokemon.Name,
            Types:     pokemon.Types,
            Abilities: pokemon.Abilities,
            Evolution: pokemon.Evolution,
        })
    }

    if err := handler.Encode(w, http.StatusOK, entries); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

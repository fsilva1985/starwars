package swapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fsilva1985/starwars/entity"
)

type Planet interface {
	GetOne(name string) (entity.PlanetSwapi, error)
}

func NewPlanet(client Client) Planet {
	return &planetImpl{
		client: client,
	}
}

type planetImpl struct {
	client Client
}

func (s *planetImpl) GetOne(name string) (entity.PlanetSwapi, error) {
	var planet entity.PlanetSwapi

	req, _ := http.NewRequest(http.MethodGet, s.client.Url()+"planets/?search="+name, nil)

	res, err := s.client.Do(req)
	if err != nil {
		return planet, err
	}

	if res.StatusCode != http.StatusOK {
		return planet, errors.New("Swapi.GetOne: Error http response")
	}

	var planets []entity.PlanetSwapi
	var body body

	body.Results = &planets

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return planet, fmt.Errorf("Swapi.GetOne: decode: %w", err)
	}

	if body.Count == 0 {
		return planet, nil
	}

	return planets[0], nil
}

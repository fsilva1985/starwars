package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/fsilva1985/starwars/entity"
)

const (
	documentType = "planets"
)

type PlanetStorage interface {
	GetAll(ctx context.Context, name string) ([]entity.Planet, error)
	GetOne(ctx context.Context, id string) (entity.Planet, error)
	Create(ctx context.Context, post entity.Planet) error
	Delete(ctx context.Context, id string) error
}

type planetRequestQuery struct {
	Query struct {
		Match struct {
			Name string `json:"name"`
		} `json:"match"`
	} `json:"query"`
}

type planetsResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source entity.Planet `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewPlanetStorage(elasticsearch Elasticsearch) PlanetStorage {
	return &planetStorageImpl{
		elasticsearch: elasticsearch,
		timeout:       time.Second * 10,
	}
}

type planetStorageImpl struct {
	elasticsearch Elasticsearch
	timeout       time.Duration
}

func (p *planetStorageImpl) GetAll(ctx context.Context, name string) ([]entity.Planet, error) {
	var planets []entity.Planet

	req := esapi.SearchRequest{
		Index:        []string{p.elasticsearch.Index()},
		DocumentType: []string{documentType},
	}

	if name != "" {
		planetNameQuery := new(planetRequestQuery)
		planetNameQuery.Query.Match.Name = name

		body, err := json.Marshal(planetNameQuery)
		if err != nil {
			return planets, fmt.Errorf("PlanetStorage.GetAll: marshall: %w", err)
		}

		req.Body = bytes.NewReader(body)
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := p.elasticsearch.Search(ctx, req)
	if err != nil {
		return planets, fmt.Errorf("PlanetStorage.GetAll: request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return planets, fmt.Errorf("PlanetStorage.GetAll: response: %s", res.String())
	}

	var bodyResponse planetsResponse

	if err := json.NewDecoder(res.Body).Decode(&bodyResponse); err != nil {
		return planets, fmt.Errorf("PlanetStorage.GetAll: decode: %w", err)
	}

	posts := make([]entity.Planet, len(bodyResponse.Hits.Hits))
	for i, v := range bodyResponse.Hits.Hits {
		posts[i] = v.Source
	}

	return posts, nil
}

func (p *planetStorageImpl) GetOne(ctx context.Context, id string) (entity.Planet, error) {
	var planet entity.Planet

	req := esapi.GetRequest{
		Index:        p.elasticsearch.Index(),
		DocumentType: documentType,
		DocumentID:   id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := p.elasticsearch.Request(ctx, req)
	if err != nil {
		return planet, fmt.Errorf("PlanetStorage.GetOne: request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return planet, fmt.Errorf("PlanetStorage.GetOne: response: %s", res.String())
	}

	var body document

	body.Source = &planet

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return planet, fmt.Errorf("PlanetStorage.GetOne: decode: %w", err)
	}

	return planet, nil
}

func (p *planetStorageImpl) Create(ctx context.Context, planet entity.Planet) error {
	body, err := json.Marshal(planet)
	if err != nil {
		return fmt.Errorf("PlanetStorage.Create: marshall: %w", err)
	}

	req := esapi.CreateRequest{
		Index:        p.elasticsearch.Index(),
		DocumentType: documentType,
		DocumentID:   planet.ID,
		Body:         bytes.NewReader(body),
		Refresh:      "true",
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := p.elasticsearch.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("PlanetStorage.Create: request: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("PlanetStorage.Create: response: %s", res.String())
	}

	return nil
}

func (p planetStorageImpl) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:        p.elasticsearch.Index(),
		DocumentType: documentType,
		DocumentID:   id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := p.elasticsearch.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("PlanetStorage.Delete: request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("PlanetStorage.Delete: response: %s", res.String())
	}

	return nil
}

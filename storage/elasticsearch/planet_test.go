package elasticsearch_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/fsilva1985/starwars/entity"
	"github.com/fsilva1985/starwars/storage/elasticsearch"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestPlanetGetAll(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientMock := elasticsearch.NewMockElasticsearch(ctrl)
	clientMock.EXPECT().Index().Times(1).Return("starwars")

	json := `{
		"took": 11,
		"timed_out": false,
		"_shards": {
			"total": 1,
			"successful": 1,
			"skipped": 0,
			"failed": 0
		},
		"hits": {
			"total": {
				"value": 9,
				"relation": "eq"
			},
			"max_score": 0.05129329,
			"hits": [{
				"_index": "starwars",
				"_type": "planets",
				"_id": "E7hJJeXr3OSEMq3SE8hvpNi7ki",
				"_score": 0.05129329,
				"_source": {
					"id": "E7hJJeXr3OSEMq3SE8hvpNi7ki",
					"name": "Tatooine",
					"climate": "arid",
					"terrain": "desert",
					"movie_apparition": 5
				}
			}]
		}
	}`

	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	response := &esapi.Response{
		StatusCode: 200,
		Body:       body,
	}

	clientMock.EXPECT().Search(gomock.Any(), gomock.Any()).Times(1).Return(response, nil)

	planetStorage := elasticsearch.NewPlanetStorage(clientMock)
	planetResponse, _ := planetStorage.GetAll(ctx, "Tatooine")

	if diff := cmp.Diff(planetResponse[0].Name, "Tatooine"); diff != "" {
		t.Error(diff)
	}
}

func TestPlanetGetOne(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientMock := elasticsearch.NewMockElasticsearch(ctrl)
	clientMock.EXPECT().Index().Times(1).Return("starwars")

	json := `{
		"_index": "starwars",
		"_type": "planets",
		"_id": "E7hJJeXr3OSEMq3SE8hvpNi7ki",
		"_version": 1,
		"_seq_no": 0,
		"_primary_term": 1,
		"found": true,
		"_source": {
			"id": "E7hJJeXr3OSEMq3SE8hvpNi7ki",
			"name": "Tatooine",
			"climate": "arid",
			"terrain": "desert",
			"movie_apparition": 5
		}
	}`

	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	response := &esapi.Response{
		StatusCode: 200,
		Body:       body,
	}

	clientMock.EXPECT().Request(gomock.Any(), gomock.Any()).Times(1).Return(response, nil)

	planetStorage := elasticsearch.NewPlanetStorage(clientMock)
	planetResponse, _ := planetStorage.GetOne(ctx, "Tatooine")

	if diff := cmp.Diff(planetResponse.Name, "Tatooine"); diff != "" {
		t.Error(diff)
	}
}

func TestPlanetInsert(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientMock := elasticsearch.NewMockElasticsearch(ctrl)
	clientMock.EXPECT().Index().Times(1).Return("starwars")

	response := &esapi.Response{
		StatusCode: 201,
	}

	clientMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(response, nil)

	planetStorage := elasticsearch.NewPlanetStorage(clientMock)
	err := planetStorage.Create(ctx, entity.Planet{})

	if diff := cmp.Diff(err, nil); diff != "" {
		t.Error(diff)
	}
}

func TestPlanetDelete(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientMock := elasticsearch.NewMockElasticsearch(ctrl)
	clientMock.EXPECT().Index().Times(1).Return("starwars")

	response := &esapi.Response{
		StatusCode: 204,
	}

	clientMock.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1).Return(response, nil)

	planetStorage := elasticsearch.NewPlanetStorage(clientMock)
	err := planetStorage.Delete(ctx, "1")

	if diff := cmp.Diff(err, nil); diff != "" {
		t.Error(diff)
	}
}

package swapi_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/fsilva1985/starwars/service/swapi"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestGetOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientMock := swapi.NewMockClient(ctrl)
	clientMock.EXPECT().Url().Times(1).Return("https://swapi.dev/api/")

	json := `{
		"count": 1,
		"next": null,
		"previous": null,
		"results": [
			{
				"name": "Tatooine",
				"rotation_period": "23",
				"orbital_period": "304",
				"diameter": "10465",
				"climate": "arid",
				"gravity": "1 standard",
				"terrain": "desert",
				"surface_water": "1",
				"population": "200000",
				"residents": [
					"https://swapi.dev/api/people/1/",
					"https://swapi.dev/api/people/2/",
					"https://swapi.dev/api/people/4/",
					"https://swapi.dev/api/people/6/",
					"https://swapi.dev/api/people/7/",
					"https://swapi.dev/api/people/8/",
					"https://swapi.dev/api/people/9/",
					"https://swapi.dev/api/people/11/",
					"https://swapi.dev/api/people/43/",
					"https://swapi.dev/api/people/62/"
				],
				"films": [
					"https://swapi.dev/api/films/1/",
					"https://swapi.dev/api/films/3/",
					"https://swapi.dev/api/films/4/",
					"https://swapi.dev/api/films/5/",
					"https://swapi.dev/api/films/6/"
				],
				"created": "2014-12-09T13:50:49.641000Z",
				"edited": "2014-12-20T20:58:18.411000Z",
				"url": "https://swapi.dev/api/planets/1/"
			}
		]
	}`

	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	response := &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	clientMock.EXPECT().Do(gomock.Any()).Times(1).Return(response, nil)

	planet := swapi.NewPlanet(clientMock)
	planetResponse, _ := planet.GetOne("Tatooine")

	if diff := cmp.Diff(planetResponse.Name, "Tatooine"); diff != "" {
		t.Error(diff)
	}
}

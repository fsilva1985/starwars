package main

import (
	"net/http"
	"time"

	"github.com/fsilva1985/starwars/entity"
	"github.com/fsilva1985/starwars/service/swapi"
	"github.com/fsilva1985/starwars/storage/elasticsearch"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

func main() {
	es := elasticsearch.New(
		[]string{
			"http://elasticsearch:9200",
		},
		"starwars",
	)

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	httpClient := &http.Client{Transport: tr}

	r := setupRouter(es, httpClient)
	r.Run(":8080")
}

func setupRouter(es elasticsearch.Elasticsearch, httpClient *http.Client) *gin.Engine {
	planetStorage := elasticsearch.NewPlanetStorage(es)
	store := persistence.NewInMemoryStore(time.Second)

	r := gin.Default()
	r.GET("/planets", cache.CachePage(store, time.Second*5, func(c *gin.Context) {
		planets, err := planetStorage.GetAll(c, c.Query("name"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, planets)
	}))

	r.GET("/planets/:id", cache.CachePage(store, time.Second*5, func(c *gin.Context) {
		planet, err := planetStorage.GetOne(c, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, planet)
	}))

	r.POST("/planets", func(c *gin.Context) {
		var planet entity.Planet

		if err := c.ShouldBindJSON(&planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		planets, err := planetStorage.GetAll(c, planet.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(planets) > 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Planet exists in database"})
			return
		}

		client := swapi.NewClient(httpClient)
		planetSwapi := swapi.NewPlanet(client)
		planetResponse, err := planetSwapi.GetOne(planet.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		planet.MovieApparition = len(planetResponse.Films)
		planet.ID = entity.GetNewId()

		if err := planetStorage.Create(c, planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, planet)
	})

	r.DELETE("/planets/:id", func(c *gin.Context) {
		if err := planetStorage.Delete(c, c.Param("id")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Writer.WriteHeader(http.StatusNoContent)
	})

	return r
}

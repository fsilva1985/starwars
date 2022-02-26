package elasticsearch

import (
	"context"
	"net/http"

	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

type Elasticsearch interface {
	Client() *elasticsearch.Client
	Index() string
	Search(ctx context.Context, request esapi.SearchRequest) (*esapi.Response, error)
	Request(ctx context.Context, request esapi.GetRequest) (*esapi.Response, error)
	Create(ctx context.Context, request esapi.CreateRequest) (*esapi.Response, error)
	Delete(ctx context.Context, request esapi.DeleteRequest) (*esapi.Response, error)
}

func New(addresses []string, index string) Elasticsearch {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	res, err := client.Indices.Exists([]string{index})
	if err != nil {
		panic(err)
	}

	if res.StatusCode == http.StatusNotFound {
		res, err = client.Indices.Create(index)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusOK {
			panic(err)
		}
	}

	return &elasticsearchImpl{
		client: client,
		index:  index,
	}
}

type elasticsearchImpl struct {
	client *elasticsearch.Client
	index  string
}

type document struct {
	Source interface{} `json:"_source"`
}

func (p *elasticsearchImpl) Client() *elasticsearch.Client {
	return p.client
}

func (p *elasticsearchImpl) Index() string {
	return p.index
}

func (p *elasticsearchImpl) Search(ctx context.Context, request esapi.SearchRequest) (*esapi.Response, error) {
	res, err := request.Do(ctx, p.client)

	return res, err
}

func (p *elasticsearchImpl) Request(ctx context.Context, request esapi.GetRequest) (*esapi.Response, error) {
	res, err := request.Do(ctx, p.client)

	return res, err
}

func (p *elasticsearchImpl) Create(ctx context.Context, request esapi.CreateRequest) (*esapi.Response, error) {
	res, err := request.Do(ctx, p.client)

	return res, err
}

func (p *elasticsearchImpl) Delete(ctx context.Context, request esapi.DeleteRequest) (*esapi.Response, error) {
	res, err := request.Do(ctx, p.client)

	return res, err
}

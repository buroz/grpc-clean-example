package storage

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/buroz/grpc-clean-example/pkg/config"
)

type ArangoClient struct {
	config *config.ArangoConfig
	Db     driver.Database
}

func NewArangoClient(conf *config.ArangoConfig) ArangoClient {
	return ArangoClient{
		config: conf,
	}
}

func (c *ArangoClient) Connect(ctx context.Context) error {
	dbUrl := fmt.Sprintf("http://%v:%d", c.config.Host, c.config.Port)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{dbUrl},
	})
	if err != nil {
		return err
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(c.config.User, c.config.Password),
	})
	if err != nil {
		return err
	}

	isExists, err := client.DatabaseExists(ctx, c.config.Name)
	if err != nil {
		return err
	}

	var db driver.Database

	if isExists {
		db, err = client.Database(ctx, c.config.Name)
		if err != nil {
			return err
		}
	} else {
		db, err = client.CreateDatabase(ctx, c.config.Name, nil)
		if err != nil {
			return err
		}
	}

	c.Db = db

	return nil
}

func (c *ArangoClient) Collection(ctx context.Context, collectionName string) (driver.Collection, error) {
	var col driver.Collection

	colExists, err := c.Db.CollectionExists(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	if colExists {
		col, err = c.Db.Collection(ctx, collectionName)
		if err != nil {
			return nil, err
		}
	} else {
		col, err = c.Db.CreateCollection(ctx, collectionName, nil)
		if err != nil {
			return nil, err
		}
	}

	return col, err
}

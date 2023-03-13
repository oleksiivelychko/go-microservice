package database

import (
	"github.com/oleksiivelychko/go-microservice/api"
	mysqlUtils "github.com/oleksiivelychko/go-utils/mysql_connection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectAllProductFromDb(t *testing.T) {
	connection, err := mysqlUtils.NewMySQLConnection("gouser", "secret", "go_microservice")
	if err != nil {
		t.Error(err)
	}

	results, err := connection.SelectAll("products")
	if err != nil {
		t.Error(err)
	}

	var products []*api.Product

	for results.Next() {
		var product api.Product

		err = results.Scan(&product.ID, &product.Name, &product.Price, &product.SKU, &product.UpdatedAt)
		if err != nil {
			t.Error(err)
		}

		products = append(products, &product)
	}

	assert.Equal(t, products[0].ID, 1)
	assert.Equal(t, products[0].Name, "Latte")
	assert.Equal(t, products[0].Price, 1.49)
	assert.Equal(t, products[0].SKU, "123-456-789")
	assert.NotEmpty(t, products[0].UpdatedAt)
	assert.Equal(t, products[1].ID, 2)
	assert.Equal(t, products[1].Name, "Espresso")
	assert.Equal(t, products[1].Price, 0.99)
	assert.Equal(t, products[1].SKU, "000-000-001")
	assert.NotEmpty(t, products[1].UpdatedAt)

	defer results.Close()
	defer connection.Close()
}

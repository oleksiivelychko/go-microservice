package database

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/api"
	mysql "github.com/oleksiivelychko/go-utils/mysql_connection"
	"testing"
)

func TestSelectAllProductFromDb(t *testing.T) {
	connection, err := mysql.NewMySQLConnection("root", "secret", "go_microservice")
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

	for _, product := range products {
		fmt.Println(product.ID)
		fmt.Println(product.Name)
		fmt.Println(product.Price)
		fmt.Println(product.SKU)
		fmt.Println(product.UpdatedAt)
	}

	defer results.Close()
	defer connection.Close()
}

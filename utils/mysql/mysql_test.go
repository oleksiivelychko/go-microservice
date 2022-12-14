package mysql

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/sdk/models"
	"github.com/oleksiivelychko/go-microservice/utils"
	"os"
	"testing"
)

func TestNewMySQLConnection(t *testing.T) {
	os.Setenv("MYSQL_USERNAME", "root")
	os.Setenv("MYSQL_PASSWORD", "secret")
	os.Setenv("MYSQL_DATABASE", "go_microservice")

	connection, err := NewMySQLConnection(utils.NewLogger())
	if err != nil {
		t.Error(err)
	}

	results, err := connection.fetchAll()
	if err != nil {
		t.Error(err)
	}

	var products []*models.Product

	for results.Next() {
		var product models.Product

		err = results.Scan(&product.ID, &product.Name, &product.Price, &product.SKU, &product.UpdatedAt)
		if err != nil {
			t.Error(err)
		}

		products = append(products, &product)
	}

	for _, product := range products {
		fmt.Println(product.ID)
		fmt.Println(*product.Name)
		fmt.Println(*product.Price)
		fmt.Println(*product.SKU)
		fmt.Println(product.UpdatedAt)
	}

	defer results.Close()
	defer connection.db.Close()
}

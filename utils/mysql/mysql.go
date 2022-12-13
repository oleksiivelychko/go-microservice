package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-hclog"
	"os"
	"time"
)

type Connection struct {
	db     *sql.DB
	logger hclog.Logger
}

func NewMySQLConnection(l hclog.Logger) (*Connection, error) {
	var mysqlUsername = os.Getenv("MYSQL_USERNAME")
	var mysqlPassword = os.Getenv("MYSQL_PASSWORD")
	var mysqlDatabase = os.Getenv("MYSQL_DATABASE")
	var dataSourceName = fmt.Sprintf("%s:%s@/%s", mysqlUsername, mysqlPassword, mysqlDatabase)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	/*
		It's required to ensure connections are closed by the driver safely before connection is closed by MySQL server, OS, or other middlewares.
		Since some middlewares close idle connections by 5 minutes, we recommend timeout shorter than 5 minutes.
		This setting helps load balancing and changing system variables too.
	*/
	db.SetConnMaxLifetime(time.Minute * 3)
	/*
		It's highly recommended to limit the number of connection used by the application.
		There is no recommended limit number because it depends on application and MySQL server.
	*/
	db.SetMaxOpenConns(10)
	/*
		It's recommended to be set same to db.SetMaxOpenConns().
		When it is smaller than SetMaxOpenConns(), connections can be opened and closed much more frequently than you expect.
		Idle connections can be closed by the db.SetConnMaxLifetime().
		If you want to close idle connections more rapidly, you can use db.SetConnMaxIdleTime() since Go 1.15.
	*/
	db.SetMaxIdleConns(10)

	return &Connection{db, l}, nil
}

func (connection *Connection) fetchAll() (*sql.Rows, error) {
	return connection.db.Query("SELECT * from products;")
}

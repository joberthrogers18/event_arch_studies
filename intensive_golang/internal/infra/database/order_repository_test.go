package database

import (
	"database/sql"

	"github.com/stretchr/testify/suite"
	"testing"
)

// initialize the suite test
type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, fina_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

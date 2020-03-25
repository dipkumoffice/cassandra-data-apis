package db

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/inf.v0"
	"reflect"
	"testing"
)

func (suite *IntegrationTestSuite) TestNumericCqlTypeMapping() {
	queries := []string{
		"CREATE TABLE ks1.tbl_numerics (id int PRIMARY KEY, bigint_value bigint, float_value float," +
			" double_value double, smallint_value smallint, tinyint_value tinyint, decimal_value decimal)",
		"INSERT INTO ks1.tbl_numerics (id, bigint_value, float_value, double_value, smallint_value, tinyint_value" +
			", decimal_value) VALUES (1, 1, 1.1, 1.1, 1, 1, 1.25)",
		"INSERT INTO ks1.tbl_numerics (id) VALUES (100)",
	}

	for _, query := range queries {
		err := suite.db.session.Execute(query, nil)
		assert.Nil(suite.T(), err)
	}

	var (
		rs  ResultSet
		err error
		row map[string]interface{}
	)

	rs, err = suite.db.session.ExecuteIter("SELECT * FROM ks1.tbl_numerics WHERE id = ?", nil, 1)
	assert.Nil(suite.T(), err)
	row = rs.Values()[0]
	assertPointerValue(suite.T(), new(string), "1", row["bigint_value"])
	assertPointerValue(suite.T(), new(float32), float32(1.1), row["float_value"])
	assertPointerValue(suite.T(), new(float64), 1.1, row["double_value"])
	assertPointerValue(suite.T(), new(int16), int16(1), row["smallint_value"])
	assertPointerValue(suite.T(), new(int8), int8(1), row["tinyint_value"])
	assertPointerValue(suite.T(), new(*inf.Dec), inf.NewDec(125, 2), row["decimal_value"])

	// Assert nil values
	rs, err = suite.db.session.ExecuteIter("SELECT * FROM ks1.tbl_numerics WHERE id = ?", nil, 100)
	assert.Nil(suite.T(), err)
	row = rs.Values()[0]
	assertNilPointer(suite.T(), new(string), row["bigint_value"])
	assertNilPointer(suite.T(), new(float32), row["float_value"])
	assertNilPointer(suite.T(), new(float64), row["double_value"])
	assertNilPointer(suite.T(), new(int16), row["smallint_value"])
	assertNilPointer(suite.T(), new(int8), row["tinyint_value"])
	assert.IsType(suite.T(), new(*inf.Dec), row["decimal_value"])
	assert.Nil(suite.T(), reflect.ValueOf(row["decimal_value"]).Elem().Interface())
}

func assertPointerValue(t *testing.T, expectedType interface{}, expected interface{}, actual interface{}) {
	assert.IsType(t, expectedType, actual)
	assert.Equal(t, expected, reflect.ValueOf(actual).Elem().Interface())
}

func assertNilPointer(t *testing.T, expectedType interface{}, actual interface{}) {
	assert.IsType(t, expectedType, actual)
	assert.Nil(t, actual)
}

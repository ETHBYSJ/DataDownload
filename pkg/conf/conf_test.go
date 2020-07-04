package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestConfParse(t *testing.T) {
	asserts := assert.New(t)
	Init("../../conf/app.ini")
	// database
	err := mapSection("database", DatabaseConfig)
	asserts.NoError(err)
	asserts.Equal(DatabaseConfig.Type, "mysql")
	asserts.Equal(DatabaseConfig.User, "root")
	asserts.Equal(DatabaseConfig.Password, "19961013")
	asserts.Equal(DatabaseConfig.Host, "127.0.0.1")
	asserts.Equal(DatabaseConfig.Name, "file_manager")
	asserts.Equal(DatabaseConfig.TablePrefix, "v1_")
	asserts.Equal(DatabaseConfig.Port, 3306)

	// system
	err = mapSection("system", SystemConfig)
	asserts.Equal(SystemConfig.Debug, true)


}



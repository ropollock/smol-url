package data

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDBProvider struct {
	db *gorm.DB
}

type PostgresDBProviderInterface interface {
	GetDB() *gorm.DB
	Connect(dbURI string)
}

func PostgresDBProvider() *postgresDBProvider {
	return &postgresDBProvider{}
}

func (provider *postgresDBProvider) Connect(dbURI string) {
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	provider.db = db
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Postgres successfully connected.")
}

func (provider *postgresDBProvider) GetDB() *gorm.DB {
	return provider.db
}

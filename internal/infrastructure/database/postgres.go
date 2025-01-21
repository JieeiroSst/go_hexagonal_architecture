package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type PostgresDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func NewPostgresDB(masterConfig, slaveConfig DBConfig) (*PostgresDB, error) {
	masterDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		masterConfig.Host, masterConfig.Port, masterConfig.User, masterConfig.Password, masterConfig.DBName)

	slaveDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		slaveConfig.Host, slaveConfig.Port, slaveConfig.User, slaveConfig.Password, slaveConfig.DBName)

	masterDB, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master DB: %v", err)
	}

	slaveDB, err := gorm.Open(postgres.Open(slaveDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to slave DB: %v", err)
	}

	// Configure connection pools
	sqlDB, err := masterDB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	sqlDB, err = slaveDB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return &PostgresDB{
		Master: masterDB,
		Slave:  slaveDB,
	}, nil
}

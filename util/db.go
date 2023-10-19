package util

import (
	common_log "common/log"
	"database/sql"
	"fmt"
	"go-hj-hospital/config"
	"log"
	syslog "log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

type Patients struct {
	ID               uint   `gorm:"primaryKey"`
	Nik              string `gorm:"unique;not null"`
	Name             string `gorm:"not null"`
	Birthdate        time.Time
	Sex              string `gorm:"not null"`
	Address          string
	Phone            string
	EmergencyContact string
}

type Doctors struct {
	ID             uint   `gorm:"primaryKey"`
	EmployeeID     string `gorm:"unique;not null"`
	Name           string `gorm:"not null"`
	Specialization string
	Phone          string
	WorkDay        string
}

type Queue struct {
	ID               uint   `gorm:"primaryKey"`
	PatientNik       string `gorm:"not null"`
	DoctorEmployeeID string `gorm:"not null"`
	QueueDate        time.Time
	Status           string `gorm:"default:Belum dilayani"`
}

func Master() *gorm.DB {
	return db
}

func CreateDB() {
	gormConf := &gorm.Config{}
	newLogger := logger.New(
		syslog.New(common_log.GetLogger().GetWriter(), "\r\n[db]", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	gormConf.Logger = newLogger
	err := OpenDB(config.Instance.MySqlUrl, gormConf, config.Instance.MySqlMaxIdle, config.Instance.MySqlMaxOpen)
	if err != nil {
		common_log.Error(err)
		panic(err)
	}
	common_log.Info("SQL Connection established")
}

func CloseMasterDB() {
	if db == nil {
		return
	}

	if err := sqlDB.Close(); nil != err {
		common_log.Errorf("Disconnect from database failed: %v", err.Error())
	}
}

func OpenDB(dsn string, config *gorm.Config, maxIdleConns, maxOpenConns int) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		}
	}

	if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		fmt.Printf("opens database failed: %v", err.Error())
		return
	}
	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		common_log.Error(err)
	}
	return
}

func MigrateDB(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&Patients{}, &Doctors{}, &Queue{})
	if err != nil {
		return err
	}
	common_log.Info("Migration successful")
	return nil
}

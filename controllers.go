package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func initializeDb() (*gorm.DB, *sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Helsinki", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	pgConf := postgres.Config{DSN: dsn, PreferSimpleProtocol: true}
	db, err := gorm.Open(postgres.New(pgConf), &gorm.Config{})
	if err != nil {
		for {
			db, err = gorm.Open(postgres.New(pgConf), &gorm.Config{})
			if err == nil {
				break
			}
			log.Printf("database connection failed, trying again")
			time.Sleep(time.Duration(10))
		}
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(1)
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		return nil, nil, err
	}
	log.Printf("connected to db: %v", db.ConnPool)
	return db, sqlDB, err
}

func (a *App) searchTodos() ([]APITodo, error) {
	var todos []APITodo
	result := a.DB.Find(&Todo{}).Where("done", false).Limit(30).Find(&todos)
	if result.Error != nil {
		return []APITodo{}, result.Error
	}
	log.Printf("getting %d todos from db", result.RowsAffected)

	return todos, nil
}

func (a *App) insertTodo(newTodo Todo) error {
	l := len(newTodo.Task)
	if l > 0 && l < 141 {
		result := a.DB.Create(&newTodo)
		if result.Error != nil {
			return result.Error
		}
		log.Printf("created a new todo %s", newTodo.Task)
		err := publish2NATS(newTodo)
		if err != nil {
			log.Printf("publishing to NATS %s channel failed: %s", os.Getenv("NATS_CHANNEL"), err)
		}
		return nil
	} else {
		return errors.New("task length should be between [1,140] characters")
	}
}

func (a *App) updateTodoDone(id int) error {
	log.Printf("trying to mark todo id: %d done", id)
	result := a.DB.Find(&Todo{}, id).Where("done", false)
	if result.RowsAffected == 1 {
		result = a.DB.Find(&Todo{}, id).Update("done", true)
		if result.RowsAffected == 1 {
			var todo Todo
			result = a.DB.Find(&Todo{}, id).First(&todo)
			if result.Error != nil {
				return result.Error
			}
			err := publish2NATS(todo)
			if err != nil {
				log.Printf("publishing to NATS %s channel failed: %s", os.Getenv("NATS_CHANNEL"), err)
			}
		}

		return nil
	}
	return errors.New("malformed request")
}

func publish2NATS(t Todo) error {
	j, err := json.Marshal(t)
	if err != nil {
		log.Printf("marshalling todo failed: %s", err)
	}
	var message = string(j[:len(j)])
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return err
	}
	err = nc.Publish(os.Getenv("NATS_CHANNEL"), []byte(message))
	if err != nil {
		return err
	}
	return nil
}

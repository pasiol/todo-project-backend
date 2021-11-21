package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func initializeDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Helsinki", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	pgConf := postgres.Config{DSN: dsn, PreferSimpleProtocol: true}
	db, err := gorm.Open(postgres.New(pgConf), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func (a *App) searchTodos() ([]APITodo, error) {
	var todos []APITodo
	result := a.DB.Model(&Todo{}).Limit(30).Find(&todos)
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
		return nil
	} else {
		return errors.New("task length should be between [1,140] characters")
	}
}

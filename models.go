package main

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Task string `json:"task"`
}

type APITodo struct {
	Task string `json:"task"`
}

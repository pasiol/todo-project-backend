package main

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type APITodo struct {
	Id   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

package main

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Task string `gorm:"check:task_checker,task <> ''"`
}

type APITodo struct {
	Task string `json:"task"`
}

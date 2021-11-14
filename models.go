package main

type Todo struct {
	Task string `bson:"task" json:"task"`
}

var Todos []Todo

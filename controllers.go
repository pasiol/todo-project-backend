package main

func searchTodos() ([]Todo, error) {
	return Todos, nil
}

func insertTodo(newTask Todo) error {
	Todos = append(Todos, newTask)
	return nil
}

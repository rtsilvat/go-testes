package main

import (
    "fmt"
    "net/http"
    "go-tasks-api/handlers" // Importa o pacote handlers
    "go-tasks-api/database"
)

func main() {

    database.Init()

    // Associa as rotas às funções exportadas pelo pacote handlers
    http.Handle("/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.ListTasks)))
    //http.HandleFunc("/tasks", handlers.ListTasks)         // GET /tasks
    http.HandleFunc("/login", handlers.Login)       // POST /login/login
    http.HandleFunc("/tasks/add", handlers.AddTask)       // POST /tasks/add
    http.HandleFunc("/tasks/update", handlers.UpdateTask) // PUT /tasks/update?id={id}
    http.HandleFunc("/tasks/delete", handlers.DeleteTask) // DELETE /tasks/delete?id={id}

    fmt.Println("Servidor rodando na porta 8080...")
    http.ListenAndServe(":8080", nil) // Inicia o servidor na porta 8080
}


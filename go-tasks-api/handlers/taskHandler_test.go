package handlers

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "go-tasks-api/database"
    "go-tasks-api/models"
)

func setupTestDB() {
    // Configura um banco de dados em memória para testes
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("Falha ao configurar banco de testes: " + err.Error())
    }

    // Migra a estrutura da tabela para o banco de dados em memória
    db.AutoMigrate(&models.Task{})

    // Adiciona tarefas iniciais para testes
    db.Create(&models.Task{Title: "Tarefa 1", Description: "Descrição 1", Completed: false})
    db.Create(&models.Task{Title: "Tarefa 2", Description: "Descrição 2", Completed: true})

    // Substitui o banco global pela instância de teste
    database.DB = db
}

func TestListTasks(t *testing.T) {
    setupTestDB() // Configura o banco de dados para testes

    req, _ := http.NewRequest("GET", "/tasks", nil)
    res := httptest.NewRecorder()

    handler := http.HandlerFunc(ListTasks)
    handler.ServeHTTP(res, req)

    if res.Code != http.StatusOK {
        t.Errorf("Esperado status 200, mas recebeu %d", res.Code)
    }
}

func TestAddTask(t *testing.T) {
    setupTestDB() // Configura o banco de dados para testes

    cases := []struct {
        name     string
        input    string
        expected int
    }{
        {"Caso válido", `{"title": "Tarefa 1"}`, http.StatusCreated},
        {"Dados inválidos", `{"invalid": "data"}`, http.StatusBadRequest},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            req, _ := http.NewRequest("POST", "/tasks/add", strings.NewReader(tc.input))
            res := httptest.NewRecorder()

            handler := http.HandlerFunc(AddTask)
            handler.ServeHTTP(res, req)

            if res.Code != tc.expected {
                t.Errorf("Esperado %d, recebido %d", tc.expected, res.Code)
            }
        })
    }
}

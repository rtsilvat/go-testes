package handlers

import (
	 "fmt"
    "encoding/json"
    "net/http"
    "strconv"
    "go-tasks-api/models" // Importa o pacote models
    "go-tasks-api/database"
    "github.com/go-playground/validator/v10"
    "time"
    "github.com/golang-jwt/jwt/v5"
)


var validate = validator.New()
var jwtKey = []byte("minha_chave_secreta123")

// GenerateJWT generates a new JWT token for the given user ID.
// The token expires in 24 hours.
// Parameters: 
// - userID: The ID of the user to generate the token for.
// Returns:
// - string: The generated JWT token.
func GenerateJWT(userID uint) (string, error) {
    claims := &jwt.MapClaims{
        "userID": userID,
        "exp":    time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}


// AuthMiddleware is a middleware function that checks for a valid JWT token in the
// Authorization header of the incoming HTTP request. If the token is valid, the request
// is passed to the next handler in the chain. If the token is invalid or missing, it
// responds with an HTTP 401 Unauthorized status.
//
// Parameters:
// - next: The next http.Handler to be called if the token is valid.
//
// Returns:
// - http.Handler: A handler that wraps the original handler with JWT authentication.
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        claims := &jwt.MapClaims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Não autorizado", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func Login(w http.ResponseWriter, r *http.Request) {
    var newUser models.User

    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }

    
    if err := validate.Struct(&newUser); err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }
    token, _ := GenerateJWT(newUser.ID);
    
    // if err := database.DB.Create(&newTask).Error; err != nil {
    //     http.Error(w, "Erro ao salvar tarefa", http.StatusInternalServerError)
    //     return
    // }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(token)
}

// ListTasks lista todas as tarefas
func ListTasks(w http.ResponseWriter, r *http.Request) {

    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

    fmt.Println("Page: ", page)
    fmt.Println("Limit: ", limit)
    if page <= 0 {
        page = 1
    }
    if limit <= 0 {
        limit = 10
    }


    var tasks []models.Task
    offset := (page - 1) * limit

    // Com Paginação
    if err := database.DB.Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
        http.Error(w, "Erro ao buscar tarefas", http.StatusInternalServerError)
        return
    }

    // SEM PAGINAÇÃO
    // if err := database.DB.Find(&tasks).Error; err != nil {
    //     http.Error(w, "Erro ao buscar tarefas", http.StatusInternalServerError)
    //     return
    // }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

// AddTask adiciona uma nova tarefa
func AddTask(w http.ResponseWriter, r *http.Request) {
    var newTask models.Task

    if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }

    if err := validate.Struct(&newTask); err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.DB.Create(&newTask).Error; err != nil {
        http.Error(w, "Erro ao salvar tarefa", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newTask)
}

// UpdateTask atualiza o status de uma tarefa
func UpdateTask(w http.ResponseWriter, r *http.Request) {

    var updateTask models.Task

    if err := json.NewDecoder(r.Body).Decode(&updateTask); err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }


    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    fmt.Println("Update ", id)

    // Converte o corpo da requisição para o tipo Task
    var task models.Task
    if err := database.DB.First(&task, id).Error; err != nil {
        http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
        return
    }

    task.Completed = !task.Completed
    task.Title = updateTask.Title
    if err := database.DB.Save(&task).Error; err != nil {
        http.Error(w, "Erro ao atualizar tarefa", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(task)
}

// DeleteTask deleta uma tarefa
func DeleteTask(w http.ResponseWriter, r *http.Request) {

    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := database.DB.Delete(&models.Task{}, id).Error; err != nil {
        http.Error(w, "Erro ao deletar tarefa", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)    
}


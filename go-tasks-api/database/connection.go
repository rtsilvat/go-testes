package database
import (
    "log"
	"github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"fmt"
	"os"
)

var DB *gorm.DB

// Inicializa a conexão com o banco
func Init() {
    var err error

    // Carregar variáveis do .env
    if err := godotenv.Load(); err != nil {
        log.Println("Erro ao carregar o arquivo .env, utilizando variáveis do ambiente do sistema")
    }

	dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_SSLMODE"),
    )
    
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }

    log.Println("Conexão com o banco de dados estabelecida!")
}

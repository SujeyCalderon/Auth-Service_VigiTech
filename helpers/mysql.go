package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)


func ConnectMySQL() (*sql.DB, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	if host == "" || port == "" || username == "" || database == "" {
		return nil, fmt.Errorf("faltan variables de entorno para la conexión MySQL")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión MySQL: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error al hacer ping a MySQL: %w", err)
	}
	return db, nil
}
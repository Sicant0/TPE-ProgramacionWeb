package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	db "proyectoweb/db/sqlc"

	_ "github.com/lib/pq"
)

var (
	dbConn  *sql.DB
	queries *db.Queries
)

func main() {
	// Las variables de entorno las carga docker compose desde el .env
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("PORT_IN"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// sql.Open no se conecta, solo valida los parámetros y prepara el pool de conexiones
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error en la configuración de la base de datos: %v", err)
	}
	defer dbConn.Close()

	// Ping es la primera conexión real: si la base no responde, el servidor no arranca
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos (Ping fallido): %v", err)
	}
	log.Println("Conexión exitosa a la base de datos")

	queries = db.New(dbConn)

	// Todas las rutas están en web_handlers.go
	registrarRutasWeb()

	port := ":8080"
	log.Printf("Servidor escuchando en el puerto %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

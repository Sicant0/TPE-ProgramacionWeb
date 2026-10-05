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
	http.HandleFunc("/usuario", UsuarioHandler)
	http.HandleFunc("/usuario/", UsuarioHandler)
	http.HandleFunc("/gasto", GastoHandler)
	http.HandleFunc("/gasto/", GastoHandler)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("PORT_IN"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	var err error
	dbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConn.Close()

	//esto solo esta para cuando se corre por docker compose
	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Error conectando a la base de datos (Ping fallido): %v", err)
	}
	log.Println("¡Conexión exitosa a la base de datos!")

	// Instancio el repositorio
	queries = db.New(dbConn)

	port := ":8080"
	fmt.Printf("Servidor escuchando en el puerto %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}

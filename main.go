package main

import (
	"fmt"
	"net/http"
)

func main() {
	//define ruta hacia directorio static
	staticDir := "./static"

	//crea fileServer apuntando hacia el directorio
	fileServer := http.FileServer(http.Dir(staticDir))

	//crea handle para root hacia fileServer
	http.Handle("/", fileServer)
	port := ":8080"
	fmt.Printf("Servidor estatico escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	//inicia servidor exponiendo port (:8080)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}

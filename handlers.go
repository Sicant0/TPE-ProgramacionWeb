package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	db "proyectoweb/db/sqlc"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// tipos para handling de requests con anotaciones para "encoding/json"
type UsuarioRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GastoRequest struct {
	IDUser      int32     `json:"id_user"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Monto       string    `json:"monto"`
	Categoria   string    `json:"categoria"`
	Fecha       time.Time `json:"fecha"`
}

// map de categorias validas, utilizado para chequear antes de convertir string categoria a db.PosiblesCategorias
var categoriasValidas = map[string]bool{
	"Comida":          true,
	"Transporte":      true,
	"Hogar":           true,
	"Salud":           true,
	"Entretenimiento": true,
	"Ropa":            true,
	"Regalos":         true,
	"Servicios":       true,
	"Otros":           true,
}

// validar interfaz
type Validable interface {
	Validar() error
}

// funciones de ayuda
func esEmailValido(email string) bool {
	// Chequeo simple: un único "@", con algo antes y un "." después.
	at := strings.Index(email, "@")
	if at <= 0 || at == len(email)-1 {
		return false
	}
	dominio := email[at+1:]
	return strings.Contains(dominio, ".") && !strings.HasPrefix(dominio, ".") && !strings.HasSuffix(dominio, ".")
}

func esMontoValido(monto string) error {
	if monto == "" {
		return errors.New("el monto es requerido")
	}
	montoConvertido, err := strconv.ParseFloat(monto, 64)
	if err != nil {
		return errors.New("el monto debe ser numérico")
	}
	if montoConvertido <= 0 {
		return errors.New("el monto debe ser mayor a cero")
	}
	return nil
}

func extraerID(path string) (int32, error) {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 || parts[1] == "" {
		return 0, errors.New("falta el id en la URL")
	}
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, errors.New("el id en la URL debe ser un entero válido")
	}
	return int32(id), nil
}

// Validaciones
func (req GastoRequest) Validar() error {
	if req.IDUser <= 0 {
		return errors.New("id_user es requerido y debe ser mayor a 0")
	}
	titulo := strings.TrimSpace(req.Titulo)
	if titulo == "" || utf8.RuneCountInString(titulo) > 15 {
		return errors.New("titulo es requerido y debe tener hasta 15 caracteres")
	}
	descripcion := strings.TrimSpace(req.Descripcion)
	if descripcion == "" || utf8.RuneCountInString(descripcion) > 255 {
		return errors.New("descripcion es requerida y debe tener hasta 255 caracteres")
	}
	if err := esMontoValido(req.Monto); err != nil {
		return err
	}
	if !categoriasValidas[req.Categoria] {
		return errors.New("categoria inválida")
	}
	if req.Fecha.IsZero() {
		return errors.New("fecha es requerida")
	}
	return nil
}

func (req UsuarioRequest) Validar() error {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if n := utf8.RuneCountInString(req.Username); n < 3 || n > 25 {
		return errors.New("username debe tener entre 3 y 25 caracteres")
	}
	if utf8.RuneCountInString(req.Email) > 50 {
		return errors.New("email debe tener hasta 50 caracteres")
	}
	if !esEmailValido(req.Email) {
		return errors.New("email inválido")
	}
	return nil
}

// Decodificar y validar
// esta funcion recibe el request
func decodificarYValidar(r *http.Request, dst Validable) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("body inválido: " + err.Error())
	}
	return dst.Validar()
}

// esta funcion recibe el responseWriter junto con el status y el resultado del validador para escribir codificar el resultado en JSON
func escribirJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		err := json.NewEncoder(w).Encode(v)
		if err != nil {
			return
		}
	}
}

// Handlers de Usuario y Gasto
func UsuarioHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsuario(w, r)
	case http.MethodPost:
		createUsuario(w, r)
	case http.MethodDelete:
		deleteUsuario(w, r)
	case http.MethodPut:
		updateUsuario(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GastoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getGasto(w, r)
	case http.MethodPost:
		createGasto(w, r)
	case http.MethodDelete:
		deleteGasto(w, r)
	case http.MethodPut:
		updateGasto(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Funciones de Usuario
func createUsuario(w http.ResponseWriter, r *http.Request) {
	var req UsuarioRequest
	if err := decodificarYValidar(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.CreateUsuarioParams{
		Username: strings.TrimSpace(req.Username),
		Email:    strings.TrimSpace(req.Email),
	}

	usuario, err := queries.CreateUsuario(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear usuario", http.StatusInternalServerError)
		return
	}

	escribirJSON(w, http.StatusCreated, usuario)
}

func deleteUsuario(w http.ResponseWriter, r *http.Request) {
	idUser, err := extraerID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//check if user has any gastos
	gastos, err := queries.GetGastosByUsuario(r.Context(), idUser)
	if err != nil {
		http.Error(w, "error al leer gastos del usuario", http.StatusInternalServerError)
		return
	}
	if len(gastos) > 0 {
		http.Error(w, "el usuario tiene gastos asociados", http.StatusConflict)
		return
	}
	//delete usuario if all checks pass
	if err := queries.DeleteUsuario(r.Context(), idUser); err != nil {
		http.Error(w, "error al borrar usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUsuario(w http.ResponseWriter, r *http.Request) {
	if strings.Trim(r.URL.Path, "/") == "usuario" { //get all users
		usuarios, err := queries.GetAllUsuarios(r.Context())
		if err != nil {
			http.Error(w, "error al leer usuarios", http.StatusInternalServerError)
			return
		}
		escribirJSON(w, http.StatusOK, usuarios)
		return
	}
	//get single user
	idUser, err := extraerID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	usuario, err := queries.GetUsuario(r.Context(), idUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	escribirJSON(w, http.StatusOK, usuario)
}

func updateUsuario(w http.ResponseWriter, r *http.Request) {
	idUser, err := extraerID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req UsuarioRequest
	if err := decodificarYValidar(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.UpdateUsuarioParams{
		IDUser:   idUser,
		Username: strings.TrimSpace(req.Username),
		Email:    strings.TrimSpace(req.Email),
	}

	usuario, err := queries.UpdateUsuario(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	escribirJSON(w, http.StatusOK, usuario)
}

func createGasto(w http.ResponseWriter, r *http.Request) {
	var req GastoRequest
	if err := decodificarYValidar(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//check user exists
	if _, err := queries.GetUsuario(r.Context(), req.IDUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "el usuario indicado no existe", http.StatusBadRequest)
			return
		}
		http.Error(w, "error al leer el usuario", http.StatusInternalServerError)
		return
	}
	//user exists, create db request
	params := db.CreateGastoParams{
		IDUser:      req.IDUser,
		Titulo:      strings.TrimSpace(req.Titulo),
		Descripcion: strings.TrimSpace(req.Descripcion),
		Monto:       req.Monto,
		Categoria:   db.PosiblesCategorias(req.Categoria),
		Fecha:       req.Fecha,
	}

	gasto, err := queries.CreateGasto(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear gasto", http.StatusInternalServerError)
		return
	}

	escribirJSON(w, http.StatusCreated, gasto)

}

func deleteGasto(w http.ResponseWriter, r *http.Request) {
	idGasto, err := extraerID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := queries.DeleteGasto(r.Context(), idGasto); err != nil {
		http.Error(w, "error al borrar gasto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func getGasto(w http.ResponseWriter, r *http.Request) {
	if strings.Trim(r.URL.Path, "/") == "gasto" { //get all gastos
		gastos, err := queries.GetAllGastos(r.Context())
		if err != nil {
			http.Error(w, "error al leer gastos", http.StatusInternalServerError)
			return
		}
		escribirJSON(w, http.StatusOK, gastos)
		return
	}
	//get single gasto
	idGasto, err := extraerID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gasto, err := queries.GetGasto(r.Context(), idGasto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "gasto no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	escribirJSON(w, http.StatusOK, gasto)
}

func updateGasto(w http.ResponseWriter, r *http.Request) {
	idGasto, err := extraerID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req GastoRequest
	if err := decodificarYValidar(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.UpdateGastoParams{
		IDGasto:     idGasto,
		IDUser:      req.IDUser,
		Titulo:      strings.TrimSpace(req.Titulo),
		Descripcion: strings.TrimSpace(req.Descripcion),
		Monto:       req.Monto,
		Categoria:   db.PosiblesCategorias(req.Categoria),
		Fecha:       req.Fecha,
	}

	gasto, err := queries.UpdateGasto(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "gasto no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "error al actualizar gasto", http.StatusInternalServerError)
		return
	}

	escribirJSON(w, http.StatusOK, gasto)
}

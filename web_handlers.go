package main

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"net/http"
	db "proyectoweb/db/sqlc"
	"proyectoweb/views"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/a-h/templ"
)

// formatoFecha es el formato que envía y espera <input type="date"> para los formularios HTML.
const formatoFecha = "2006-01-02"

// registrarRutasWeb tiene todas las rutas hacia las paginas renderizadas mediante templ
func registrarRutasWeb() {
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("GET /{$}", indexPage)

	http.HandleFunc("GET /usuarios", listarUsuariosPage)
	http.HandleFunc("POST /usuarios", crearUsuarioForm)
	http.HandleFunc("GET /usuarios/{id}", mostrarUsuarioItem)
	http.HandleFunc("GET /usuarios/{id}/editar", editarUsuarioItem)
	http.HandleFunc("POST /usuarios/{id}/editar", actualizarUsuarioForm)
	http.HandleFunc("DELETE /usuarios/{id}", eliminarUsuario)

	http.HandleFunc("GET /gastos", listarGastosPage)
	http.HandleFunc("POST /gastos", crearGastoForm)
	http.HandleFunc("GET /gastos/{id}", mostrarGastoItem)
	http.HandleFunc("GET /gastos/{id}/editar", editarGastoItem)
	http.HandleFunc("POST /gastos/{id}/editar", actualizarGastoForm)
	http.HandleFunc("DELETE /gastos/{id}", eliminarGasto)
}

// VALIDACION

type UsuarioRequest struct {
	Username string
	Email    string
}

type GastoRequest struct {
	IDUser      int32
	Titulo      string
	Descripcion string
	Monto       string
	Categoria   string
	Fecha       time.Time
}

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
		return errors.New("El monto debe ser numérico")
	}
	if montoConvertido <= 0 {
		return errors.New("El monto debe ser mayor a cero")
	}
	return nil
}

func (req UsuarioRequest) Validar() error {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if n := utf8.RuneCountInString(req.Username); n < 3 || n > 25 {
		return errors.New("El nombre de usuario debe tener entre 3 y 25 caracteres")
	}
	if utf8.RuneCountInString(req.Email) > 50 {
		return errors.New("El email debe tener hasta 50 caracteres")
	}
	if !esEmailValido(req.Email) {
		return errors.New("El email debe tener el formato nombre@dominio.com")
	}
	return nil
}

func (req GastoRequest) Validar() error {
	if req.IDUser <= 0 {
		return errors.New("id_user es requerido y debe ser mayor a 0")
	}
	titulo := strings.TrimSpace(req.Titulo)
	if titulo == "" || utf8.RuneCountInString(titulo) > 15 {
		return errors.New("El titulo es requerido y debe tener hasta 15 caracteres")
	}
	descripcion := strings.TrimSpace(req.Descripcion)
	if descripcion == "" || utf8.RuneCountInString(descripcion) > 255 {
		return errors.New("La descripcion es requerida y debe tener hasta 255 caracteres")
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

// FUNCIONES AUXILIARES (utilizadas por las funciones de rendering de paginas)

func render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	var buf bytes.Buffer
	if err := c.Render(r.Context(), &buf); err != nil {
		log.Printf("error al renderizar: %v", err)
		http.Error(w, "error al renderizar la página", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

func idDeRuta(r *http.Request) (int32, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("el id en la URL debe ser un entero válido")
	}
	return int32(id), nil
}

// SECCION DE RENDERING DE PAGINAS

func indexPage(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, views.IndexPage())
}

// rendering de paginas de usuarios

// Renderiza la pagina de usuarios
func listarUsuariosPage(w http.ResponseWriter, r *http.Request) {
	usuarios, err := queries.GetAllUsuarios(r.Context())
	if err != nil {
		http.Error(w, "error al leer usuarios", http.StatusInternalServerError)
		return
	}
	render(w, r, http.StatusOK, views.UsersPage(usuarios, views.UsuarioFormData{}))
}

// leerYValidadFormUsuario lee el formulario y arma el request para crear/actualizar un usuario
func leerYValidarFormUsuario(r *http.Request) (UsuarioRequest, views.UsuarioFormData, error) {
	if err := r.ParseForm(); err != nil {
		return UsuarioRequest{}, views.UsuarioFormData{}, errors.New("formulario inválido")
	}
	req := UsuarioRequest{
		Username: strings.TrimSpace(r.PostFormValue("username")),
		Email:    strings.TrimSpace(r.PostFormValue("email")),
	}
	form := views.UsuarioFormData{Username: req.Username, Email: req.Email}
	return req, form, req.Validar()
}

// Recibe el form de creacion de un usuario, valida y si es exitoso, crea el usuario y renderiza sobre la misma pagina con el nuevo usuario listado
func crearUsuarioForm(w http.ResponseWriter, r *http.Request) {
	req, form, err := leerYValidarFormUsuario(r)
	if err != nil {
		form.Error = err.Error()
		render(w, r, http.StatusBadRequest, views.UserCreateDivForm(form))
		return
	}

	params := db.CreateUsuarioParams{Username: req.Username, Email: req.Email}
	if _, err := queries.CreateUsuario(r.Context(), params); err != nil {
		form.Error = "error al crear usuario"
		render(w, r, http.StatusInternalServerError, views.UserCreateDivForm(form))
		return
	}

	//En vez de redireccionar a /usuarios, se reconsultan los usuarios para realizar la recarga parcial
	usuarios, err := queries.GetAllUsuarios(r.Context())
	if err != nil {
		http.Error(w, "error al leer usuarios", http.StatusInternalServerError)
		return
	}

	render(
    	w,
    	r,
    	http.StatusOK,
    	templ.Join(
        	views.UserList(usuarios),
        	views.UserCreateDivFormOOB(views.UsuarioFormData{}),
    	),
	)
}

// Funcion auxiliar que busca el usuario de una ruta (el {id} de la URL), utilizado para la edicion/cancelacion de edicion de un usuario
func usuarioDeRuta(w http.ResponseWriter, r *http.Request) (db.Usuario, bool){
	id, err := idDeRuta(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return db.Usuario{}, false
	}
	u, err := queries.GetUsuario(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "error al leer usuario", http.StatusInternalServerError)
		}
		return db.Usuario{}, false
	}
	return u, true
}

// Devuelve la fila de un usuario (utilizada al cancelar una edicion de un usuario)
func mostrarUsuarioItem(w http.ResponseWriter, r *http.Request){
	u, ok := usuarioDeRuta(w, r)
	if !ok {
		return
	}
	render(w, r, http.StatusOK, views.UserItem(u))
}

// Devuelve la fila de un usuario en modo edicion (con formulario precargado con info del usuario)
func editarUsuarioItem(w http.ResponseWriter, r *http.Request) {
	u, ok := usuarioDeRuta(w, r)
	if !ok {
		return
	}
	form := views.UsuarioFormData{ID: u.IDUser, Username: u.Username, Email: u.Email}
	render(w, r, http.StatusOK, views.UserEditItem(form))
}

// Procesa la actualizacion de un usuario: si es valida, retorna la fila editada en modo lectura; si no es valida, devuelve la fila en modo edicion con el error indicado
func actualizarUsuarioForm(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//valida cambios sobre usuario
	req, form, err := leerYValidarFormUsuario(r)
	form.ID = id
	if err != nil {
		form.Error = err.Error()
		render(w, r, http.StatusBadRequest, views.UserEditItem(form))
		return
	}
	//actualiza usuario
	params := db.UpdateUsuarioParams{IDUser: id, Username: req.Username, Email: req.Email}
	u, err := queries.UpdateUsuario(r.Context(), params)
	//si hay errores, renderiza el form de edicion con el error apropiado
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			form.Error = "el usuario ya no existe"
			render(w, r, http.StatusNotFound, views.UserEditItem(form))
			return
		}
		form.Error = "error al actualizar usuario"
		render(w, r, http.StatusInternalServerError, views.UserEditItem(form))
		return
	}
	//sin errores, renderiza la fila del usuario ya editado
	render(w, r, http.StatusOK, views.UserItem(u))
}

// Procesa la eliminacion de un usuario, si es exitosa, elimina al usuario y su fila en el render, eliminadolo de la lista (via hx-delete)
func eliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := queries.DeleteUsuario(r.Context(), id); err != nil {
		http.Error(w, "error al borrar usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) //ahora solo se envia el mensaje de OK, no hay redireccion debido a HTMX
}

// rendering de paginas de gastos

// Renderiza la pagina de gastos, permite editar o eliminar los gastos
func listarGastosPage(w http.ResponseWriter, r *http.Request) {
	gastos, err := queries.GetAllGastos(r.Context())
	if err != nil {
		http.Error(w, "error al leer gastos", http.StatusInternalServerError)
		return
	}
	//rebusca los usuarios nuevamente para mandarselos al createForm que ahora vive en la misma pagina
	usuarios, err := queries.GetAllUsuarios(r.Context())
	if err != nil {
		http.Error(w, "error al leer usuarios", http.StatusInternalServerError)
		return
	}
	form := views.GastoFormData{Fecha: time.Now().Format(formatoFecha)}
	render(w, r, http.StatusOK, views.GastosPage(gastos, form, usuarios))
}

// Funcion auxiliar para poder reemplazar el div del form con el error via hx-target-error
func renderErrorFormGasto(w http.ResponseWriter, r *http.Request, status int, form views.GastoFormData) {
	usuarios, err := queries.GetAllUsuarios(r.Context())
	if err != nil {
		http.Error(w, "error al leer usuarios", http.StatusInternalServerError)
		return
	}
	render(w, r, status, views.GastoCreateDivForm(form, usuarios))
}

// Funcion auxiliar que obtiene el gasto de la URL
func gastoDeRuta(w http.ResponseWriter, r *http.Request) (db.Gasto, bool) {
	id, err := idDeRuta(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return db.Gasto{}, false
	}
	g, err := queries.GetGasto(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "error al leer gasto", http.StatusInternalServerError)
		}
		return db.Gasto{}, false
	}
	return g, true
}

// Devuelve la fila de gasto, utilizado si se cancela la edicion de un gasto
func mostrarGastoItem(w http.ResponseWriter, r *http.Request) {
	g, ok := gastoDeRuta(w, r)
	if !ok {
		return
	}
	render(w, r, http.StatusOK, views.GastoItem(g))
}

// Devuelve un gasto en modo edicion, con los datos precargados del gasto en el formulario
func editarGastoItem(w http.ResponseWriter, r *http.Request) {
	g, ok := gastoDeRuta(w, r)
	if !ok {
		return
	}
	form := views.GastoFormData{
		ID:          g.IDGasto,
		IDUser:      g.IDUser,
		Titulo:      g.Titulo,
		Descripcion: g.Descripcion,
		Monto:       g.Monto,
		Categoria:   string(g.Categoria),
		Fecha:       g.Fecha.Format(formatoFecha),
	}
	render(w, r, http.StatusOK, views.GastoEditItem(form))
}

// leerYValidarFormGasto lee el formulario y arma el request para crear/actualizar un gasto
func leerYValidarFormGasto(r *http.Request) (GastoRequest, views.GastoFormData, error) {
	if err := r.ParseForm(); err != nil {
		return GastoRequest{}, views.GastoFormData{}, errors.New("formulario inválido")
	}
	form := views.GastoFormData{
		Titulo:      strings.TrimSpace(r.PostFormValue("titulo")),
		Descripcion: strings.TrimSpace(r.PostFormValue("descripcion")),
		Monto:       strings.TrimSpace(r.PostFormValue("monto")),
		Categoria:   r.PostFormValue("categoria"),
		Fecha:       r.PostFormValue("fecha"),
	}

	idUser, err := strconv.ParseInt(r.PostFormValue("id_user"), 10, 32)
	if err != nil || idUser <= 0 {
		return GastoRequest{}, form, errors.New("tenés que elegir un usuario") //revisar el comentario del error mas adelante
	}
	form.IDUser = int32(idUser)

	fecha, err := time.Parse(formatoFecha, form.Fecha)
	if err != nil {
		return GastoRequest{}, form, errors.New("fecha inválida")
	}

	req := GastoRequest{
		IDUser:      form.IDUser,
		Titulo:      form.Titulo,
		Descripcion: form.Descripcion,
		Monto:       form.Monto,
		Categoria:   form.Categoria,
		Fecha:       fecha,
	}
	return req, form, req.Validar()
}

// Recibe el form de creacion de un gasto, valida y si es exitoso, crea el gasto y renderiza sobre la misma pagina con el nuevo gasto listado
func crearGastoForm(w http.ResponseWriter, r *http.Request) {
	req, form, err := leerYValidarFormGasto(r)
	if err != nil {
		form.Error = err.Error()
		renderErrorFormGasto(w, r, http.StatusBadRequest, form)
		return
	}
	if _, err := queries.GetUsuario(r.Context(), req.IDUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			form.Error = "el usuario indicado no existe"
			renderErrorFormGasto(w, r, http.StatusBadRequest, form)
			return
		}
		form.Error = "error al leer el usuario"
		renderErrorFormGasto(w, r, http.StatusInternalServerError, form)
		return
	}

	params := db.CreateGastoParams{
		IDUser:      req.IDUser,
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Monto:       req.Monto,
		Categoria:   db.PosiblesCategorias(req.Categoria),
		Fecha:       req.Fecha,
	}
	if _, err := queries.CreateGasto(r.Context(), params); err != nil {
		form.Error = "error al crear gasto"
		renderErrorFormGasto(w, r, http.StatusInternalServerError, form)
		return
	}

	//nueva parte con recarga parcial
	gastos, err := queries.GetAllGastos(r.Context())
	if err != nil {
		http.Error(w, "error al leer gastos", http.StatusInternalServerError)
		return
	}
	usuarios, err := queries.GetAllUsuarios(r.Context())
	if err != nil {
		http.Error(w, "error al leer usuarios", http.StatusInternalServerError)
		return
	}
	render(w, r, http.StatusOK, templ.Join(
		views.GastoList(gastos),
		views.GastoCreateDivFormOOB(views.GastoFormData{Fecha: time.Now().Format(formatoFecha)}, usuarios),
	))
}

// Procesa la actualizacion de un gasto: si es valida, retorna la fila editada en modo lectura; si no es valida, devuelve la fila en modo edicion con el error indicado
func actualizarGastoForm(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//valida cambios sobre gasto
	req, form, err := leerYValidarFormGasto(r)
	form.ID = id
	if err != nil {
		form.Error = err.Error()
		render(w, r, http.StatusBadRequest, views.GastoEditItem(form))
		return
	}
	//actualiza gasto
	params := db.UpdateGastoParams{
		IDGasto:     id,
		IDUser:      req.IDUser,
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Monto:       req.Monto,
		Categoria:   db.PosiblesCategorias(req.Categoria),
		Fecha:       req.Fecha,
	}
	g, err := queries.UpdateGasto(r.Context(), params)
	//si hay errores al actualizar, renderiza el form con el error de actualizacion
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			form.Error = "el gasto ya no existe"
			render(w, r, http.StatusNotFound, views.GastoEditItem(form))
			return
		}
		form.Error = "error al actualizar gasto"
		render(w, r, http.StatusInternalServerError, views.GastoEditItem(form))
		return
	}
	//sin errores, renderiza la fila del gasto ya editado
	render(w, r, http.StatusOK, views.GastoItem(g))
}

// Procesa la eliminacion de un gasto, si es exitosa, elimina al gasto y su fila en el render, eliminadolo de la lista (via hx-delete)
func eliminarGasto(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := queries.DeleteGasto(r.Context(), id); err != nil {
		http.Error(w, "error al borrar gasto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) //ahora solo se envia el mensaje de OK, no hay redireccion debido a HTMX
}
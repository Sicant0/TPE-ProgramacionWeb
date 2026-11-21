package views

import "fmt"

// Structs que utilizan los formularios y tambien para mostrar al usuario lo escrito si la validacion del
// formulario falla

type UsuarioFormData struct {
	ID       int32
	Username string
	Email    string
	Error    string
}

type GastoFormData struct {
	ID          int32
	IDUser      int32
	Titulo      string
	Descripcion string
	Monto       string
	Categoria   string
	Fecha       string
	Error       string
}

var Categorias = []string{
	"Comida",
	"Transporte",
	"Hogar",
	"Salud",
	"Entretenimiento",
	"Ropa",
	"Regalos",
	"Servicios",
	"Otros",
}

// Funcion auxiliar utilizada al editar usuario/gasto in-place. Genera los distintos id de cada formulario de creacion o edicion de una entidad, ya que pueden convivir muchos activos en un mismo momento.
func campoID(nombre string, id int32) string {
	if id == 0 {
		return nombre
	}
	return fmt.Sprintf("%s-%d", nombre, id)
}
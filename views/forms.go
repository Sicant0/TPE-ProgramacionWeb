package views

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

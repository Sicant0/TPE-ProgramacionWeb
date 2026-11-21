# TPE-ProgramacionWeb

Administrador de gastos: una aplicación web para registrar usuarios y sus gastos, categorizarlos y mantener un registro de ellos.

## Ejecución rápida

```bash
git clone https://github.com/Sicant0/TPE-ProgramacionWeb.git
cd TPE-ProgramacionWeb
git checkout tp6
make run-app
```

La aplicación queda disponible en <http://localhost:8080>. Para detenerla, presionar **Ctrl+C**: se borran automáticamente los contenedores y volúmenes creados.

Para correr las pruebas automáticas: `make test`.

## Información sobre persistencia

La aplicación usa PostgreSQL con dos modelos:

- **Usuario** (IDUser, Username, Email): el username tiene entre 3 y 25 caracteres y el email hasta 50.
- **Gasto** (IDGasto, IDUser, Titulo, Descripcion, Monto, Categoria, Fecha): el título tiene hasta 15 caracteres, la descripción hasta 255 y el monto tiene que ser mayor a 0.

Cada gasto pertenece a un usuario. Al borrar un usuario se borran también todos sus gastos (`ON DELETE CASCADE`). La categoría es un tipo enumerado con valores fijos: Comida, Transporte, Hogar, Salud, Entretenimiento, Ropa, Regalos, Servicios y Otros.

## Arquitectura y funcionamiento

La aplicación se renderiza del lado del servidor: el servidor Go genera el HTML de cada página con componentes [templ](https://templ.guide/) y accede a la base de datos mediante el código generado por [sqlc](https://sqlc.dev/). La base de datos y el servidor corren en contenedores de Docker.

- `db/`: esquema de la base de datos (`schema/schema.sql`), consultas SQL (`queries/queries.sql`) y código Go generado por sqlc (`sqlc/`)
- `views/`: componentes templ (`layout.templ`, `index.templ`, `user.templ`, `gasto.templ`) y structs de los formularios (`forms.go`)
- `static/`: hoja de estilos
- `main.go`: conexión a la base de datos y arranque del servidor
- `web_handlers.go`: rutas, validación y handlers HTTP que consultan la base y renderizan las páginas o fragmentos
- `hurl/`: pruebas automáticas del recorrido HTML

### Interactividad con HTMX

La interactividad se resuelve con [HTMX](https://htmx.org/) y su extensión oficial `response-targets`, sin JavaScript propio: el proyecto no contiene archivos `.js`. Ambas se cargan desde un CDN (unpkg), por lo que el navegador necesita conexión a internet.

Las páginas `/usuarios` y `/gastos` muestran el formulario de creación arriba y el listado abajo. Todas las operaciones actualizan solo el fragmento de la página que cambia, sin recargarla:

| Acción | Petición | El servidor responde | Dónde se inserta |
|---|---|---|---|
| Crear | `POST /usuarios` | La lista actualizada y un formulario vacío marcado con `hx-swap-oob` | La lista reemplaza a la anterior (`hx-target`). El formulario vacío reemplaza al formulario usado, ubicado por su `id` (*out-of-band swap*). |
| Editar | `GET /usuarios/{id}/editar` | La fila con un formulario precargado | Reemplaza la fila (`closest li`) |
| Guardar edición | `POST /usuarios/{id}/editar` | La fila actualizada | Reemplaza la fila |
| Cancelar edición | `GET /usuarios/{id}` | La fila original | Reemplaza la fila |
| Borrar | `DELETE /usuarios/{id}` | Respuesta vacía con `200 OK` | Reemplaza la fila por nada, es decir, la elimina. Antes se pide confirmación. |

Gastos funciona igual, con las rutas `/gastos/...`.

### Validación y manejo de errores

Los datos se validan en el navegador, con los atributos de los formularios HTML (`required`, `minlength`, `maxlength`, `min`, `type="email"`), y de nuevo en el servidor, después de quitar los espacios sobrantes.

Si la validación del servidor o la base de datos fallan, el servidor devuelve el formulario con el mensaje de error y los valores ingresados, con código 4xx (validación) o 5xx (error interno). La extensión `response-targets` usa ese código para mostrar la respuesta en un lugar distinto al de una operación exitosa (`hx-target-error`):

- **Al crear:** el error se muestra sobre el formulario de creación y la lista no cambia.
- **Al editar:** el error se muestra dentro de la fila que se está editando, que sigue en modo edición.

## Información sobre ejecución

### Requisitos

- Go (desarrollado en Go 1.26)
- sqlc (instrucciones en <https://docs.sqlc.dev/en/latest/overview/install.html>)
- Docker y Docker Compose (instrucciones en <https://docs.docker.com/get-started/>)

templ no requiere instalación aparte: se ejecuta con `go tool templ`.

### Ejecutar la aplicación: `make run-app`

Realiza todo el ciclo de ejecución con un único comando:

1. **Tareas previas:**
   - Borra contenedores y volúmenes de ejecuciones anteriores.
   - Crea `.env` a partir de `.env.sample` si no existe.
   - Genera el código de templ y sqlc.
   - Compila la aplicación (dentro de la imagen de Docker).
   - Levanta los contenedores y espera a que la base de datos esté lista para recibir conexiones.
2. **Ejecución:** deja la aplicación corriendo en primer plano, mostrando los logs del servidor.
3. **Tareas posteriores:** al presionar Ctrl+C borra contenedores y volúmenes. La limpieza también se hace si el servidor se detiene o si falla cualquiera de los pasos previos.

Cada ejecución arranca con la base de datos vacía.

### Pruebas automáticas: `make test`

Sigue los mismos pasos previos que `make run-app`, corre las pruebas de la carpeta `hurl/` contra la aplicación y al terminar borra contenedores y volúmenes, aunque alguna prueba falle. El comando termina con error si alguna prueba no pasa.

Las pruebas usan [Hurl](https://hurl.dev/), que se ejecuta en un contenedor (no requiere instalación). Recorren las operaciones de usuarios (`usuarios.hurl`) y de gastos (`gastos.hurl`) y verifican:

- los códigos de estado HTTP, incluidos los casos de error (400 para datos o ids inválidos, 404 para registros inexistentes);
- el contenido HTML de las páginas;
- que las respuestas de las operaciones HTMX contengan los fragmentos esperados: la lista actualizada y el formulario *out-of-band* al crear, el formulario con el error al fallar la validación, la fila en modo edición o lectura al editar, y la respuesta vacía al borrar.

### Otros comandos (desarrollo)

| Comando | Qué hace |
|---|---|
| `make server-up` | Crea `.env` si no existe, genera el código de templ y sqlc, y levanta la base de datos y el servidor en segundo plano. Si el arranque falla, muestra los logs y borra contenedores y volúmenes. |
| `make server-down` | Baja los contenedores, conservando los datos cargados. |
| `make server-reset` | Baja los contenedores y borra el volumen de la base de datos. |
| `make generate` | Solo genera el código de templ y sqlc. |
| `make make-env` | Crea `.env` a partir de `.env.sample` si no existe. |
# TPE-ProgramacionWeb

## Información sobre persistencia

Esta aplicación registra gastos de usuarios, permitiendo categorizarlos y mantener un registro de ellos. Para ello existen los modelos de Usuario (IDUser, Username, Email) y Gasto (IDGasto, IDUser, Titulo, Descripcion, Monto, Categoria, Fecha).

## Arquitectura y funcionamiento

La aplicación se renderiza del lado del servidor: el servidor Go genera el HTML de cada página con componentes [templ](https://templ.guide/) y accede a la base de datos PostgreSQL mediante el código generado por [sqlc](https://sqlc.dev/).

- `db/`: codigo de esquema de base de datos y codigo generado de sqlc (`/queries/queries.sql`, `/schema/schema.sql`, `/sqlc`)
- `static/`: hoja de estilos
- `views/`: componentes templ (`layout.templ`, `index.templ`, `user.templ`, `gasto.templ`)
- `main.go` y `web_handlers.go`: handlers HTTP que consultan la base y renderizan las páginas

Los formularios siguen el patrón **Post/Redirect/Get**: cada formulario hace un POST y el servidor, después de guardar los cambios, redirige a la lista correspondiente, que se vuelve a pedir con GET y muestra los datos actualizados.

## Información sobre ejecución

### Requisitos
- Tener instalado Go (desarrollado en Go-1.26)
- Tener instalado sqlc (archivo e instrucciones en <https://sqlc.dev/>)
- Tener instalado Docker y Docker Compose (instrucciones en <https://docs.docker.com/get-started/>)

templ no requiere instalación aparte: se ejecuta con `go tool templ`.

### Comandos

| Comando | Qué hace |
|---|---|
| `make server-up` | Crea `.env` a partir de `.env.sample` si no existe, genera el código de templ y sqlc, y levanta la base de datos y el servidor. Si el arranque falla, muestra los logs y baja los contenedores. |
| `make server-down` | Baja los contenedores, conservando los datos cargados. |
| `make server-reset` | Baja los contenedores y borra el volumen de la base de datos. |
| `make generate` | Solo genera el código de templ y sqlc. |
| `make env` | Crea `.env` a partir de `.env.sample` si no existe. |


Una vez levantado, la aplicación está disponible en <http://localhost:8080>.
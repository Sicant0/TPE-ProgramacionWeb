# TPE-ProgramacionWeb

## Informacion sobre persistencia

Esta aplicacion planea registrar gastos de usuarios, permitiendo categorizar distintos gastos y mantener registro de ellos. Para ello, existen los modelos de Usuario (IDUser, Username, Email) y Gasto (IDGasto, IDUser, Titulo, Descripcion, Monto, Categoria, Fecha).

## Informacion sobre ejecucion
### Requisitos
- Tener instalado Go (desarrollado en Go-1.27.1)
- Tener instalado sqlc (archivo e instrucciones en <https://sqlc.dev/>)
- Tener instalado Docker & Docker Compose (instrucciones en <https://docs.docker.com/get-started/>)

### Endpoints de API
- Para crear un usuario o gasto utilice "/usuario" y "/gasto" respectivamente, utilizando metodo PUT en el request
- Para listar todos los usuarios o gastos utilice "/usuario" y "/gasto" respectivamente, utilizando el metodo GET en el request
- Los usuarios y gastos pueden ser consultados para actualizacion o eliminacion utilizando el endpoint "/usuario/{id}" y "/gasto/{id}" con los metodos PUT y DELETE respectivamente en el request.
 
### Ejecucion de tests de forma automatica

Para crear archivo ".env" a partir de ".env.sample" si aun no existe, levantar los contenedores de testing y base de datos, correr los tests y luego dar de baja estos contenedores, simplemente debera ejecutar el siguiente comando ejecutando el Makefile presente en el directorio: ` make hurl-test `

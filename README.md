# TPE-ProgramacionWeb

## Informacion sobre persistencia

Esta aplicacion planea registrar gastos de usuarios, permitiendo categorizar distintos gastos y mantener registro de ellos. Para ello, existen los modelos de Usuario (IDUser, Username, Email) y Gasto (IDGasto, IDUser, Titulo, Descripcion, Monto, Categoria, Fecha).

## Informacion sobre ejecucion
### Requisitos
- Tener instalado Go (desarrollado en Go-1.27.1)
- Tener instalado sqlc (archivo e instrucciones en <https://sqlc.dev/>)
- Tener instalado Docker & Docker Compose (instrucciones en <https://docs.docker.com/get-started/>)
 
### Ejecucion de tests de forma automatica

Para crear archivo ".env" a partir de ".env.sample" si aun no existe, levantar los contenedores de testing y base de datos, correr los tests y luego dar de baja estos contenedores, simplemente debera ejecutar el siguiente comando ejecutando el Makefile presente en el directorio: ` make test `

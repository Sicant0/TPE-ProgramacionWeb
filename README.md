# TPE-ProgramacionWeb

## Informacion sobre persistencia

Esta aplicacion planea registrar gastos de usuarios, permitiendo categorizar distintos gastos y mantener registro de ellos. Para ello, existen los modelos de Usuario (IDUser, Username, Email) y Gasto (IDGasto, IDUser, Titulo, Descripcion, Monto, Categoria, Fecha).

## Informacion sobre ejecucion

### Ejecucion de tests de forma automatica

1. Copiar el archivo '.env.sample' y crear un archivo '.env'. Para hacerlo ejecutar el comando `cp .env.sample .env`. Puede cambiar su contenido si asi lo desea.
2. Para levantar los contenedores de testing y base de datos, correr los tests y luego dar de baja estos contenedores, simplemente debera ejecutar el siguiente comando ejecutando el Makefile presente en el directorio:<br>
` make test ` <br>

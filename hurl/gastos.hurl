# Pruebas de gastos.
# Se asume la base vacia al empezar (usuarios.hurl la deja vacia al terminar).
# El archivo deja la base vacia al terminar.


# Sin usuarios no se puede cargar un gasto

GET {{base_url}}/gastos
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
xpath "//p[@class='empty']" exists
xpath "//div[@id='gasto-form']" not exists
xpath "count(//ul[@id='gasto-list']/li)" == 0


# Preparacion: un usuario duenio de los gastos
# Obtiene el id_user del usuario creado a traves del hx-get del boton editar del unico user en la lista de users en este momento

POST {{base_url}}/usuarios
[FormParams]
username: inaki2026
email: inakiuv@gmail.com
HTTP 200
[Captures]
id_inaki: xpath "string(//ul[@id='user-list']/li[1]//button[text()='Editar']/@hx-get)" regex "/usuarios/(\\d+)/editar"


# Pagina Inicial

# Con un usuario cargado, aparece el form con el usuario en el select
# Chequea que existe el form de creacion de gasto y la opcion de seleccionar el unico user, ademas de que la lista de gastos esta vacia
GET {{base_url}}/gastos
HTTP 200
[Asserts]
xpath "//div[@id='gasto-form']/form[@hx-post='/gastos']" exists
xpath "string(//div[@id='gasto-form']/form/@hx-target)" == "#gasto-list"
xpath "string(//div[@id='gasto-form']/form/@hx-target-error)" == "#gasto-form"
xpath "string(//select[@name='id_user']/option[@value='{{id_inaki}}'])" == "inaki2026"
xpath "string(//input[@name='fecha']/@value)" matches "^\\d{4}-\\d{2}-\\d{2}$"
xpath "count(//ul[@id='gasto-list']/li)" == 0


# Creacion

# Crear: responde la lista con el gasto nuevo + el form vacio out-of-band
POST {{base_url}}/gastos
[FormParams]
id_user: {{id_inaki}}
titulo: Supermercado
descripcion: Compra semanal
monto: 15000.50
fecha: 2026-09-24
categoria: Comida
HTTP 200
[Captures]
id_super: xpath "string(//ul[@id='gasto-list']/li[1]//button[text()='Editar']/@hx-get)" regex "/gastos/(\\d+)/editar"
[Asserts]
xpath "count(//ul[@id='gasto-list']/li)" == 1
xpath "string(//ul[@id='gasto-list']/li[1]//span[@class='item-title'])" == "Supermercado"
xpath "string(//ul[@id='gasto-list']/li[1]//span[@class='tag'])" == "Comida"
xpath "normalize-space(//ul[@id='gasto-list']/li[1]//span[@class='amount'])" == "$ 15000.50"
xpath "normalize-space(//ul[@id='gasto-list']/li[1]//span[@class='item-meta'])" == "24/09/2026 · Usuario #{{id_inaki}}"
xpath "//div[@id='gasto-form'][@hx-swap-oob='true']" exists
xpath "string(//div[@id='gasto-form']//input[@name='titulo']/@value)" == ""
xpath "count(//div[@id='gasto-form']//select[@name='id_user']/option)" == 2

# Crear un segundo gasto
POST {{base_url}}/gastos
[FormParams]
id_user: {{id_inaki}}
titulo: Colectivo
descripcion: Carga de SUBE
monto: 5000
fecha: 2026-09-25
categoria: Transporte
HTTP 200
[Asserts]
xpath "count(//ul[@id='gasto-list']/li)" == 2


# Errores de validacion en la creacion (400)

# Titulo de mas de 15 caracteres
POST {{base_url}}/gastos
[FormParams]
id_user: {{id_inaki}}
titulo: Titulo demasiado largo
descripcion: Prueba
monto: 100
fecha: 2026-09-25
categoria: Otros
HTTP 400
[Asserts]
xpath "//div[@id='gasto-form'][not(@hx-swap-oob)]" exists
xpath "//div[@id='gasto-form']//p[@class='alert']" exists
xpath "count(//select[@name='id_user']/option)" == 2
xpath "string(//textarea[@name='descripcion'])" == "Prueba"
xpath "//ul[@id='gasto-list']" not exists

# Monto cero
POST {{base_url}}/gastos
[FormParams]
id_user: {{id_inaki}}
titulo: Prueba
descripcion: Prueba
monto: 0
fecha: 2026-09-25
categoria: Otros
HTTP 400
[Asserts]
xpath "//div[@id='gasto-form']//p[@class='alert']" exists

# Categoria que no esta en la lista
POST {{base_url}}/gastos
[FormParams]
id_user: {{id_inaki}}
titulo: Prueba
descripcion: Prueba
monto: 100
fecha: 2026-09-25
categoria: Viajes
HTTP 400
[Asserts]
xpath "//div[@id='gasto-form']//p[@class='alert']" exists

# Fecha invalida
POST {{base_url}}/gastos
[FormParams]
id_user: {{id_inaki}}
titulo: Prueba
descripcion: Prueba
monto: 100
fecha: 2026-13-40
categoria: Otros
HTTP 400
[Asserts]
xpath "//div[@id='gasto-form']//p[@class='alert']" exists

# Sin usuario elegido
POST {{base_url}}/gastos
[FormParams]
id_user:
titulo: Prueba
descripcion: Prueba
monto: 100
fecha: 2026-09-25
categoria: Otros
HTTP 400
[Asserts]
xpath "//div[@id='gasto-form']//p[@class='alert']" exists

# Usuario inexistente
POST {{base_url}}/gastos
[FormParams]
id_user: 99999
titulo: Prueba
descripcion: Prueba
monto: 100
fecha: 2026-09-25
categoria: Otros
HTTP 400
[Asserts]
xpath "//div[@id='gasto-form']//p[@class='alert']" exists

# Los errores no crearon nada
GET {{base_url}}/gastos
HTTP 200
[Asserts]
xpath "count(//ul[@id='gasto-list']/li)" == 2


# Edicion in-place

# Abrir la edicion: devuelve solo la fila con el form precargado
GET {{base_url}}/gastos/{{id_super}}/editar
HTTP 200
[Asserts]
xpath "count(//li)" == 1
xpath "//li/form[@hx-post='/gastos/{{id_super}}/editar']" exists
xpath "string(//li/form/@hx-target)" == "closest li"
xpath "string(//li/form/@hx-target-error)" == "closest li"
xpath "string(//input[@type='hidden'][@name='id_user']/@value)" == "{{id_inaki}}"
xpath "string(//input[@name='titulo']/@value)" == "Supermercado"
xpath "string(//input[@name='monto']/@value)" == "15000.50"
xpath "string(//input[@name='fecha']/@value)" == "2026-09-24"
xpath "string(//select[@name='categoria']/option[@selected]/@value)" == "Comida"
xpath "//input[@id='titulo-{{id_super}}']" exists

# Cancelar: devuelve la fila original en modo lectura
GET {{base_url}}/gastos/{{id_super}}
HTTP 200
[Asserts]
xpath "//li/form" not exists
xpath "string(//li//span[@class='item-title'])" == "Supermercado"

# Guardar con error: devuelve la fila en modo edicion con el mensaje
POST {{base_url}}/gastos/{{id_super}}/editar
[FormParams]
id_user: {{id_inaki}}
titulo: Supermercado
descripcion: Compra semanal
monto: -10
fecha: 2026-09-24
categoria: Comida
HTTP 400
[Asserts]
xpath "//li/form" exists
xpath "//li//p[@class='alert']" exists

# Guardar con exito: devuelve la fila actualizada en modo lectura
POST {{base_url}}/gastos/{{id_super}}/editar
[FormParams]
id_user: {{id_inaki}}
titulo: Supermercado
descripcion: Compra del mes
monto: 42000
fecha: 2026-09-30
categoria: Hogar
HTTP 200
[Asserts]
xpath "//li/form" not exists
xpath "string(//li//span[@class='item-sub'])" == "Compra del mes"
xpath "string(//li//span[@class='tag'])" == "Hogar"
xpath "normalize-space(//li//span[@class='amount'])" == "$ 42000.00"
xpath "normalize-space(//li//span[@class='item-meta'])" == "30/09/2026 · Usuario #{{id_inaki}}"


# Ids inexistentes o invalidos

GET {{base_url}}/gastos/99999
HTTP 404

GET {{base_url}}/gastos/99999/editar
HTTP 404

POST {{base_url}}/gastos/99999/editar
[FormParams]
id_user: {{id_inaki}}
titulo: Nada
descripcion: Nada
monto: 100
fecha: 2026-09-25
categoria: Otros
HTTP 404
[Asserts]
xpath "//li//p[@class='alert']" exists

GET {{base_url}}/gastos/abc
HTTP 400


# Borrado

# Borrar un gasto: 200 con cuerpo vacio
DELETE {{base_url}}/gastos/{{id_super}}
HTTP 200
[Asserts]
body == ""

GET {{base_url}}/gastos
HTTP 200
[Asserts]
xpath "count(//ul[@id='gasto-list']/li)" == 1

# Borrar el usuario borra tambien sus gastos (ON DELETE CASCADE)
DELETE {{base_url}}/usuarios/{{id_inaki}}
HTTP 200

GET {{base_url}}/gastos
HTTP 200
[Asserts]
xpath "count(//ul[@id='gasto-list']/li)" == 0
xpath "//p[@class='empty']" exists

# Pruebas de usuarios.
# Se asume la base vacia al empezar (make test hace server-reset antes).
# El archivo deja la base vacia al terminar, para no afectar a gastos.hurl.


# Pagina inicial

# La pagina de usuarios trae el form de creacion (dentro de #user-form) y la lista vacia
GET {{base_url}}/usuarios
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
xpath "//div[@id='user-form']/form[@hx-post='/usuarios']" exists
xpath "string(//div[@id='user-form']/form/@hx-target)" == "#user-list"
xpath "string(//div[@id='user-form']/form/@hx-target-error)" == "#user-form"
xpath "count(//ul[@id='user-list']/li)" == 0


# Creacion

# Crear: responde la lista con el usuario nuevo + el form vacio out-of-band
POST {{base_url}}/usuarios
[FormParams]
username: juan123
email: juan@gmail.com
HTTP 200
[Captures]
id_juan: xpath "string(//ul[@id='user-list']/li[1]//button[text()='Editar']/@hx-get)" regex "/usuarios/(\\d+)/editar"
[Asserts]
xpath "count(//ul[@id='user-list']/li)" == 1
xpath "string(//ul[@id='user-list']/li[1]//span[@class='item-title'])" == "juan123"
xpath "string(//ul[@id='user-list']/li[1]//span[@class='item-sub'])" == "juan@gmail.com"
xpath "//div[@id='user-form'][@hx-swap-oob='true']" exists
xpath "string(//div[@id='user-form']//input[@name='username']/@value)" == ""

# Crear un segundo usuario: la lista ahora tiene dos filas
POST {{base_url}}/usuarios
[FormParams]
username: inaki123
email: inakiuv@gmail.com
HTTP 200
[Captures]
id_inaki: xpath "string(//ul[@id='user-list']/li[2]//button[text()='Editar']/@hx-get)" regex "/usuarios/(\\d+)/editar"
[Asserts]
xpath "count(//ul[@id='user-list']/li)" == 2

# La pagina completa refleja los dos usuarios
GET {{base_url}}/usuarios
HTTP 200
[Asserts]
xpath "count(//ul[@id='user-list']/li)" == 2


# Errores de validacion en la creacion (400)

# Responde el form (sin OOB) con el error y los valores cargados.
POST {{base_url}}/usuarios
[FormParams]
username: a
email: a@gmail.com
HTTP 400
[Asserts]
xpath "//div[@id='user-form'][not(@hx-swap-oob)]" exists
xpath "//div[@id='user-form']//p[@class='alert']" exists
xpath "string(//input[@name='email']/@value)" == "a@gmail.com"
xpath "//ul[@id='user-list']" not exists

# Email sin punto en el dominio
POST {{base_url}}/usuarios
[FormParams]
username: pedro
email: pedro@gmail
HTTP 400
[Asserts]
xpath "//div[@id='user-form']//p[@class='alert']" exists

# Campos vacios
POST {{base_url}}/usuarios
[FormParams]
username:
email:
HTTP 400
[Asserts]
xpath "//div[@id='user-form']//p[@class='alert']" exists

# Los errores no crearon nada
GET {{base_url}}/usuarios
HTTP 200
[Asserts]
xpath "count(//ul[@id='user-list']/li)" == 2


# Edicion in-place

# Abrir la edicion: devuelve solo la fila con el form precargado
GET {{base_url}}/usuarios/{{id_juan}}/editar
HTTP 200
[Asserts]
xpath "count(//li)" == 1
xpath "//li/form[@hx-post='/usuarios/{{id_juan}}/editar']" exists
xpath "string(//li/form/@hx-target)" == "closest li"
xpath "string(//li/form/@hx-target-error)" == "closest li"
xpath "string(//input[@name='username']/@value)" == "juan123"
xpath "string(//input[@name='email']/@value)" == "juan@gmail.com"
xpath "//input[@id='username-{{id_juan}}']" exists

# Cancelar: devuelve la fila original en modo lectura
GET {{base_url}}/usuarios/{{id_juan}}
HTTP 200
[Asserts]
xpath "count(//li)" == 1
xpath "//li/form" not exists
xpath "string(//li//span[@class='item-title'])" == "juan123"

# Guardar con error: devuelve la fila en modo edicion con el mensaje
POST {{base_url}}/usuarios/{{id_juan}}/editar
[FormParams]
username: jj
email: juan@gmail.com
HTTP 400
[Asserts]
xpath "//li/form" exists
xpath "//li//p[@class='alert']" exists

# Guardar con exito: devuelve la fila actualizada en modo lectura
POST {{base_url}}/usuarios/{{id_juan}}/editar
[FormParams]
username: juannuevo
email: juannuevo@gmail.com
HTTP 200
[Asserts]
xpath "//li/form" not exists
xpath "string(//li//span[@class='item-title'])" == "juannuevo"
xpath "string(//li//span[@class='item-sub'])" == "juannuevo@gmail.com"

# El cambio quedo guardado
GET {{base_url}}/usuarios
HTTP 200
[Asserts]
body contains "juannuevo"
body not contains "juan123"


# Ids inexistentes o invalidos

# Usuario inexistente: 404
GET {{base_url}}/usuarios/99999
HTTP 404

GET {{base_url}}/usuarios/99999/editar
HTTP 404

# Actualizar un usuario inexistente: 404 con el form y el error
POST {{base_url}}/usuarios/99999/editar
[FormParams]
username: nadie
email: nadie@gmail.com
HTTP 404
[Asserts]
xpath "//li//p[@class='alert']" exists

# Id no numerico o no positivo: 400
GET {{base_url}}/usuarios/abc
HTTP 400

GET {{base_url}}/usuarios/0/editar
HTTP 400


# Borrado

# Borrar: 200 con cuerpo vacio (HTMX reemplaza la fila por nada)
DELETE {{base_url}}/usuarios/{{id_juan}}
HTTP 200
[Asserts]
body == ""

DELETE {{base_url}}/usuarios/{{id_inaki}}
HTTP 200
[Asserts]
body == ""

# La lista queda vacia
GET {{base_url}}/usuarios
HTTP 200
[Asserts]
xpath "count(//ul[@id='user-list']/li)" == 0

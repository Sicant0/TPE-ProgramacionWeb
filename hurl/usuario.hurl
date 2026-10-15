#insertar user1
POST {{base_url}}/usuario
Content-Type: application/json
{
    "username": "juan123",
    "email": "juan@gmail.com"
}
HTTP 201
[Asserts]
jsonpath "$.username" == "juan123"
jsonpath "$.email" == "juan@gmail.com"
[Captures]
id_usuario1: jsonpath "$['id_user']"

#insertar user2
POST {{base_url}}/usuario
Content-Type: application/json
{
    "username": "Inaki",
    "email": "inaki@gmail.com"
}
HTTP 201
[Asserts]
jsonpath "$.username" == "Inaki"
jsonpath "$.email" == "inaki@gmail.com"
[Captures]
id_usuario2: jsonpath "$['id_user']"

# Obtener user1
GET {{base_url}}/usuario/{{id_usuario1}}
HTTP 200
[Asserts]
jsonpath "$.username" == "juan123"
jsonpath "$.email" == "juan@gmail.com"

# Actualizar user2
PUT {{base_url}}/usuario/{{id_usuario2}}
Content-Type: application/json
{
    "id_user": {{id_usuario2}},
    "username": "inakinuevo",
    "email": "inakinuevo@gmail.com"
}
HTTP 200
[Asserts]
jsonpath "$.username" == "inakinuevo"
jsonpath "$.email" == "inakinuevo@gmail.com"

# Listar todos los usuarios
GET {{base_url}}/usuario
HTTP 200
[Asserts]
jsonpath "$[*]" count == 2

# Eliminar user1
DELETE {{base_url}}/usuario/{{id_usuario1}}
HTTP 204

# Eliminar user2
DELETE {{base_url}}/usuario/{{id_usuario2}}
HTTP 204

# Verifico que la lista esté vacía
GET {{base_url}}/usuario
HTTP 200
[Asserts]
jsonpath "$" count == 0
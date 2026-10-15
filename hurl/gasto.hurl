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

# Insertar gasto1
POST {{base_url}}/gasto
Content-Type: application/json
{
    "id_user": {{id_usuario1}},
    "titulo": "CompraCarrefour",
    "descripcion": "Compra semanal",
    "monto": "100000.00",
    "categoria": "Comida",
    "fecha": "2026-09-24T00:00:00Z"
}
HTTP 201
[Asserts]
jsonpath "$.id_user" == {{id_usuario1}}
jsonpath "$.titulo" == "CompraCarrefour"
jsonpath "$.descripcion" == "Compra semanal"
jsonpath "$.monto" == "100000.00"
jsonpath "$.categoria" == "Comida"
jsonpath "$.fecha" == "2026-09-24T00:00:00Z"
[Captures]
id_gasto1: jsonpath "$['id_gasto']"

# Insertar gasto2
POST {{base_url}}/gasto
Content-Type: application/json
{
    "id_user": {{id_usuario1}},
    "titulo": "CompraJean",
    "descripcion": "Compra de jean, el anterior ya estaba bastante desgastado",
    "monto": "120000.00",
    "categoria": "Ropa",
    "fecha": "2026-09-27T00:00:00Z"
}
HTTP 201
[Asserts]
jsonpath "$.id_user" == {{id_usuario1}}
jsonpath "$.titulo" == "CompraJean"
jsonpath "$.descripcion" == "Compra de jean, el anterior ya estaba bastante desgastado"
jsonpath "$.monto" == "120000.00"
jsonpath "$.categoria" == "Ropa"
jsonpath "$.fecha" == "2026-09-27T00:00:00Z"
[Captures]
id_gasto2: jsonpath "$['id_gasto']"

# Obtener gasto1
GET {{base_url}}/gasto/{{id_gasto1}}
HTTP 200
[Asserts]
jsonpath "$.id_user" == {{id_usuario1}}
jsonpath "$.titulo" == "CompraCarrefour"
jsonpath "$.descripcion" == "Compra semanal"
jsonpath "$.monto" == "100000.00"
jsonpath "$.categoria" == "Comida"
jsonpath "$.fecha" == "2026-09-24T00:00:00Z"

# Actualizar gasto2
PUT {{base_url}}/gasto/{{id_gasto2}}
Content-Type: application/json

{
    "id_user": {{id_usuario1}},
    "titulo": "CompraJean",
    "descripcion": "Compra de jean, el anterior ya estaba bastante desgastado. Ademas, necesitaba un color diferente",
    "monto": "120000.00",
    "categoria": "Ropa",
    "fecha": "2026-07-15T00:00:00Z"
}
HTTP 200
[Asserts]
jsonpath "$.id_user" == {{id_usuario1}}
jsonpath "$.titulo" == "CompraJean"
jsonpath "$.descripcion" == "Compra de jean, el anterior ya estaba bastante desgastado. Ademas, necesitaba un color diferente"
jsonpath "$.monto" == "120000.00"
jsonpath "$.categoria" == "Ropa"
jsonpath "$.fecha" == "2026-07-15T00:00:00Z"

# Listar todos los gastos
GET {{base_url}}/gasto
HTTP 200
[Asserts]
jsonpath "$[*]" count == 2

# Eliminar gasto1
DELETE {{base_url}}/gasto/{{id_gasto1}}
HTTP 204

# Eliminar gasto2
DELETE {{base_url}}/gasto/{{id_gasto2}}
HTTP 204

# Verifico que la lista de gastos este vacía
GET {{base_url}}/gasto
HTTP 200
[Asserts]
jsonpath "$" count == 0

# Eliminar usuario
DELETE {{base_url}}/usuario/{{id_usuario1}}
HTTP 204

#Verifico que no hay usuarios remanentes
GET {{base_url}}/usuario
HTTP 200
[Asserts]
jsonpath "$" count == 0
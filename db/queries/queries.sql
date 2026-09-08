-- name: GetUsuario :one
SELECT *
FROM Usuario
WHERE id_user = $1;


-- name: CreateUsuario :one
INSERT INTO Usuario (
    username,
    email
)
VALUES (
    $1,
    $2
)
RETURNING *;


-- name: UpdateUsuario :exec
UPDATE Usuario
SET
    username = $2,
    email = $3
WHERE id_user = $1;


-- name: DeleteUsuario :exec
DELETE FROM Usuario
WHERE id_user = $1;


-- name: GetGasto :one
SELECT *
FROM Gasto
WHERE id_gasto = $1;


-- name: GetGastosByUsuario :many
SELECT *
FROM Gasto
WHERE id_user = $1
ORDER BY fecha DESC;


-- name: CreateGasto :one
INSERT INTO Gasto (
    id_user,
    titulo,
    descripcion,
    monto,
    categoria,
    fecha
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


-- name: UpdateGasto :exec
UPDATE Gasto
SET
    titulo = $3,
    descripcion = $4,
    monto = $5,
    categoria = $6,
    fecha = $7
WHERE id_gasto = $1
  AND id_user = $2;


-- name: DeleteGasto :exec
DELETE FROM Gasto
WHERE id_gasto = $1
  AND id_user = $2;

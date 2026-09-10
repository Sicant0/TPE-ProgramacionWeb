package tests

import (
	"context"
	db "proyectoweb/db/sqlc"
	"testing"
	"time"
)

func TestGastoCRUD(t *testing.T) {
	ctx := context.Background()
	//CREATE USER
	usuario, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Username: "testuser",
		Email:    "testuser@example.com",
	})

	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	//CREATE GASTO
	gasto, err := queries.CreateGasto(ctx, db.CreateGastoParams{
		IDUser:      usuario.IDUser,
		Titulo:      "Test Gasto",
		Descripcion: "Este es un gasto de prueba",
		Monto:       "100.0",
		Categoria:   db.PosiblesCategoriasEntretenimiento,
		Fecha:       time.Now(),
	})

	if err != nil {
		t.Fatalf("error creating gasto: %v", err)
	}

	//UPDATE GASTO
	err = queries.UpdateGasto(ctx, db.UpdateGastoParams{
		IDGasto:     gasto.IDGasto,
		IDUser:      gasto.IDUser,
		Titulo:      "GastoUpdated",
		Descripcion: "Este es un gasto actualizado",
		Monto:       "150.0",
		Categoria:   db.PosiblesCategoriasComida,
		Fecha:       time.Now(),
	})

	if err != nil {
		t.Fatalf("error updating gasto: %v", err)
	}

	//GET GASTO
	updatedGasto, err := queries.GetGasto(ctx, gasto.IDGasto)

	if err != nil {
		t.Fatalf("error getting gasto: %v", err)
	}

	if updatedGasto.Titulo != "GastoUpdated" || updatedGasto.Descripcion != "Este es un gasto actualizado" || updatedGasto.Monto != "150.00" || updatedGasto.Categoria != db.PosiblesCategoriasComida {
		t.Fatalf("gasto not updated correctly")
	}

	//DELETE GASTO
	err = queries.DeleteGasto(ctx, db.DeleteGastoParams{
		IDGasto: gasto.IDGasto,
		IDUser:  gasto.IDUser,
	})

	if err != nil {
		t.Fatalf("error to delete gasto: %v", err)
	}

	//DELETE USER
	err = queries.DeleteUsuario(ctx, usuario.IDUser)

	if err != nil {
		t.Fatalf("error to delete user: %v", err)
	}
}

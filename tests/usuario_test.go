package tests

import (
	"context"
	db "proyectoweb/db/sqlc"
	"testing"
	"time"
)

func TestUsuarioCrud(t *testing.T) {
	ctx := context.Background()
	//CREATE USER
	usuario, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Username: "testuser",
		Email:    "testuser@example.com",
	})
	if err != nil {
		t.Fatalf("error to create user: %v", err)
	}

	//UPDATE USER
	err = queries.UpdateUsuario(ctx, db.UpdateUsuarioParams{
		IDUser:   usuario.IDUser,
		Username: "updateduser",
		Email:    "updateduser@example.com",
	})
	if err != nil {
		t.Fatalf("error to update user: %v", err)
	}

	//CREATE GASTOS AND LIST GASTOS
	gasto1, err := queries.CreateGasto(ctx, db.CreateGastoParams{
		IDUser:      usuario.IDUser,
		Titulo:      "Test Gasto",
		Descripcion: "Este es un gasto de prueba",
		Monto:       "100.0",
		Categoria:   db.PosiblesCategoriasEntretenimiento,
		Fecha:       time.Now(),
	})
	if err != nil {
		t.Fatalf("error to create gasto: %v", err)
	}

	gasto2, err := queries.CreateGasto(ctx, db.CreateGastoParams{
		IDUser:      usuario.IDUser,
		Titulo:      "Test Gasto 2",
		Descripcion: "Este es un gasto de prueba",
		Monto:       "100.0",
		Categoria:   db.PosiblesCategoriasOtros,
		Fecha:       time.Now(),
	})
	if err != nil {
		t.Fatalf("error to create gasto: %v", err)
	}

	listGastos, err := queries.GetGastosByUsuario(ctx, usuario.IDUser)
	if err != nil {
		t.Fatalf("error to list gastos: %v", err)
	}

	if len(listGastos) != 2 {
		t.Fatalf("expected 2 gastos, got %d", len(listGastos))
	}

	//DELETE GASTOS
	err = queries.DeleteGasto(ctx, db.DeleteGastoParams{
		IDGasto: gasto1.IDGasto,
		IDUser:  usuario.IDUser,
	})
	if err != nil {
		t.Fatalf("error to delete gasto: %v", err)
	}

	err = queries.DeleteGasto(ctx, db.DeleteGastoParams{
		IDGasto: gasto2.IDGasto,
		IDUser:  usuario.IDUser,
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

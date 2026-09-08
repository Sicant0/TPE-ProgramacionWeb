-- Created by Redgate Data Modeler (https://datamodeler.redgate-platform.com)
-- Last modification date: 2026-09-07 19:10:50.099

-- tipos enumerados
CREATE TYPE posibles_categorias AS ENUM (
    'Comida',
    'Transporte',
    'Hogar',
    'Salud',
    'Entretenimiento',
    'Ropa',
    'Regalos',
    'Servicios',
    'Otros'
);

-- tables
-- Table: Gasto
CREATE TABLE Gasto (
    id_gasto serial  NOT NULL,
    id_user int  NOT NULL,
    titulo varchar(15)  NOT NULL,
    descripcion text  NOT NULL,
    monto decimal(12,2)  NOT NULL,
    categoria posibles_categorias  NOT NULL,
    fecha date  NOT NULL,
    CONSTRAINT Gasto_PK PRIMARY KEY (id_gasto)
);

-- Table: Usuario
CREATE TABLE Usuario (
    id_user serial  NOT NULL,
    username varchar(25)  NOT NULL,
    email varchar(50)  NOT NULL,
    CONSTRAINT Usuario_PK PRIMARY KEY (id_user)
);

-- foreign keys
-- Reference: Gasto_Usuario_FK (table: Gasto)
ALTER TABLE Gasto ADD CONSTRAINT Gasto_Usuario_FK
    FOREIGN KEY (id_user)
    REFERENCES Usuario (id_user)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- End of file.


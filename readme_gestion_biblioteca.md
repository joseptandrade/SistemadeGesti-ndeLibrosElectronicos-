
# 📚 Sistema de Gestión de Biblioteca Digital

Este es un sistema completo de gestión de libros electrónicos desarrollado en Go, con frontend en HTML/JavaScript y backend en MySQL. Permite administrar libros, usuarios, préstamos y reportes de forma sencilla y visual.

---

## Características

- CRUD de Libros
- CRUD de Usuarios
- Registro y devolución de préstamos
- Filtros por estado de préstamo
- Reportes:
  - Libros disponibles/prestados
  - Usuarios con préstamos activos
- Validaciones de formularios
- Interfaz web

---

## Tecnologías Utilizadas

- **Lenguaje Backend**: Go (Golang)
- **Base de Datos**: MySQL
- **Frontend**: HTML + JavaScript + CSS puro
- **Router**: Gorilla Mux

---

## Estructura del Proyecto

SistemaLibros/
├── db/              # Configuración y conexión a MySQL
├── libros/          # Lógica de gestión de libros
├── usuarios/        # Lógica de gestión de usuarios
├── prestamos/       # Lógica de préstamos
├── reportes/        # Módulo de reportes
├── Templates/       # Interfaz web (index.html, app.js)
├── main.go          # Punto de entrada de la app

---

## Requisitos

- Go 1.20 o superior
- MySQL 5.7+ o compatible
- Navegador moderno (Chrome, Firefox, Edge)

---

## Instalación y Uso

### 1. Clonar el repositorio

git clone https://github.com/tu-usuario/SistemaLibros.git
cd SistemaLibros

### 2. Crear la base de datos y las tablas

Ejecuta el siguiente script en tu servidor MySQL:

-- Crear base de datos
CREATE DATABASE IF NOT EXISTS biblioteca;
USE biblioteca;

-- Tabla de Usuarios
CREATE TABLE Usuarios (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Nombre VARCHAR(100) NOT NULL
);

-- Tabla de Libros
CREATE TABLE Libros (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Titulo VARCHAR(255) NOT NULL,
    Autor VARCHAR(100) NOT NULL,
    Estado ENUM('disponible', 'prestado') DEFAULT 'disponible'
);

-- Tabla de Préstamos
CREATE TABLE Prestamos (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    LibroID INT NOT NULL,
    UsuarioID INT NOT NULL,
    FechaPrestamo DATE NOT NULL,
    Estado ENUM('vigente', 'devuelto') DEFAULT 'vigente',
    FOREIGN KEY (LibroID) REFERENCES Libros(ID),
    FOREIGN KEY (UsuarioID) REFERENCES Usuarios(ID)
);

### 3. Configurar conexión en db/db.go

Credenciales:

usuario := "root"
contraseña := ""
servidor := "localhost"
puerto := "3306"
baseDatos := "biblioteca"

### 4. Ejecutar el backend

go run main.go

### 5. Abrir en el navegador

http://localhost:8080

---

## Autor

Desarrollado por José Leandro Andrade Zalamea  

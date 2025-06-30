//PROYECTO: Sistema de Gestión de Biblioteca Digital
// FECHA: 28 de junio de 2025
// DESARROLLADOR: Jose Andrade Zalamea

package main

import (
	"fmt"
	"log"
	"net/http"

	"SistemaLibros/db"
	"SistemaLibros/libros"
	"SistemaLibros/prestamos"
	"SistemaLibros/reportes"
	"SistemaLibros/usuarios"

	"github.com/gorilla/mux"
)

func main() {
	// ✅ Conexión a la base de datos
	db.Conectar()
	defer db.DB.Close()

	// Enrutador con gorilla/mux
	r := mux.NewRouter()

	// ============================
	// RUTAS - GESTIÓN DE LIBROS
	// ============================
	r.HandleFunc("/libros", libros.GetLibros).Methods("GET")
	r.HandleFunc("/libros/{id}", libros.GetLibro).Methods("GET")
	r.HandleFunc("/libros", libros.AddLibro).Methods("POST")
	r.HandleFunc("/libros/{id}", libros.UpdateLibro).Methods("PUT")
	r.HandleFunc("/libros/{id}", libros.DeleteLibro).Methods("DELETE")
	r.HandleFunc("/libros/filter", libros.FilterLibros).Methods("GET")

	// ============================
	// RUTAS - GESTIÓN DE USUARIOS
	// ============================
	r.HandleFunc("/usuarios", usuarios.GetUsuarios).Methods("GET")
	r.HandleFunc("/usuarios", usuarios.AddUsuario).Methods("POST")
	r.HandleFunc("/usuarios/{id}", usuarios.UpdateUsuario).Methods("PUT")
	r.HandleFunc("/usuarios/{id}", usuarios.DeleteUsuario).Methods("DELETE")

	// ============================
	// RUTAS - GESTIÓN DE PRÉSTAMOS
	// ============================
	r.HandleFunc("/prestamos", prestamos.GetPrestamos).Methods("GET")
	r.HandleFunc("/prestamos", prestamos.AddPrestamo).Methods("POST")
	r.HandleFunc("/prestamos/return/{id}", prestamos.ReturnPrestamo).Methods("PUT")
	r.HandleFunc("/prestamos/filter", prestamos.FilterPrestamos).Methods("GET")

	// ============================
	// RUTAS - REPORTES
	// ============================
	r.HandleFunc("/reportes/libros", reportes.ReporteLibros).Methods("GET")
	r.HandleFunc("/reportes/usuarios/prestamos-activos", reportes.ReporteUsuariosConPrestamosActivos).Methods("GET")

	// ============================
	// HTML Estático (Frontend)
	// ============================
	fs := http.FileServer(http.Dir("./Templates"))
	r.PathPrefix("/").Handler(fs)

	// ============================
	// Iniciar servidor
	// ============================
	fmt.Println("🚀 Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

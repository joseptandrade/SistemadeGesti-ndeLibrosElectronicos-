//PROYECTO: Sistema de Gestión de Biblioteca Digital
// FECHA: 28 de junio de 2025
// DESARROLLADOR: Jose Andrade Zalamea

package prestamos

import (
	"SistemaLibros/db" // Conexión a la base de datos
	"database/sql"     // Para manejar sql.Rows
	"encoding/json"    // Para codificar/decodificar JSON
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Prestamo representa un préstamo registrado en la base de datos
type Prestamo struct {
	ID        int       `json:"id"`
	LibroID   int       `json:"libro_id"`
	UsuarioID int       `json:"usuario_id"`
	Fecha     time.Time `json:"fecha"`  // Fecha del préstamo
	Estado    string    `json:"estado"` // "vigente" o "devuelto"
}

// GetPrestamos obtiene todos los préstamos desde MySQL
func GetPrestamos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT ID, LibroID, UsuarioID, FechaPrestamo, Estado FROM Prestamos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var prestamos []Prestamo
	for rows.Next() {
		var p Prestamo
		if err := rows.Scan(&p.ID, &p.LibroID, &p.UsuarioID, &p.Fecha, &p.Estado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		prestamos = append(prestamos, p)
	}

	if len(prestamos) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prestamos)
}

// AddPrestamo agrega un nuevo préstamo a la base de datos
func AddPrestamo(w http.ResponseWriter, r *http.Request) {
	var p Prestamo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Valores por defecto
	p.Fecha = time.Now()
	p.Estado = "vigente"

	result, err := db.DB.Exec(
		"INSERT INTO Prestamos (LibroID, UsuarioID, FechaPrestamo, Estado) VALUES (?, ?, ?, ?)",
		p.LibroID, p.UsuarioID, p.Fecha, p.Estado,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// ReturnPrestamo cambia el estado de un préstamo a "devuelto"
func ReturnPrestamo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("UPDATE Prestamos SET Estado = 'devuelto' WHERE ID = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Préstamo no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FilterPrestamos permite filtrar préstamos por estado (o listar todos si se omite)
func FilterPrestamos(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")

	var rows *sql.Rows
	var err error

	// Si no se especifica estado, se consultan todos los préstamos
	if estado == "" {
		rows, err = db.DB.Query("SELECT ID, LibroID, UsuarioID, FechaPrestamo, Estado FROM Prestamos")
	} else {
		rows, err = db.DB.Query("SELECT ID, LibroID, UsuarioID, FechaPrestamo, Estado FROM Prestamos WHERE Estado = ?", estado)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var prestamos []Prestamo
	for rows.Next() {
		var p Prestamo
		if err := rows.Scan(&p.ID, &p.LibroID, &p.UsuarioID, &p.Fecha, &p.Estado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		prestamos = append(prestamos, p)
	}

	if len(prestamos) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prestamos)
}

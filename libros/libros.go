//PROYECTO: Sistema de Gestión de Biblioteca Digital
// FECHA: 28 de junio de 2025
// DESARROLLADOR: Jose Andrade Zalamea

package libros

import (
	"SistemaLibros/db"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Libro representa la estructura de un libro en la base de datos
type Libro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Estado string `json:"estado"` // disponible o prestado
}

// GetLibros obtiene todos los libros desde MySQL
func GetLibros(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT ID, Titulo, Autor, Estado FROM Libros")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		if err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Estado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		libros = append(libros, l)
	}

	if len(libros) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libros)
}

// GetLibro obtiene un solo libro por ID
func GetLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var l Libro
	err = db.DB.QueryRow("SELECT ID, Titulo, Autor, Estado FROM Libros WHERE ID = ?", id).
		Scan(&l.ID, &l.Titulo, &l.Autor, &l.Estado)
	if err == sql.ErrNoRows {
		http.Error(w, "Libro no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(l)
}

// AddLibro agrega un nuevo libro a la base de datos
func AddLibro(w http.ResponseWriter, r *http.Request) {
	var l Libro
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if l.Estado == "" {
		l.Estado = "disponible"
	}

	result, err := db.DB.Exec("INSERT INTO Libros (Titulo, Autor, Estado) VALUES (?, ?, ?)", l.Titulo, l.Autor, l.Estado)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	l.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(l)
}

// UpdateLibro actualiza un libro existente por ID
func UpdateLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var l Libro
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE Libros SET Titulo = ?, Autor = ?, Estado = ? WHERE ID = ?", l.Titulo, l.Autor, l.Estado, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l.ID = id
	json.NewEncoder(w).Encode(l)
}

// DeleteLibro elimina un libro solo si no tiene préstamos registrados
func DeleteLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar si hay préstamos asociados a este libro
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM Prestamos WHERE LibroID = ? AND Estado = 'vigente'", id).Scan(&count)
	if err != nil {
		http.Error(w, "Error al verificar préstamos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "No se puede eliminar el libro porque tiene préstamos vigentes", http.StatusForbidden)
		return
	}

	// Eliminar el libro si no tiene préstamos
	result, err := db.DB.Exec("DELETE FROM Libros WHERE ID = ?", id)
	if err != nil {
		http.Error(w, "Error al eliminar libro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Libro no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FilterLibros permite filtrar libros por estado o autor
func FilterLibros(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")
	autor := r.URL.Query().Get("autor")

	query := "SELECT ID, Titulo, Autor, Estado FROM Libros WHERE 1=1"
	var params []interface{}
	if estado != "" {
		query += " AND Estado = ?"
		params = append(params, estado)
	}
	if autor != "" {
		query += " AND Autor LIKE ?"
		params = append(params, "%"+autor+"%")
	}

	rows, err := db.DB.Query(query, params...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		if err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.Estado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		libros = append(libros, l)
	}

	if len(libros) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libros)
}

package usuarios

import (
	"SistemaLibros/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Usuario representa un usuario del sistema
type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// GetUsuarios obtiene todos los usuarios
func GetUsuarios(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT ID, Nombre FROM Usuarios")
	if err != nil {
		http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var u Usuario
		if err := rows.Scan(&u.ID, &u.Nombre); err != nil {
			http.Error(w, "Error al leer usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}
		usuarios = append(usuarios, u)
	}

	if len(usuarios) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

// AddUsuario agrega un nuevo usuario
func AddUsuario(w http.ResponseWriter, r *http.Request) {
	var u Usuario
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error en el formato del JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO Usuarios (Nombre) VALUES (?)", u.Nombre)
	if err != nil {
		http.Error(w, "Error al insertar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	u.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// UpdateUsuario actualiza un usuario existente
func UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var u Usuario
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error al leer datos del usuario: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE Usuarios SET Nombre = ? WHERE ID = ?", u.Nombre, id)
	if err != nil {
		http.Error(w, "Error al actualizar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	u.ID = id
	json.NewEncoder(w).Encode(u)
}

// DeleteUsuario elimina un usuario solo si no tiene préstamos
func DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar si el usuario tiene préstamos registrados
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM Prestamos WHERE UsuarioID = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Error al verificar préstamos del usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "No se puede eliminar el usuario porque tiene préstamos registrados", http.StatusForbidden)
		return
	}

	result, err := db.DB.Exec("DELETE FROM Usuarios WHERE ID = ?", id)
	if err != nil {
		http.Error(w, "Error al eliminar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

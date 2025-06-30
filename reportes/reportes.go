//PROYECTO: Sistema de Gestión de Biblioteca Digital
// FECHA: 28 de junio de 2025
// DESARROLLADOR: Jose Andrade Zalamea

package reportes

import (
	"SistemaLibros/db" // Importa la conexión global a MySQL
	"encoding/json"    // Para codificar respuestas en JSON
	"net/http"         // Para manejar peticiones HTTP
)

// ReporteLibros consulta la base de datos y cuenta los libros disponibles y prestados
func ReporteLibros(w http.ResponseWriter, r *http.Request) {
	// Consulta solo el estado de los libros desde la base de datos
	rows, err := db.DB.Query("SELECT Estado FROM Libros")
	if err != nil {
		http.Error(w, "Error al consultar libros: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	disponibles := 0
	prestados := 0

	// Recorre los resultados y cuenta según estado
	for rows.Next() {
		var estado string
		if err := rows.Scan(&estado); err != nil {
			http.Error(w, "Error al leer datos: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if estado == "disponible" {
			disponibles++
		} else {
			prestados++
		}
	}

	// Construye el mapa con los resultados y responde en JSON
	report := map[string]int{
		"disponibles": disponibles,
		"prestados":   prestados,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// ReporteUsuariosConPrestamosActivos devuelve los usuarios con préstamos vigentes
func ReporteUsuariosConPrestamosActivos(w http.ResponseWriter, r *http.Request) {
	// Consulta SQL que une usuarios y préstamos donde el estado sea 'vigente'
	query := `
		SELECT DISTINCT u.ID, u.Nombre
		FROM Usuarios u
		INNER JOIN Prestamos p ON u.ID = p.UsuarioID
		WHERE p.Estado = 'vigente'
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Error al consultar usuarios con préstamos activos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Define una estructura local para representar al usuario
	type Usuario struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
	}

	var usuarios []Usuario

	// Recorre los resultados de la consulta
	for rows.Next() {
		var u Usuario
		if err := rows.Scan(&u.ID, &u.Nombre); err != nil {
			http.Error(w, "Error al leer usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		usuarios = append(usuarios, u)
	}

	// Devuelve la lista de usuarios en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

//PROYECTO: Sistema de Gestión de Biblioteca Digital
// FECHA: 28 de junio de 2025
// DESARROLLADOR: Jose Andrade Zalamea

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
)

// DB es la conexión global accesible desde todo el sistema
var DB *sql.DB

// Conectar inicializa la conexión a MySQL
func Conectar() {
	var err error

	// Configuración de la conexión
	username := "root"     // Usuario MySQL
	password := ""         // Contraseña esta vacía
	host := "localhost"    // Host local
	port := 3306           // Puerto por defecto de MySQL
	dbname := "Biblioteca" // Mi base de datos

	// Cadena de conexión DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		username, password, host, port, dbname)

	// Abre conexión
	tempDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error abriendo conexión con MySQL:", err)
	}

	// Verifica que la conexión
	err = tempDB.Ping()
	if err != nil {
		log.Fatal("❌ No se pudo conectar a MySQL:", err)
	}

	DB = tempDB
	log.Println("✅ Conectado a MySQL exitosamente")
}

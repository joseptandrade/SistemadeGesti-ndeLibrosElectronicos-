package main

import (
	"bufio"                // Para lectura de entradas desde la consola
	"fmt"                  // Para mostrar mensajes en consola
	"os"                   // Para acceder al sistema operativo (entrada estándar)
	"sistemalibros/libros" // Importa el módulo donde está la lógica del sistema (encapsulada)
	"strings"              // Para procesar texto (quitar espacios, comparar)
)

func main() {
	// Crea un lector de entrada de texto desde consola
	reader := bufio.NewReader(os.Stdin)

	// Se inicializa la "clase" GestorLibros, responsable de toda la lógica.
	// ENCAPSULACIÓN: toda la lógica de libros está contenida en este objeto.
	gestor := libros.NewGestorLibros()

	// Menú interactivo en consola
	for {
		// Mostrar las opciones disponibles
		fmt.Println("\n1. Registrar Libro")
		fmt.Println("2. Buscar Libro")
		fmt.Println("3. Ver Libros Disponibles")
		fmt.Println("4. Editar Libro")
		fmt.Println("5. Eliminar Libro")
		fmt.Println("6. Filtrar por Estado")
		fmt.Println("7. Salir")
		fmt.Print("Seleccione una opción: ")

		// Leer opción seleccionada
		opcion, _ := reader.ReadString('\n')
		opcion = strings.TrimSpace(opcion)

		switch opcion {
		case "1":
			// REGISTRAR LIBRO
			fmt.Print("Título: ")
			titulo, _ := reader.ReadString('\n')
			fmt.Print("Autor: ")
			autor, _ := reader.ReadString('\n')
			fmt.Print("ISBN: ")
			isbn, _ := reader.ReadString('\n')

			// Crear una instancia del struct Libro
			libro := libros.Libro{
				Titulo: strings.TrimSpace(titulo),
				Autor:  strings.TrimSpace(autor),
				ISBN:   strings.TrimSpace(isbn),
			}

			// Llamada al método del gestor para registrar
			// MANEJO DE ERRORES: si el libro ya existe, se informa
			err := gestor.RegistrarLibro(libro)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Libro registrado.")
			}

		case "2":
			// BUSCAR LIBRO
			fmt.Print("Ingrese título o autor: ")
			filtro, _ := reader.ReadString('\n')

			// Llamada al método de búsqueda (insensible a mayúsculas)
			encontrados := gestor.BuscarLibro(strings.TrimSpace(filtro))

			// Mostrar resultados
			for _, l := range encontrados {
				fmt.Printf(" %s - %s [%s]\n", l.Titulo, l.Autor, l.Estado)
			}

		case "3":
			// VER LIBROS DISPONIBLES
			disponibles := gestor.FiltrarPorEstado(libros.Disponible)
			for _, l := range disponibles {
				fmt.Printf(" %s - %s\n", l.Titulo, l.Autor)
			}

		case "4":
			// EDITAR LIBRO
			fmt.Print("Ingrese ISBN del libro a editar: ")
			isbn, _ := reader.ReadString('\n')

			fmt.Print("Nuevo Título: ")
			nuevoTitulo, _ := reader.ReadString('\n')
			fmt.Print("Nuevo Autor: ")
			nuevoAutor, _ := reader.ReadString('\n')

			// Llama al método que edita los datos de un libro
			// MANEJO DE ERRORES: se informa si el ISBN no existe
			err := gestor.EditarLibro(strings.TrimSpace(isbn), strings.TrimSpace(nuevoTitulo), strings.TrimSpace(nuevoAutor))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Libro actualizado correctamente.")
			}

		case "5":
			// ELIMINAR LIBRO
			fmt.Print("Ingrese ISBN del libro a eliminar: ")
			isbn, _ := reader.ReadString('\n')

			// Llama al método para eliminar el libro
			err := gestor.EliminarLibro(strings.TrimSpace(isbn))
			if err != nil {
				fmt.Println("Error:", err) // MANEJO DE ERRORES: ISBN no encontrado
			} else {
				fmt.Println("Libro eliminado correctamente.")
			}

		case "6":
			// FILTRAR POR ESTADO
			fmt.Print("Ingrese estado (Disponible o Prestado): ")
			estadoStr, _ := reader.ReadString('\n')
			estadoStr = strings.TrimSpace(estadoStr)

			var estado libros.Estado
			if estadoStr == "Disponible" {
				estado = libros.Disponible
			} else if estadoStr == "Prestado" {
				estado = libros.Prestado
			} else {
				fmt.Println("Estado inválido.")
				continue
			}

			// Mostrar resultados filtrados por estado
			resultado := gestor.FiltrarPorEstado(estado)
			for _, l := range resultado {
				fmt.Printf(" %s - %s [%s]\n", l.Titulo, l.Autor, l.Estado)
			}

		case "7":
			// SALIR DEL SISTEMA
			fmt.Println("Saliendo...")
			return

		default:
			// OPCIÓN NO VÁLIDA
			fmt.Println("Opción inválida")
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sistemalibros/libros" // Importa el módulo que gestiona los libros
	"strings"
)

func main() {
	// Crea un lector para leer entradas del usuario desde la consola
	reader := bufio.NewReader(os.Stdin)

	// Se crea una instancia del gestor de libros (usa encapsulación)
	gestor := libros.NewGestorLibros()

	// Bucle principal del menú de opciones
	for {
		// Mostrar menú al usuario
		fmt.Println("\n1. Registrar Libro")
		fmt.Println("2. Buscar Libro")
		fmt.Println("3. Ver Libros Disponibles")
		fmt.Println("4. Editar Libro")
		fmt.Println("5. Eliminar Libro")
		fmt.Println("6. Filtrar por Estado")
		fmt.Println("7. Salir")
		fmt.Print("Seleccione una opción: ")

		// Leer la opción ingresada
		opcion, _ := reader.ReadString('\n')
		opcion = strings.TrimSpace(opcion)

		switch opcion {
		case "1":
			// Registrar nuevo libro
			fmt.Print("Título: ")
			titulo, _ := reader.ReadString('\n')
			fmt.Print("Autor: ")
			autor, _ := reader.ReadString('\n')
			fmt.Print("ISBN: ")
			isbn, _ := reader.ReadString('\n')

			libro := libros.Libro{
				Titulo: strings.TrimSpace(titulo),
				Autor:  strings.TrimSpace(autor),
				ISBN:   strings.TrimSpace(isbn),
			}

			err := gestor.RegistrarLibro(libro)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Libro registrado.")
			}

		case "2":
			// Buscar libro por título o autor
			fmt.Print("Ingrese título o autor: ")
			filtro, _ := reader.ReadString('\n')
			encontrados := gestor.BuscarLibro(strings.TrimSpace(filtro))
			for _, l := range encontrados {
				fmt.Printf(" %s - %s [%s]\n", l.Titulo, l.Autor, l.Estado)
			}

		case "3":
			// Ver libros disponibles
			disponibles := gestor.FiltrarPorEstado(libros.Disponible)
			for _, l := range disponibles {
				fmt.Printf(" %s - %s\n", l.Titulo, l.Autor)
			}

		case "4":
			// Editar libro existente
			fmt.Print("Ingrese ISBN del libro a editar: ")
			isbn, _ := reader.ReadString('\n')
			isbn = strings.TrimSpace(isbn)

			fmt.Print("Nuevo Título: ")
			nuevoTitulo, _ := reader.ReadString('\n')
			fmt.Print("Nuevo Autor: ")
			nuevoAutor, _ := reader.ReadString('\n')

			err := gestor.EditarLibro(isbn, strings.TrimSpace(nuevoTitulo), strings.TrimSpace(nuevoAutor))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Libro actualizado correctamente.")
			}

		case "5":
			// Eliminar libro por ISBN
			fmt.Print("Ingrese ISBN del libro a eliminar: ")
			isbn, _ := reader.ReadString('\n')
			err := gestor.EliminarLibro(strings.TrimSpace(isbn))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Libro eliminado correctamente.")
			}

		case "6":
			// Filtrar libros por estado
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

			resultado := gestor.FiltrarPorEstado(estado)
			for _, l := range resultado {
				fmt.Printf(" %s - %s [%s]\n", l.Titulo, l.Autor, l.Estado)
			}

		case "7":
			// Salir del programa
			fmt.Println("Saliendo...")
			return

		default:
			fmt.Println("Opción inválida")
		}
	}
}

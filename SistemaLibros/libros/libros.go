package libros

import (
	"errors"
	"strings"
)

// Gestor es una interfaz que define los métodos para gestionar libros.
// Esto permite desacoplar la implementación de la lógica de negocio.
type Gestor interface {
	RegistrarLibro(libro Libro) error
	BuscarLibro(filtro string) []Libro
	FiltrarPorEstado(estado Estado) []Libro
}

// Estado representa el estado de un libro (Disponible o Prestado)
type Estado string

const (
	Disponible Estado = "Disponible"
	Prestado   Estado = "Prestado"
)

// Clase Libro representa la estructura de un libro en el sistema
type Libro struct {
	Titulo string
	Autor  string
	ISBN   string
	Estado Estado
}

// Clase GestorLibros es la estructura que almacena y gestiona una colección de libros
type GestorLibros struct {
	libros []Libro
}

// NewGestorLibros crea e inicializa un nuevo GestorLibros vacío
func NewGestorLibros() *GestorLibros {
	return &GestorLibros{}
}

// RegistrarLibro agrega un libro a la colección si el ISBN no está duplicado.
// Si el libro ya existe, retorna un error.
func (g *GestorLibros) RegistrarLibro(libro Libro) error {
	for _, l := range g.libros {
		if l.ISBN == libro.ISBN {
			return errors.New("libro ya registrado con ese ISBN")
		}
	}
	libro.Estado = Disponible // Se registra como disponible por defecto
	g.libros = append(g.libros, libro)
	return nil
}

// BuscarLibro busca libros por título o autor, ignorando mayúsculas/minúsculas.
// Devuelve una lista de libros que coincidan con el filtro.
func (g *GestorLibros) BuscarLibro(filtro string) []Libro {
	var encontrados []Libro
	for _, l := range g.libros {
		if strings.Contains(strings.ToLower(l.Titulo), strings.ToLower(filtro)) ||
			strings.Contains(strings.ToLower(l.Autor), strings.ToLower(filtro)) {
			encontrados = append(encontrados, l)
		}
	}
	return encontrados
}

// FiltrarPorEstado devuelve los libros que tienen el estado solicitado
// (por ejemplo: solo los disponibles o solo los prestados)
func (g *GestorLibros) FiltrarPorEstado(estado Estado) []Libro {
	var filtrados []Libro
	for _, l := range g.libros {
		if l.Estado == estado {
			filtrados = append(filtrados, l)
		}
	}
	return filtrados
}

// EditarLibro permite modificar los datos de un libro existente por su ISBN.
// Si no se encuentra el libro, devuelve un error.
func (g *GestorLibros) EditarLibro(isbn string, nuevoTitulo, nuevoAutor string) error {
	for i, l := range g.libros {
		if l.ISBN == isbn {
			g.libros[i].Titulo = nuevoTitulo
			g.libros[i].Autor = nuevoAutor
			return nil
		}
	}
	return errors.New("libro no encontrado")
}

// EliminarLibro elimina un libro de la colección según su ISBN.
// Si no se encuentra el libro, devuelve un error.
func (g *GestorLibros) EliminarLibro(isbn string) error {
	for i, l := range g.libros {
		if l.ISBN == isbn {
			// Remueve el libro de la lista usando slicing
			g.libros = append(g.libros[:i], g.libros[i+1:]...)
			return nil
		}
	}
	return errors.New("libro no encontrado")
}

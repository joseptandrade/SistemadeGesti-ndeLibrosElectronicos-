package libros

import (
	"errors"  // Para manejo de errores personalizados
	"strings" // Para manipulación de texto (búsquedas, limpieza, etc.)
)

// INTERFAZ: Gestor
// Define los métodos que cualquier tipo gestor debe implementar.
// Permite desacoplar implementación del uso.
type Gestor interface {
	RegistrarLibro(libro Libro) error
	BuscarLibro(filtro string) []Libro
	FiltrarPorEstado(estado Estado) []Libro
	EditarLibro(isbn string, nuevoTitulo, nuevoAutor string) error
	EliminarLibro(isbn string) error
}

// ENUMERACIÓN: Estado
// Representa el estado lógico del libro.
type Estado string

// Constantes que definen los posibles estados de un libro
const (
	Disponible Estado = "Disponible"
	Prestado   Estado = "Prestado"
)

// CLASE (struct): Libro
// Representa un objeto tipo Libro con sus atributos principales.
type Libro struct {
	Titulo string
	Autor  string
	ISBN   string
	Estado Estado
}

// CLASE (struct): GestorLibros
// Aplica ENCAPSULACIÓN mediante el campo privado 'libros'.
type GestorLibros struct {
	libros []Libro // Campo privado que almacena los libros registrados
}

// Constructor de la clase GestorLibros.
// Inicializa el gestor con una colección vacía de libros.
func NewGestorLibros() *GestorLibros {
	return &GestorLibros{}
}

// Método que registra un libro nuevo si el ISBN no está repetido.
// MANEJO DE ERROR: control de duplicados
func (g *GestorLibros) RegistrarLibro(libro Libro) error {
	for _, l := range g.libros {
		if l.ISBN == libro.ISBN {
			return errors.New("libro ya registrado con ese ISBN") // MANEJO DE ERROR: control de duplicados
		}
	}
	libro.Estado = Disponible // Todo libro nuevo se registra como disponible
	g.libros = append(g.libros, libro)
	return nil
}

// Método para buscar libros por título o autor (no distingue mayúsculas).
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

// Método para obtener los libros que tengan un estado específico (Disponible o Prestado).
func (g *GestorLibros) FiltrarPorEstado(estado Estado) []Libro {
	var filtrados []Libro
	for _, l := range g.libros {
		if l.Estado == estado {
			filtrados = append(filtrados, l)
		}
	}
	return filtrados
}

// Método para editar los datos de un libro identificado por su ISBN.
func (g *GestorLibros) EditarLibro(isbn string, nuevoTitulo, nuevoAutor string) error {
	for i, l := range g.libros {
		if l.ISBN == isbn {
			g.libros[i].Titulo = nuevoTitulo
			g.libros[i].Autor = nuevoAutor
			return nil
		}
	}
	return errors.New("libro no encontrado") // MANEJO DE ERROR: control de inexistencia
}

// Método para eliminar un libro por su ISBN.
// Remueve el libro de la colección si lo encuentra.
// MANEJO DE ERROR: control de inexistencia
func (g *GestorLibros) EliminarLibro(isbn string) error {
	for i, l := range g.libros {
		if l.ISBN == isbn {
			// Remueve el libro usando slicing
			g.libros = append(g.libros[:i], g.libros[i+1:]...)
			return nil
		}
	}
	return errors.New("libro no encontrado") // MANEJO DE ERROR: control de inexistencia
}

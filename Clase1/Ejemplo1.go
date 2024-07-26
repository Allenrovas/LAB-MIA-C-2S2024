package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Profesor struct {
	Tipo        int32    //4 bytes
	Id_profesor int32    //4 bytes
	CUI         [13]byte //13 bytes
	Nombre      [25]byte //25 bytes
	Cursos      [25]byte //25 bytes
}

type Estudiante struct {
	Tipo          int32    //4 bytes
	Id_estudiante int32    //4 bytes
	CUI           [13]byte //13 bytes
	Nombre        [25]byte //25 bytes
	Carnet        [25]byte //25 bytes
}

func main() {
	crearArchivo()
	Menu()
}

//Ejecutar el menu

func Menu() {
	var ValorMenu string

	fmt.Println("")
	fmt.Println("Sistema de Registro de estudiantes y profesores")
	fmt.Println("1. Registro de Profesores")
	fmt.Println("2. Registro de Estudiantes")
	fmt.Println("3. Ver registros")
	fmt.Println("4. Salir")
	fmt.Println("")
	fmt.Println("Por favor seleccione una opcion: ")
	fmt.Scan(&ValorMenu)

	if ValorMenu == "1" {
		RegistroProfesor()
	} else if ValorMenu == "2" {
		RegistroEstudiante()
	} else if ValorMenu == "3" {
		VerRegistros()
	} else if ValorMenu == "4" {
		os.Exit(0)
	} else {
		fmt.Println("Opcion no valida")
	}

	Menu()
}

// Funcion para registrar profesor
func RegistroProfesor() {
	var id int32
	var cui, nombre, curso string

	//Abrir archivo en modo escritura
	arch, err := os.OpenFile("Registros.bin", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	arch.Seek(0, io.SeekEnd)

	var profesorNuevo Profesor
	profesorNuevo.Tipo = int32(1)

	fmt.Println("Ingrese el ID del profesor: ")
	fmt.Scan(&id)
	profesorNuevo.Id_profesor = id

	fmt.Println("Ingrese el CUI del profesor: ")
	fmt.Scan(&cui)
	copy(profesorNuevo.CUI[:], cui)

	fmt.Println("Ingrese el nombre del profesor: ")
	fmt.Scan(&nombre)
	copy(profesorNuevo.Nombre[:], nombre)

	fmt.Println("Ingrese el curso del profesor: ")
	fmt.Scan(&curso)
	copy(profesorNuevo.Cursos[:], []byte(curso))

	//Escribir la estructura completa en el archivo binario
	binary.Write(arch, binary.LittleEndian, &profesorNuevo)
	arch.Close()
	fmt.Println("Se registro un nuevo profesor")
}

// Funcion para registrar estudiante
func RegistroEstudiante() {
	var id int32
	var cui, nombre, carnet string

	//Abrir archivo en modo escritura
	arch, err := os.OpenFile("Registros.bin", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	arch.Seek(0, io.SeekEnd)

	var estudianteNuevo Estudiante
	estudianteNuevo.Tipo = int32(2)

	fmt.Println("Ingrese el ID del estudiante: ")
	fmt.Scan(&id)
	estudianteNuevo.Id_estudiante = id

	fmt.Println("Ingrese el CUI del estudiante: ")
	fmt.Scan(&cui)
	copy(estudianteNuevo.CUI[:], cui)

	fmt.Println("Ingrese el nombre del estudiante: ")
	fmt.Scan(&nombre)
	copy(estudianteNuevo.Nombre[:], nombre)

	fmt.Println("Ingrese el carnet del estudiante: ")
	fmt.Scan(&carnet)
	copy(estudianteNuevo.Carnet[:], []byte(carnet))

	//Escribir la estructura completa en el archivo binario
	binary.Write(arch, binary.LittleEndian, &estudianteNuevo)
	arch.Close()
	fmt.Println("Se registro un nuevo estudiante")
}

// Ver los registros
func VerRegistros() {

	fmt.Println("")
	fmt.Println("Registros")

	//Abrir el archivo en modo lectura
	arch, err := os.OpenFile("Registros.bin", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	//Leer el archivo con un bucle
	for {
		//Leer como profesor y el tipo
		var profesor Profesor
		err = binary.Read(arch, binary.LittleEndian, &profesor)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		//Si el tipo 1, es profesor
		if profesor.Tipo == 1 {
			fmt.Println("Profesor")
			fmt.Println("ID: ", profesor.Id_profesor)
			fmt.Println("CUI: ", string(profesor.CUI[:]))
			fmt.Println("Nombre: ", string(profesor.Nombre[:]))
			fmt.Println("Curso: ", string(profesor.Cursos[:]))
			fmt.Println("")
		} else if profesor.Tipo == 2 {
			fmt.Println("Estudiante")
			fmt.Println("ID: ", profesor.Id_profesor)
			fmt.Println("CUI: ", string(profesor.CUI[:]))
			fmt.Println("Nombre: ", string(profesor.Nombre[:]))
			fmt.Println("Carnet: ", string(profesor.Cursos[:]))
			fmt.Println("")
		}
	}

}

// Crea un archivo binario
func crearArchivo() {
	if _, err := os.Stat("Registros.bin"); os.IsNotExist(err) {
		arch, err := os.Create("Registros.bin")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer arch.Close()
	}
}

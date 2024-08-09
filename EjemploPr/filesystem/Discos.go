package filesystem

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// Funcion que crea los discos binarios
// respuesta += CrearDisco(sizeInt, valorUnit, valorFit, valorPath)
func CrearDiscos(size int, unit string, fit string, pathValor string) string {
	var respuesta string
	//Eliminar el nombre del disco, path, por ejemplo> /home/user/disco1.mia
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Path: " + dirPath)
	fmt.Println("Nombre: " + fileName)
	//Tamano en bytes del disco
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	} else {
		fmt.Println("Error: Unit no reconocido")
		respuesta += "Error: Unit no reconocido\n"
		return respuesta
	}
	//Crear el directorio
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		fmt.Println("Error: No se pudo crear el directorio")
		respuesta += "Error: No se pudo crear el directorio\n"
		return respuesta
	}
	//Crear el archivo
	archivo, err := os.Create(pathValor)
	if err != nil {
		fmt.Println("Error: No se pudo crear el archivo")
		respuesta += "Error: No se pudo crear el archivo\n"
		return respuesta
	}
	defer archivo.Close()
	//Escribir en el archivo
	randomNum := rand.Intn(99) + 1
	var disk MBR

	disk.Mbr_tamano = int32(size)
	disk.Mbr_disk_signature = int32(randomNum)
	fitAux := []byte(fit)
	disk.Dsk_fit = [1]byte{fitAux[0]}
	fechaA := time.Now()
	fecha := fechaA.Format("2006-01-02 15:04:05")
	copy(disk.Mbr_fecha_creacion[:], fecha)

	disk.Mbr_partition_1.Part_status = [1]byte{'0'}
	disk.Mbr_partition_2.Part_status = [1]byte{'0'}
	disk.Mbr_partition_3.Part_status = [1]byte{'0'}
	disk.Mbr_partition_4.Part_status = [1]byte{'0'}

	disk.Mbr_partition_1.Part_fit = [1]byte{'0'}
	disk.Mbr_partition_2.Part_fit = [1]byte{'0'}
	disk.Mbr_partition_3.Part_fit = [1]byte{'0'}
	disk.Mbr_partition_4.Part_fit = [1]byte{'0'}

	disk.Mbr_partition_1.Part_type = [1]byte{'0'}
	disk.Mbr_partition_2.Part_type = [1]byte{'0'}
	disk.Mbr_partition_3.Part_type = [1]byte{'0'}
	disk.Mbr_partition_4.Part_type = [1]byte{'0'}

	disk.Mbr_partition_1.Part_start = 0
	disk.Mbr_partition_2.Part_start = 0
	disk.Mbr_partition_3.Part_start = 0
	disk.Mbr_partition_4.Part_start = 0

	disk.Mbr_partition_1.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition_2.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition_3.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition_4.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}

	buffer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		buffer.WriteByte(0)
	}

	var totalBytes int = 0
	for totalBytes < size {
		c, err := archivo.Write(buffer.Bytes())
		if err != nil {
			fmt.Println("Error: No se pudo escribir en el archivo")
			respuesta += "Error: No se pudo escribir en el archivo\n"
			return respuesta
		}
		totalBytes += c
	}
	fmt.Println("Archivo llenado")
	//Escribir el MBR en el archivo
	archivo.Seek(0, 0)
	if err != nil {
		fmt.Println("Error: No se pudo escribir en el archivo")
		respuesta += "Error: No se pudo escribir en el archivo\n"
		return respuesta
	}
	fmt.Println("Disco " + fileName + " creado correctamente")
	respuesta += "Disco " + fileName + " creado correctamente\n"
	return respuesta
}

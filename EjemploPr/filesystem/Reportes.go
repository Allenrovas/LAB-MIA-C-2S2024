package filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// respuesta += ReporteDisk(idValor, pathValor)
func ReporteDisk(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]

	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte Disk \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "node[shape=record, color=lightgrey]a0[label=\"MBR"
	//Leer el MBR
	disk := MBR{}
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		respuesta += "Error al leer el MBR\n"
		fmt.Println("Error al leer el MBR")
		return respuesta
	}
	sizeMBR := int(disk.Mbr_tamano)
	libreMBR := int(disk.Mbr_tamano)

	//Crear el MBR
	if disk.Mbr_partition_1.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_1.Part_size)
		Dot += "|"
		if disk.Mbr_partition_1.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_1.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_1.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_1.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}

	}
	if disk.Mbr_partition_2.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_2.Part_size)
		Dot += "|"
		if disk.Mbr_partition_2.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_2.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_2.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_2.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition_3.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_3.Part_size)
		Dot += "|"
		if disk.Mbr_partition_3.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_3.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_3.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_3.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}

	}
	if disk.Mbr_partition_4.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_4.Part_size)
		Dot += "|"
		if disk.Mbr_partition_4.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_4.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_4.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_4.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if libreMBR > 0 {
		Dot += "|Libre"
		porcentaje := (float64(libreMBR) * float64(100)) / float64(sizeMBR)
		Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
	}

	Dot += "\"];\n}"
	//Crear el archivo dot
	//-path=/home/user/reports/report2.pdf
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
	//Crear el archivo .dot
	file, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot")
		respuesta += "Error al crear el archivo .dot\n"
		return respuesta
	}
	defer file.Close()
	//Escribir el archivo .dot
	_, err = file.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot")
		respuesta += "Error al escribir el archivo .dot\n"
		return respuesta
	}
	fmt.Println("Archivo .dot creado")

	//Quitar el punto a la extension
	extension = extension[1:]

	//Crear el reporte
	cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
	fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al crear el reporte")
		respuesta += "Error al crear el reporte\n"
		return respuesta
	}

	return "Reporte Disk creado con exito\n"
}

func ReporteSB(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte SuperBlock \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">SuperBlock</TD><TD></TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_filesystem_type</TD><TD>" + strconv.Itoa(int(superBloque.S_filesystem_type)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inodes_count</TD><TD>" + strconv.Itoa(int(superBloque.S_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_blocks_count</TD><TD>" + strconv.Itoa(int(superBloque.S_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_blocks_count</TD><TD>" + strconv.Itoa(int(superBloque.S_free_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_inodes_count</TD><TD>" + strconv.Itoa(int(superBloque.S_free_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mtime</TD><TD>" + string(superBloque.S_mtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_umtime</TD><TD>" + string(superBloque.S_umtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mnt_count</TD><TD>" + strconv.Itoa(int(superBloque.S_mnt_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_magic</TD><TD>" + strconv.Itoa(int(superBloque.S_magic)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_size</TD><TD>" + strconv.Itoa(int(superBloque.S_inode_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_size</TD><TD>" + strconv.Itoa(int(superBloque.S_block_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_ino</TD><TD>" + strconv.Itoa(int(superBloque.S_first_ino)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_blo</TD><TD>" + strconv.Itoa(int(superBloque.S_first_blo)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_inode_start</TD><TD>" + strconv.Itoa(int(superBloque.S_bm_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_block_start</TD><TD>" + strconv.Itoa(int(superBloque.S_bm_block_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_start</TD><TD>" + strconv.Itoa(int(superBloque.S_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_start</TD><TD>" + strconv.Itoa(int(superBloque.S_block_start)) + "</TD></TR>\n"
	Dot += "</TABLE>>];\n"
	Dot += "}"

	//Crear el archivo dot
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
	//Crear el archivo .dot
	file, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot")
		respuesta += "Error al crear el archivo .dot\n"
		return respuesta
	}
	defer file.Close()
	//Escribir el archivo .dot
	_, err = file.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot")
		respuesta += "Error al escribir el archivo .dot\n"
		return respuesta
	}
	fmt.Println("Archivo .dot creado")

	//Quitar el punto a la extension
	extension = extension[1:]

	//Crear el reporte
	cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
	fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al crear el reporte")
		respuesta += "Error al crear el reporte\n"
		return respuesta
	}

	return "Reporte SuperBlock creado con exito\n"

}

func ReporteFile(idValor string, pathValor string, rutaValor string) string {
	var respuesta string
	//Buscar la particion montada
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	indice := VerificarParticionMontada(idValor)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}
	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	//Buscar el inodo de la ruta
	numeroInodo := BuscarInodo(rutaValor, MountActual, superBloque, archivo)
	if numeroInodo == -1 {
		respuesta += "La ruta no existe\n"
		return respuesta
	}
	//Leer el inodo
	cadena := LeerArchivo(numeroInodo, superBloque, archivo)
	if len(cadena) == 0 {
		respuesta += "El archivo esta vacio\n"
		return respuesta
	}
	Dot := "digraph G{\n"
	Dot += "a[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD colspan=\"2\" bgcolor=\"lightgrey\" >" + rutaValor + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Contenido</TD></TR>\n"
	Dot += "<TR><TD>" + cadena + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"
	//Crear el archivo dot
	extension := path.Ext(pathValor)

	if extension == ".txt" {
		//Crear el archivo .txt
		file, err := os.Create(pathValor)
		if err != nil {
			fmt.Println("Error al crear el archivo .txt")
			respuesta += "Error al crear el archivo .txt\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .txt
		_, err = file.WriteString(cadena)
		if err != nil {
			fmt.Println("Error al escribir el archivo .txt")
			respuesta += "Error al escribir el archivo .txt\n"
			return respuesta
		}
		fmt.Println("Archivo .txt creado")
		return "Reporte File creado con exito\n"

	} else {
		//Archivo sin extension
		fileName = strings.TrimSuffix(fileName, extension)
		DotName := dirPath + fileName + ".dot"
		//Crear el archivo .dot
		file, err := os.Create(DotName)
		if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .dot
		_, err = file.WriteString(Dot)
		if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
		}
		fmt.Println("Archivo .dot creado")

		//Quitar el punto a la extension
		extension = extension[1:]

		//Crear el reporte
		cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
		fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error al crear el reporte")
			respuesta += "Error al crear el reporte\n"
			return respuesta
		}

		return "Reporte File creado con exito\n"
	}

}

func ReporteBMInode(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	//Leer el bitmap de inodos, teniendo 20 registros por fila
	Desplazamiento := int(superBloque.S_bm_inode_start)
	BmString := ""

	for i := 0; i < int(superBloque.S_inodes_count); i++ {
		var bit byte
		archivo.Seek(int64(Desplazamiento+i), 0)
		err = binary.Read(archivo, binary.LittleEndian, &bit)
		if err != nil {
			respuesta += "Error al leer el bitmap de inodos\n"
			fmt.Println("Error al leer el bitmap de inodos")
			return respuesta
		}
		if bit == 0 {
			BmString += "0"
		} else {
			BmString += "1"
		}
		if (i+1)%20 == 0 {
			BmString += "\n"
		}
	}
	Dot := "digraph G{\n"
	Dot += "a[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Bitmap de Inodos</TD></TR>\n"
	Dot += "<TR><TD>" + BmString + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"
	//Crear el archivo dot
	extension := path.Ext(pathValor)
	if extension == ".txt" {
		//Crear el archivo .txt
		file, err := os.Create(pathValor)
		if err != nil {
			fmt.Println("Error al crear el archivo .txt")
			respuesta += "Error al crear el archivo .txt\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .txt
		_, err = file.WriteString(BmString)
		if err != nil {
			fmt.Println("Error al escribir el archivo .txt")
			respuesta += "Error al escribir el archivo .txt\n"
			return respuesta
		}
		fmt.Println("Archivo .txt creado")
		return "Reporte Bitmap de Inodos creado con exito\n"

	} else {
		//Archivo sin extension
		fileName = strings.TrimSuffix(fileName, extension)
		DotName := dirPath + fileName + ".dot"
		//Crear el archivo .dot
		file, err := os.Create(DotName)
		if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .dot
		_, err = file.WriteString(Dot)
		if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
		}
		fmt.Println("Archivo .dot creado")

		//Quitar el punto a la extension
		extension = extension[1:]

		//Crear el reporte
		cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
		fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error al crear el reporte")
			respuesta += "Error al crear el reporte\n"
			return respuesta
		}

		return "Reporte Bitmap de Inodos creado con exito\n"
	}
}

func ReporteBMBlock(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	//Leer el bitmap de bloques, teniendo 20 registros por fila
	Desplazamiento := int(superBloque.S_bm_block_start)
	BmString := ""

	for i := 0; i < int(superBloque.S_blocks_count); i++ {
		var bit byte
		archivo.Seek(int64(Desplazamiento+i), 0)
		err = binary.Read(archivo, binary.LittleEndian, &bit)
		if err != nil {
			respuesta += "Error al leer el bitmap de bloques\n"
			fmt.Println("Error al leer el bitmap de bloques")
			return respuesta
		}
		if bit == 0 {
			BmString += "0"
		} else {
			BmString += "1"
		}
		if (i+1)%20 == 0 {
			BmString += "\n"
		}
	}
	Dot := "digraph G{\n"
	Dot += "a[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Bitmap de Bloques</TD></TR>\n"
	Dot += "<TR><TD>" + BmString + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"
	//Crear el archivo dot
	extension := path.Ext(pathValor)
	if extension == ".txt" {
		//Crear el archivo .txt
		file, err := os.Create(pathValor)
		if err != nil {
			fmt.Println("Error al crear el archivo .txt")
			respuesta += "Error al crear el archivo .txt\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .txt
		_, err = file.WriteString(BmString)
		if err != nil {
			fmt.Println("Error al escribir el archivo .txt")
			respuesta += "Error al escribir el archivo .txt\n"
			return respuesta
		}
		fmt.Println("Archivo .txt creado")
		return "Reporte Bitmap de Bloques creado con exito\n"

	} else {
		//Archivo sin extension
		fileName = strings.TrimSuffix(fileName, extension)
		DotName := dirPath + fileName + ".dot"
		//Crear el archivo .dot
		file, err := os.Create(DotName)
		if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .dot
		_, err = file.WriteString(Dot)
		if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
		}
		fmt.Println("Archivo .dot creado")

		//Quitar el punto a la extension
		extension = extension[1:]

		//Crear el reporte
		cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
		fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error al crear el reporte")
			respuesta += "Error al crear el reporte\n"
			return respuesta
		}

		return "Reporte Bitmap de Bloques creado con exito\n"
	}
}

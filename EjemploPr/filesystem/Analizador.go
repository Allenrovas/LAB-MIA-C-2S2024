package filesystem

import (
	"fmt"
	"strconv"
	"strings"
)

// DividirComando recibe un comando y lo divide en un arreglo de strings
func DividirComando(comando string) string {
	var respuesta string
	//se divide el comando en un arreglo de strings por el enter
	comandos := strings.Split(comando, "\n")
	//se recorre el arreglo de strings
	for i := 0; i < len(comandos); i++ {
		//imprime el comando
		fmt.Println("Comando: " + comandos[i])
		//se analiza el comando
		respuesta += AnalizarComando(comandos[i])
	}
	return respuesta
}

// AnalizarComando recibe un comando y lo analiza
func AnalizarComando(comando string) string {
	var respuesta string
	//se divide el comando en un arreglo de strings por el espacio
	comandoSeparado := strings.Split(comando, " ")
	//Si encuentra el # en la posicion 0, es un comentario
	if strings.Contains(comandoSeparado[0], "#") {
		//imprime el comentario
		fmt.Println("Comentario: ")
		//Eliminiar el #
		comandoSeparado[0] = strings.Replace(comandoSeparado[0], "#", "", -1)
		respuesta += "Comentario: "
		//se recorre el arreglo de strings
		for i := 0; i < len(comandoSeparado); i++ {
			respuesta += comandoSeparado[i] + " "
		}
		respuesta += "\n"
		fmt.Println(respuesta)
	} else {
		//Si no es un comentario, entonces es un comando
		//Iterar sobre el comando
		for _, valor := range comandoSeparado {
			//el primer valor del comando lo pasamos a minusculas
			valor = strings.ToLower(valor)
			//Si el valor es mkdisk, entonces es un comando para crear un disco
			if valor == "mkdisk" {
				fmt.Println("Comando mkdisk")
				respuesta += "Ejecutando mkdisk\n"
				//Analizar el comando mkdisk
				respuesta += AnalizarMkdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "rmdisk" {
				fmt.Println("Comando rmdisk")
				respuesta += "Ejecutando rmdisk\n"
				//Analizar el comando rmdisk
				//respuesta += AnalizarRmdisk(comandoSeparado)
				//Pasar a string el comando separado
				//comandoSeparadoString := strings.Join(comandoSeparado, " ")
				//respuesta += AnalizarComando(comandoSeparadoString)
				//return respuesta
			} else if valor == "\n" {
				continue
			} else if valor == "\r" {
				continue
			} else if valor == "\t" {
				continue
			} else if valor == "" {
				continue
			} else {
				//Si no es ningun comando, entonces es un error
				fmt.Println("Error: Comando no reconocido")
				respuesta += "Error: Comando no reconocido\n"
			}
		}
	}
	return respuesta
}

// AnalizarMkdisk recibe un comando mkdisk y lo analiza
func AnalizarMkdisk(comando *[]string) string {
	//mkdisk -unit=M -path="/home 1/mis discos/Disco3.mia"
	//0 		1     2     3"/home/mis     4         5discos/Disco3.mia"
	var respuesta string
	//Eliminar el primer valor del comando
	*comando = (*comando)[1:]
	//-size=5 -unit=M -path="/home/mis discos/Disco3.mia"
	//Booleanos para saber si se encontro el size, unit, fit, path
	var size, unit, path, fit bool
	//Variables para guardar el valor del size, unit, fit, path
	var valorSize, valorUnit, valorFit, valorPath string
	//Iterar sobre el comando
	valorFit = "f"
	valorUnit = "m"
	for _, valor := range *comando {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-size" {
			size = true
			valorSize = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-unit" {
			unit = true
			valorUnit = banderaValor
			valorUnit = strings.ToLower(valorUnit)
			*comando = (*comando)[1:]
		} else if bandera == "-fit" {
			fit = true
			valorFit = banderaValor
			valorFit = strings.ToLower(valorFit)
			*comando = (*comando)[1:]
		} else if bandera == "-path" {
			path = true
			//Verificar si el path tiene comillas
			//-path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
				//Eliminar las comillas del inicio
				banderaValor = strings.Replace(banderaValor, "\"", "", -1)
				//Eliminar el primer valor del comandoSeparado
				*comando = (*comando)[1:]
				//Iterar sobre el comando
				Contador := 0
				for _, valor := range *comando {
					//Si el valor contiene comillas
					if strings.Contains(valor, "\"") {
						//Eliminar las comillas del final
						valor = strings.Replace(valor, "\"", "", -1)
						//Agregar el valor al path
						valorPath += valor
						break
					} else {
						//Agregar el valor al path
						valorPath += valor + " "
						Contador++
					}
				}
				//Eliminar los valores del comando
				*comando = (*comando)[Contador:]
			} else {
				valorPath = banderaValor
				*comando = (*comando)[1:]
			}
		} else {
			fmt.Println("Error: Parametro no reconocida")
			respuesta += "Error: Parametro no reconocida\n"
		}

	}
	if !size {
		fmt.Println("Error: Falta el parametro size")
		respuesta += "Error: Falta el parametro size\n"
		return respuesta
	} else if !path {
		fmt.Println("Error: Falta el parametro path")
		respuesta += "Error: Falta el parametro path\n"
		return respuesta
	} else {
		if fit {
			if valorFit != "bf" && valorFit != "ff" && valorFit != "wf" {
				fmt.Println("Error: Fit no reconocido")
				respuesta += "Error: Fit no reconocido\n"
				return respuesta
			} else {
				if valorFit == "bf" {
					valorFit = "b"
				} else if valorFit == "ff" {
					valorFit = "f"
				} else if valorFit == "wf" {
					valorFit = "w"
				}
			}
		}
		if unit {
			if valorUnit != "k" && valorUnit != "m" {
				fmt.Println("Error: Unit no reconocido")
				respuesta += "Error: Unit no reconocido\n"
				return respuesta
			}
		}
		//Pasar a int el size
		sizeInt, err := strconv.Atoi(valorSize)
		if err != nil {
			fmt.Println("Error: Size no es un numero")
			respuesta += "Error: Size no es un numero\n"
			return respuesta
		}
		//Verificar que el size sea mayor a 0
		if sizeInt <= 0 {
			fmt.Println("Error: Size debe ser mayor a 0")
			respuesta += "Error: Size debe ser mayor a 0\n"
			return respuesta
		}
		//Imprimir los valores
		fmt.Println("Size: " + valorSize)
		fmt.Println("Unit: " + valorUnit)
		fmt.Println("Fit: " + valorFit)
		fmt.Println("Path: " + valorPath)
		//Llamar a la funcion para crear el disco
		respuesta += CrearDiscos(sizeInt, valorUnit, valorFit, valorPath)
		return respuesta
	}
}

func obtenerBandera(bandera string) string {
	//-size
	var banderaValor string
	for _, valor := range bandera {
		if valor == '=' {
			break
		}
		banderaValor += string(valor)
	}
	banderaValor = strings.ToLower(banderaValor)
	return banderaValor
}

func obtenerValor(bandera string) string {
	var banderaValor string
	var boolBandera bool
	for _, valor := range bandera {
		if boolBandera {
			banderaValor += string(valor)
		}
		if valor == '=' {
			boolBandera = true
		}
	}
	return banderaValor
}

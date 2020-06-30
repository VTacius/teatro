package utils

import (
	"strings"
)

// Conexion : El punto de entrada para configurarse desde la construcción
var Conexion string

// ObtenerCadenaConexion : La cadena no puede tener espacios, según parece
func ObtenerCadenaConexion() string {
	resultado := strings.Replace(Conexion, ";", " ", -1)
	return resultado
}

package utils

import (
	"fmt"
	"os"
)

// Credenciales El punto de entrada para las contraseñas
var Credenciales string

// Salida Se asegura de salir con dignidad
func Salida(mensaje string, err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", mensaje, err)

	os.Exit(1)
}

// ObtenerCredenciales : Para estandarizar un poco esto
func ObtenerCredenciales() (string, string) {
	usuario := "uid=zimbra,cn=admins,cn=zimbra"
	return usuario, Credenciales
}

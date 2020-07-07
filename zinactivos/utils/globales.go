package utils

import (
	"fmt"
	"os"
)

// UsuarioLdap El punto de entrada para el usuario
var UsuarioLdap string

// PasswordLdap El punto de entrada para el usuario
var PasswordLdap string

// ServidorLdap El punto de entrada para el usuario
var ServidorLdap string

// UsuarioCorreo El punto de entrada para el usuario de envio de correos
var UsuarioCorreo string

// PasswordCorreo Punto de entrada para la contraseña de correo
var PasswordCorreo string

// ServidorCorreo P.E para la configuración del servidor de correo
var ServidorCorreo string

// Salida Se asegura de salir con dignidad
func Salida(mensaje string, err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", mensaje, err)

	os.Exit(1)
}

// ConfiguracionAccesoLdap : Para estandarizar un poco esto
func ConfiguracionAccesoLdap() (string, string, string) {
	url := "ldap://" + ServidorLdap + ":389"
	return UsuarioLdap, PasswordLdap, url 
}

// ConfiguracionEnvioCorreo : Todo lo necesario para el envío de Correo
func ConfiguracionEnvioCorreo() (string, string, string) {
	return UsuarioCorreo, PasswordCorreo, ServidorCorreo
}

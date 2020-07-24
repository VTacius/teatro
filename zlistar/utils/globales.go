package utils

// UsuarioLdap : DN de usuario con el que se conecta
var UsuarioLdap string

// PasswordLdap : Contraseña de usuario para conectarse
var PasswordLdap string

// ServidorLdap : Servidor LDAP
var ServidorLdap string

// ParametrosAccesoLdap : Los valores necesarios para conectarnos
func ParametrosAccesoLdap() (string, string, string) {
	url := "ldap://" + ServidorLdap + ":389"
	return url, UsuarioLdap, PasswordLdap
}

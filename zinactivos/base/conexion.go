package base

import (
	"crypto/tls"

	"github.com/go-ldap/ldap/v3"
)

// Reemplazo : Un simple clave-valor sobre los atributos a reemplazar
type Reemplazo struct {
	Clave string
	Valor string
}

// Conectar : Inicia comunicaciones con el servidor LDAP
func Conectar(url string, usuario string, contrasenia string) (*ldap.Conn, error) {
	// Conexión inicial
	config := ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	conexion, err := ldap.DialURL(url, config)
	if err != nil {
		return conexion, err
	}

	// Nos autenticamos
	err = conexion.Bind(usuario, contrasenia)
	if err != nil {
		return conexion, err
	}

	return conexion, nil
}

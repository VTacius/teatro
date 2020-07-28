package base

import (
	"crypto/tls"

	"github.com/go-ldap/ldap/v3"
)

// Conectar : Inicia la comunicación con el servidor LDAP
func Conectar(url string, usuario string, contrasenia string) (*ldap.Conn, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	config := ldap.DialWithTLSConfig(tlsConfig)
	conexion, err := ldap.DialURL(url, config)
	if err != nil {
		return conexion, err
	}

	conexion.StartTLS(tlsConfig)
	err = conexion.Bind(usuario, contrasenia)
	if err != nil {
		return conexion, err
	}

	return conexion, nil
}

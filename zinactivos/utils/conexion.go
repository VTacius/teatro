package utils

import (
	"crypto/tls"

	"github.com/go-ldap/ldap/v3"
)

// Peticionar : Realiza una petición al servidor LDAP
func Peticionar(conexion *ldap.Conn, base string, filtro string, atributos []string) (*ldap.SearchResult, error) {

	peticion := ldap.NewSearchRequest(base,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false, filtro, atributos, nil)

	respuesta, err := conexion.Search(peticion)
	if err != nil {
		return respuesta, err
	}

	return respuesta, nil
}

// Reemplazo : Un simple clave-valor sobre los atributos a reemplazar
type Reemplazo struct {
	Clave string
	Valor string
}

// ModificarReemplazo : Realiza una modificación en base al reemplazo de un atributo
func ModificarReemplazo(conexion *ldap.Conn, dn string, reemplazos []Reemplazo) error {
	modificacion := ldap.NewModifyRequest(dn, []ldap.Control{})
	for _, r := range reemplazos {
		modificacion.Replace(r.Clave, []string{r.Valor})
	}

	err := conexion.Modify(modificacion)

	return err
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

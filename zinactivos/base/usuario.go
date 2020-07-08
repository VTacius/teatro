package base

import (
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
	"xibalba.com/vtacius/zinactivos/utils"
)

// Usuario Contiene la definición actual de usuario
type Usuario struct {
	DN                       string
	Username                 string
	Nombre                   string
	ZimbraLastLogonTimestamp time.Time
	Description              string
}

// AccesoUsuario representa las acciones que realizamos sobre alguna base
type AccesoUsuario struct {
	url      string
	base     string
	username string
	password string
	Datos    []Usuario
}

// NewAccesoUsuario Factoria para AccesoUsuario
func NewAccesoUsuario(url string, username string, password string) *AccesoUsuario {
	return &AccesoUsuario{
		url:      url,
		base:     "ou=people,dc=salud,dc=gob,dc=sv",
		username: username,
		password: password,
		Datos:    []Usuario{},
	}
}

type transformador func(*ldap.SearchResult) []Usuario

// ListarUsuarios Lista los usuarios según filtro
func (a *AccesoUsuario) ListarUsuarios(filtro string, atributos []string, fn transformador) {
	conexion, err := utils.Conectar(a.url, a.username, a.password)
	if err != nil {
		utils.Salida("Error en la conexión", err)
	}

	respuesta, err := utils.Peticionar(conexion, a.base, filtro, atributos)
	if err != nil {
		utils.Salida("Error en la peticion", err)
	}

	a.Datos = fn(respuesta)

}

// ModificarUsuario Modificamos UN usuario
func (a *AccesoUsuario) ModificarUsuario(dn string, modificacion []utils.Reemplazo) {
	conexion, err := utils.Conectar(a.url, a.username, a.password)
	if err != nil {
		utils.Salida("Error en la conexión", err)
	}

	err = utils.ModificarReemplazo(conexion, dn, modificacion)
	if err != nil {
		utils.Salida(fmt.Sprintf("Error modificando %s", dn), err)
	}

}

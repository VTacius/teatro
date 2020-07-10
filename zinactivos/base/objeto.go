package base

import (
	"time"

	"github.com/go-ldap/ldap/v3"
	"xibalba.com/vtacius/zinactivos/utils"
)

type operacional func(*ldap.SearchResult) []Objeto

// Objeto es el tipo para mockeo
type Objeto struct {
	DN                       string
	Username                 string
	Nombre                   string
	ZimbraLastLogonTimestamp time.Time
	Description              string
}

// ObjetoAcceso es una prueba sobre el Mockeo
type ObjetoAcceso struct {
	Base    string
	Cliente ldap.Client
	Datos   []Objeto
}

// ListarUsuarios Lista los usuarios según filtro
func (a *ObjetoAcceso) ListarUsuarios(filtro string, atributos []string, fn operacional) error {

	peticion := ldap.NewSearchRequest(a.Base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filtro, atributos, nil)
	respuesta, err := a.Cliente.Search(peticion)
	if err != nil {
		return err
	}

	a.Datos = fn(respuesta)
	return nil
}

// ModificarUsuario Modificamos UN usuario
func (a *ObjetoAcceso) ModificarUsuario(dn string, reemplazos []utils.Reemplazo) error {

	modificacion := ldap.NewModifyRequest(dn, []ldap.Control{})
	for _, reemplazo := range reemplazos {
		modificacion.Replace(reemplazo.Clave, []string{reemplazo.Valor})
	}

	err := a.Cliente.Modify(modificacion)
	if err != nil {
		return err
	}

	return nil
}

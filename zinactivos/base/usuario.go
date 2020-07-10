package base

import (
	"strings"
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
	Data    *ldap.SearchResult
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

func obtenerValor(entrada *ldap.Entry, clave string) (bool, string) {
	valor := strings.TrimSpace(entrada.GetAttributeValue(clave))

	return valor != "", valor

}

// ObtenerUsuarios Lista los usuarios según filtro
func (a *ObjetoAcceso) ObtenerUsuarios(filtro string, atributos []string) error {

	peticion := ldap.NewSearchRequest(a.Base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filtro, atributos, nil)
	respuesta, err := a.Cliente.Search(peticion)
	if err != nil {
		return err
	}

	a.Data = respuesta
	return nil
}

// ListarBusqueda : Convierte un ldap.SearchResult a un []Objectos
func (a *ObjetoAcceso) ListarBusqueda() *ObjetoAcceso {
	resultado := []Objeto{}

	for _, item := range a.Data.Entries {
		usuario := Objeto{}
		usuario.DN = item.DN
		usuario.ZimbraLastLogonTimestamp = utils.Fechador(item.GetAttributeValue("zimbraLastLogonTimestamp"))
		if existe, valor := obtenerValor(item, "uid"); existe {
			usuario.Username = valor
		}
		if existe, valor := obtenerValor(item, "cn"); existe {
			usuario.Nombre = valor
		}
		if existe, valor := obtenerValor(item, "description"); existe {
			usuario.Description = valor
		}
		resultado = append(resultado, usuario)
	}
	a.Datos = resultado
	return a
}

// FiltrarInactivos : Filtra aquellas entradas que no se hayan logueado en tanto tiempo
func (a *ObjetoAcceso) FiltrarInactivos(periodoInactividad float64) (resultado []Objeto) {
	hoy := time.Now()
	for _, item := range a.Datos {
		if utils.RevisarIntervalo(periodoInactividad, hoy, item.ZimbraLastLogonTimestamp) == 1 {
			resultado = append(resultado, item)
		}
	}

	return
}

// ModificarUsuario Modificamos UN usuario
func (a *ObjetoAcceso) ModificarUsuario(dn string, reemplazos []Reemplazo) error {

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

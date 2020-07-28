package base

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// Acceso : Establece la conexión y operaciones LDAP
type Acceso struct {
	Cliente ldap.Client
	Base    string
	attrs   []string
	Err     error
	Datos   []Objeto
	Data    *ldap.SearchResult
}

// Buscar : Lista todos los objetos tales solicitados
func (acceso *Acceso) Buscar(filtro string, atributos []string) *Acceso {
	acceso.attrs = atributos
	peticion := ldap.NewSearchRequest(acceso.Base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filtro, atributos, nil)
	respuesta, err := acceso.Cliente.Search(peticion)
	if err != nil {
		acceso.Err = err
		return acceso
	}

	acceso.Data = respuesta
	return acceso
}

func obtenerObjeto(entrada *ldap.Entry, atributos []string) map[string]Atributo {
	resultado := make(map[string]Atributo)
	for _, clave := range atributos {
		valor := strings.TrimSpace(entrada.GetAttributeValue(clave))
		if valor != "" {
			resultado[clave] = Atributo{valor, len(valor)}
		}
	}
	return resultado
}

// Listar : Devuelve las entradas en un formato más accesible
func (acceso *Acceso) Listar() *Acceso {
	if acceso.Err != nil {
		return acceso
	}
	var resultado []Objeto
	for _, item := range acceso.Data.Entries {
		contenido := obtenerObjeto(item, acceso.attrs)
		resultado = append(resultado, Objeto{item.DN, contenido})
	}

	acceso.Datos = resultado

	return acceso
}

func obtenerLongitudes(attrs *[]string, datos *[]Objeto) map[string]int {
	longitudes := make(map[string]int)
	for _, clave := range *attrs {
		longitudes[clave] = 0
	}

	for _, item := range *datos {
		for _, clave := range *attrs {
			if item.Atributos[clave].Longitud > longitudes[clave] {
				longitudes[clave] = item.Atributos[clave].Longitud
			}
		}
	}
	return longitudes
}

// ParaCSV : Produce una salida en CSV
func (acceso *Acceso) ParaCSV(salida io.Writer) {
	for _, item := range acceso.Datos {
		fmt.Fprintln(salida, item.Enumerar(acceso.attrs))
	}
}

// Imprimir : Muestra en pantalla el resultado
func (acceso *Acceso) Imprimir(salida io.Writer) {
	longitudes := obtenerLongitudes(&acceso.attrs, &acceso.Datos)
	for _, item := range acceso.Datos {
		fmt.Fprintln(salida, item.Tabular(acceso.attrs, longitudes))
	}
}

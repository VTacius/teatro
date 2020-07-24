package base

import (
	"fmt"
	"strings"
)

// Atributo : Valor y metadata por cada atributos
type Atributo struct {
	Valor    string
	Longitud int
}

// Objeto : Cada item que se obtiene de LDAP
type Objeto struct {
	DN        string
	Atributos map[string]Atributo
}

// Enumerar : Crea una lista de los valores del item
func (objeto *Objeto) Enumerar(atributos []string) string {
	var resultado strings.Builder
	for _, clave := range atributos {
		attr := "'" + objeto.Atributos[clave].Valor + "';"
		resultado.WriteString(attr)
	}
	return resultado.String()
}

// Tabular : Crea una simpática tabla para mostrar
func (objeto *Objeto) Tabular(atributos []string, longitudes map[string]int) string {
	var resultado strings.Builder
	resultado.WriteString("|")
	for _, clave := range atributos {
		celda := fmt.Sprintf(" %-*s |", longitudes[clave], objeto.Atributos[clave].Valor)
		resultado.WriteString(celda)
	}

	return resultado.String()
}

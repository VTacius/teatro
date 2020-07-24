package base_test

import (
	"testing"

	"xibalba.com/vtacius/zlistar/base"
)

func TestTabular(t *testing.T) {
	for _, caso := range casosTabular {
		objeto := base.Objeto{"", caso.atributos}
		resultado := objeto.Tabular(caso.attrs, caso.longitudes)
		if resultado != caso.esperado {
			t.Errorf("Error en %s: %s != %s", caso.mensaje, resultado, caso.esperado)
		}

	}
}

func TestEnumerar(t *testing.T) {
	for _, caso := range casosEnumerar {
		objeto := base.Objeto{"", caso.atributos}
		resultado := objeto.Enumerar(caso.attrs)
		if resultado != caso.esperado {
			t.Errorf("Error en %s: %s != %s", caso.mensaje, resultado, caso.esperado)
		}
	}
}

var casosTabular = []struct {
	attrs      []string
	longitudes map[string]int
	atributos  map[string]base.Atributo
	esperado   string
	mensaje    string
}{
	{
		[]string{"atributo"},
		map[string]int{
			"atributo": 5,
		},
		map[string]base.Atributo{
			"atributo": base.Atributo{"valor", 5},
		},
		"| valor |",
		"Caso mínimo",
	},
	{
		[]string{"atributo", "attrs"},
		map[string]int{
			"atributo": 7,
			"attrs":    7,
		},
		map[string]base.Atributo{
			"atributo": base.Atributo{"valor", 5},
			"attrs":    base.Atributo{"value", 5},
		},
		"| valor   | value   |",
		"Verifica que el tamaño dependa de longitudes",
	},
}

var casosEnumerar = []struct {
	attrs     []string
	atributos map[string]base.Atributo
	esperado  string
	mensaje   string
}{
	{
		[]string{"atributo"},
		map[string]base.Atributo{
			"atributo": base.Atributo{"valor", 5},
		},
		"'valor';",
		"Caso mínimo",
	},
	{
		[]string{"atributo", "attrs"},
		map[string]base.Atributo{
			"atributo": base.Atributo{"valor", 5},
			"attrs":    base.Atributo{"value", 5},
		},
		"'valor';'value';",
		"Primer caso funcional",
	},
}

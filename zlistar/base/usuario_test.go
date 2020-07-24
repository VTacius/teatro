package base_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"xibalba.com/vtacius/zlistar/base"
)

func TestObtenerLongitudes(t *testing.T) {
	for _, caso := range casosObtenerLongitudes {
		resultado := base.ObtenerLongitudes(&caso.listaAtributos, &caso.datos)
		if !cmp.Equal(resultado, caso.esperado) {
			t.Errorf("Error con %s", caso.mensaje)
			t.Error(cmp.Diff(resultado, caso.esperado))
		}
	}
}

var casosObtenerLongitudes = []struct {
	datos          []base.Objeto
	listaAtributos []string
	esperado       map[string]int
	mensaje        string
}{
	{
		listaAtributos: []string{"atributo", "attr"},
		datos: []base.Objeto{
			base.Objeto{
				Atributos: map[string]base.Atributo{
					"atributo": base.Atributo{"valor", 5},
					"attr":     base.Atributo{"value", 5},
				},
			},
		},
		esperado: map[string]int{"atributo": 5, "attr": 5},
		mensaje:  "Caso base",
	},
}

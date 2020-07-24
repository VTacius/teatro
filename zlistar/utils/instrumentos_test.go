package utils_test

import (
	"testing"

	"xibalba.com/vtacius/zlistar/utils"
)

func TestConstruirBase(t *testing.T) {
	for _, caso := range casosConstruirBase {
		resultado := utils.ConstruirBase(caso.base)
		if resultado != caso.esperado {
			t.Errorf("El caso %s tiene problemas: %s != %s", caso.mensaje, resultado, caso.esperado)
		}
	}
}

var casosConstruirBase = []struct {
	base     string
	esperado string
	mensaje  string
}{
	{"sv", "dc=sv", "El caso por defecto"},
	{"hnm.gob.sv", "dc=hnm,dc=gob,dc=sv", "El caso por defecto"},
}

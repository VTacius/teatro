package utils_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"xibalba.com/vtacius/zlistar/utils"
)

func TestParametrosAccesoLdap(t *testing.T) {
	for _, caso := range casosParametros {
		utils.UsuarioLdap = caso.usuario
		utils.PasswordLdap = caso.contrasenia
		utils.ServidorLdap = caso.servidor
		url, user, pass := utils.ParametrosAccesoLdap()
		respuesta := []string{url, user, pass}
		if !cmp.Equal(respuesta, caso.esperado) {
			t.Errorf("Problemas con %s", caso.mensaje)
			t.Error(cmp.Diff(respuesta, caso.esperado))
		}
	}
}

var casosParametros = []struct {
	usuario     string
	contrasenia string
	servidor    string
	esperado    []string
	mensaje     string
}{
	{
		"alortiz", "passwordito", "192.168.2.1", []string{"ldap://192.168.2.1:389", "alortiz", "passwordito"}, "Caso base",
	},
}

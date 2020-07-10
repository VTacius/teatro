package base_test

import (
	"fmt"
	"testing"

	"github.com/go-ldap/ldap/v3"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"xibalba.com/vtacius/zinactivos/base"
	"xibalba.com/vtacius/zinactivos/mocks"
)

func operacion(entradas *ldap.SearchResult) (resultado []base.Objeto) {
	for _, entrada := range entradas.Entries {
		nombre := entrada.GetAttributeValue("cn")
		descripcion := entrada.GetAttributeValue("description")
		item := base.Objeto{Nombre: nombre, Description: descripcion}
		fmt.Println(item)
		resultado = append(resultado, item)
	}
	return
}

var datosEntrada = []map[string][]string{
	{
		"cn":          []string{"Alexander Ortíz"},
		"description": []string{"USERMODWEB"},
	},
}

func crearResultado(datos []map[string][]string) *ldap.SearchResult {
	resultado := new(ldap.SearchResult)

	for _, datosItem := range datos {
		entrada := ldap.NewEntry("", datosItem)
		resultado.Entries = append(resultado.Entries, entrada)

	}
	return resultado
}

func TestAccesoUsuario(t *testing.T) {
	ctrl := gomock.NewController(t)

	conexionMock := mocks.NewMockClient(ctrl)
	acceso := base.ObjetoAcceso{Cliente: conexionMock, Base: ""}

	peticion := ldap.NewSearchRequest("", ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "", []string{}, nil)

	resultado := crearResultado(datosEntrada)
	conexionMock.
		EXPECT().
		Search(peticion).
		Return(resultado, nil).
		Times(1)

	// Comienza propiamente el test
	item := base.Objeto{Nombre: "Alexander Ortíz", Description: "USERMODWEB"}
	esperado := []base.Objeto{item}
	acceso.ListarUsuarios("", []string{}, operacion)
	if !cmp.Equal(acceso.Datos, esperado) {
		t.Error(cmp.Diff(acceso.Datos, esperado))
	}
}

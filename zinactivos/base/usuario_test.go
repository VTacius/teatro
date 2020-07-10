package base_test

import (
	"errors"
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

func estaConfigurado(clave string, item map[string][]string) (bool, string) {
	valor, ok := item[clave]
	if ok {
		return ok, valor[0]
	}
	return ok, ""
}

func crearEsperado(datos []map[string][]string) (resultado []base.Objeto) {
	for _, item := range datos {
		usuario := base.Objeto{}
		if ok, valor := estaConfigurado("cn", item); ok {
			usuario.Nombre = valor
		}
		if ok, valor := estaConfigurado("description", item); ok {
			usuario.Description = valor
		}
		resultado = append(resultado, usuario)
	}

	return
}

func crearResultado(datos []map[string][]string) *ldap.SearchResult {
	resultado := new(ldap.SearchResult)

	for _, datosItem := range datos {
		entrada := ldap.NewEntry("", datosItem)
		resultado.Entries = append(resultado.Entries, entrada)

	}
	return resultado
}

func TestAccesoUsuarioNil(t *testing.T) {
	ctrl := gomock.NewController(t)

	conexionMock := mocks.NewMockClient(ctrl)
	acceso := base.ObjetoAcceso{Cliente: conexionMock, Base: ""}

	peticion := ldap.NewSearchRequest("", ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "", []string{}, nil)
	errorPrueba := errors.New("Error improbable")

	conexionMock.
		EXPECT().
		Search(peticion).
		Return(new(ldap.SearchResult), errorPrueba).
		Times(1)

	// Empieza el dichoso test
	err := acceso.ListarUsuarios("", []string{}, operacion)
	if err == nil {
		t.Error("Debe saber devolver un error")
	}

}

func TestModificarUsuario(t *testing.T) {
	ctrl := gomock.NewController(t)

	conexionMock := mocks.NewMockClient(ctrl)
	acceso := base.ObjetoAcceso{Cliente: conexionMock, Base: ""}

	modificacion := ldap.NewModifyRequest("", []ldap.Control{})
	errorPrueba := errors.New("Otro error improbable")

	conexionMock.
		EXPECT().
		Modify(modificacion).
		Return(errorPrueba).
		Times(1)

	// Empieza el test propiamente dicho
	err := acceso.ModificarUsuario("", []base.Reemplazo{})
	if err == nil {
		t.Errorf("Debe saber devolver un error")
	}
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
	esperado := crearEsperado(datosEntrada)
	acceso.ListarUsuarios("", []string{}, operacion)
	if !cmp.Equal(acceso.Datos, esperado) {
		t.Error(cmp.Diff(acceso.Datos, esperado))
	}
}

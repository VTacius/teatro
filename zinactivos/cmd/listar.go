package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
	"xibalba.com/vtacius/zinactivos/base"
	"xibalba.com/vtacius/zinactivos/utils"
)

// listarCmd represents the listar command
var listarCmd = &cobra.Command{
	Use:   "listar",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		listar()
	},
}

func init() {
	rootCmd.AddCommand(listarCmd)
}

func filtrarInactivosObjetos(respuesta *ldap.SearchResult) (resultado []base.Objeto) {
	hoy := time.Now()
	for _, item := range respuesta.Entries {
		DN := item.DN
		ZimbraLastLogonTimestamp := utils.Fechador(item.GetAttributeValue("zimbraLastLogonTimestamp"))
		if utils.RevisarIntervalo(15552000.0, hoy, ZimbraLastLogonTimestamp) == 1 {
			Username := strings.TrimSpace(item.GetAttributeValue("uid"))
			Nombre := strings.TrimSpace(item.GetAttributeValue("cn"))
			Description := strings.TrimSpace(item.GetAttributeValue("description"))
			resultado = append(resultado,
				base.Objeto{
					DN:                       DN,
					Username:                 Username,
					Nombre:                   Nombre,
					ZimbraLastLogonTimestamp: ZimbraLastLogonTimestamp,
					Description:              Description})
		}
	}

	return
}

func encontrarLongitud(entradas []base.Objeto) (resultUsername int, resultNombre int) {
	for _, item := range entradas {
		longUsername := len(item.Username)
		longNombre := len(item.Nombre)

		if longUsername > resultUsername {
			resultUsername = longUsername
		}

		if longNombre > resultNombre {
			resultNombre = longNombre
		}
	}
	return
}

func listar() {

	dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	conexion, err := utils.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}
	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}

	filtro := "(zimbraLastLogonTimestamp=*)"
	atributos := []string{"cn", "uid", "zimbraLastLogonTimestamp", "description"}
	err = acceso.ListarUsuarios(filtro, atributos, filtrarInactivosObjetos)
	if err != nil {
		utils.Salida("Error al ejecutar búsqueda", err)
	}

	longitudUsername, longitudNombre := encontrarLongitud(acceso.Datos)

	for _, usuario := range acceso.Datos {
		fmt.Printf("| %-*s | %-*s | %s | %-20s\n", longitudUsername, usuario.Username, longitudNombre, usuario.Nombre, usuario.ZimbraLastLogonTimestamp, usuario.Description)
	}
}

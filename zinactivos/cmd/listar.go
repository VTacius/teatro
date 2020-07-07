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

func filtrarInactivos(respuesta *ldap.SearchResult) (resultado []base.Usuario) {
	hoy := time.Now()

	for _, item := range respuesta.Entries {
		ZimbraLastLogonTimestamp := utils.Fechador(item.GetAttributeValue("zimbraLastLogonTimestamp"))

		if utils.RevisarIntervalo(15552000.0, hoy, ZimbraLastLogonTimestamp) {

			Username := strings.TrimSpace(item.GetAttributeValue("uid"))
			Nombre := strings.TrimSpace(item.GetAttributeValue("cn"))
			Description := strings.TrimSpace(item.GetAttributeValue("description"))

			resultado = append(resultado, base.Usuario{Username, Nombre, ZimbraLastLogonTimestamp, Description})

		}

	}

	return
}

func listar() {

	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	acceso := base.NewAccesoUsuario(url, usuario, contrasenia)

	filtro := "(zimbraLastLogonTimestamp=*)"
	atributos := []string{"cn", "uid", "zimbraLastLogonTimestamp", "description"}
	acceso.ListarUsuarios(filtro, atributos, filtrarInactivos)

	longitud := 50

	for _, usuario := range acceso.Datos {
		fecha := fmt.Sprintf("%02d/%02d/%d",
			usuario.ZimbraLastLogonTimestamp.Day(),
			usuario.ZimbraLastLogonTimestamp.Month(),
			usuario.ZimbraLastLogonTimestamp.Year())
		fmt.Printf("|%-20s| %-*s | %s | %-20s\n", usuario.Username, longitud, usuario.Nombre, fecha, usuario.Description)
	}
}

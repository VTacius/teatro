package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/zinactivos/base"
	"xibalba.com/vtacius/zinactivos/utils"
)

// listarCmd represents the listar command
var listarCmd = &cobra.Command{
	Use:   "listar",
	Short: "Lista todos los usuarios inactivos por más de seis meses",
	Run: func(cmd *cobra.Command, args []string) {
		dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
		usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
		// Actualmente, 6 meses * 30 dias * 24 horas * 60 minutos * 60 segundos
		listar(os.Stdout, 15552000.0, dnBase, usuario, contrasenia, url)
	},
}

func init() {
	rootCmd.AddCommand(listarCmd)
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

func listar(salida io.Writer, periodo float64, dnBase string, usuario string, contrasenia string, url string) {
	conexion, err := base.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}

	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}
	filtro := "(zimbraLastLogonTimestamp=*)"
	atributos := []string{"cn", "uid", "zimbraLastLogonTimestamp", "description"}

	err = acceso.ObtenerUsuarios(filtro, atributos)
	if err != nil {
		utils.Salida("Error al ejecutar búsqueda", err)
	}

	usuarios := acceso.ListarBusqueda().FiltrarInactivos(periodo)
	longitudUsername, longitudNombre := encontrarLongitud(usuarios)

	for _, usuario := range usuarios {
		fmt.Fprintf(salida, "| %-*s | %-*s | %s | %-20s\n", longitudUsername, usuario.Username, longitudNombre, usuario.Nombre, usuario.ZimbraLastLogonTimestamp, usuario.Description)
	}
}

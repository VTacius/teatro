package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/zinactivos/base"
	"xibalba.com/vtacius/zinactivos/utils"
)

// bloquearCmd represents the bloquear command
var bloquearCmd = &cobra.Command{
	Use:   "bloquear",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		bloquear()
	},
}

func init() {
	rootCmd.AddCommand(bloquearCmd)
}

func bloquearUsuario(fecha time.Time, DNUsuario string) {
	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	acceso := base.NewAccesoUsuario(url, usuario, contrasenia)

	marcaBorrado := fmt.Sprintf("ENCOLABORRADO : %s", fecha.Format(time.RFC822))
	modificacion := []utils.Reemplazo{
		{Clave: "description", Valor: marcaBorrado},
		{Clave: "zimbraMailStatus", Valor: "disabled"},
		{Clave: "zimbraAccountStatus", Valor: "closed"},
	}

	acceso.ModificarUsuario(DNUsuario, modificacion)

}

func bloquear() {
	fecha := time.Now()

	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	acceso := base.NewAccesoUsuario(url, usuario, contrasenia)

	filtro := "(&(zimbraLastLogonTimestamp=*)(!(description=ENCOLABORRADO*)))"
	atributos := []string{"uid"}
	acceso.ListarUsuarios(filtro, atributos, filtrarInactivos)

	for _, usuario := range acceso.Datos {
		bloquearUsuario(fecha, usuario.DN)
	}
}

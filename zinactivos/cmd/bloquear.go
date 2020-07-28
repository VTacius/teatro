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
	Short: "Bloquea los usuarios inactivos por más de seis meses",
	Run: func(cmd *cobra.Command, args []string) {
		ejecutar, _ := cmd.Flags().GetBool("ejecutar")
		dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
		usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
		// Actualmente, 6 meses * 30 dias * 24 horas * 60 minutos * 60 segundos
		bloquear(ejecutar, 15552000.0, dnBase, usuario, contrasenia, url)
	},
}

func init() {
	bloquearCmd.Flags().BoolP("ejecutar", "e", false, "")
	rootCmd.AddCommand(bloquearCmd)
}

func bloquearUsuario(fecha time.Time, DNUsuario string) error {

	dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	conexion, err := base.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}
	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}

	marcaBorrado := fmt.Sprintf("ENCOLABORRADO : %s", fecha.Format(time.RFC822))
	modificacion := []base.Reemplazo{
		{Clave: "description", Valor: marcaBorrado},
		{Clave: "zimbraMailStatus", Valor: "disabled"},
		{Clave: "zimbraAccountStatus", Valor: "closed"},
	}

	err = acceso.ModificarUsuario(DNUsuario, modificacion)
	return err
}

func bloquear(ejecutar bool, periodo float64, dnBase string, usuario string, contrasenia string, url string) {
	conexion, err := base.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}

	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}
	filtro := "(&(zimbraLastLogonTimestamp=*)(!(description=ENCOLABORRADO*)))"
	atributos := []string{"uid", "zimbraLastLogonTimestamp"}

	err = acceso.ObtenerUsuarios(filtro, atributos)
	if err != nil {
		utils.Salida("Error al ejecutar búsqueda", err)
	}

	usuarios := acceso.ListarBusqueda().FiltrarInactivos(periodo)

	if ejecutar {
		mensaje := fmt.Sprintf("\n\033[1m%d\033[0m usuarios serán bloqueados\n", len(usuarios))
		fmt.Println(mensaje)
	} else {
		mensaje := fmt.Sprintf("\nLos siguientes \033[1m%d\033[0m usuarios podrían ser bloqueados. \nUse la opción \033[1m-e\033[0m para ejecutar\n ", len(usuarios))
		fmt.Println(mensaje)
	}

	fecha := time.Now()

	for _, usuario := range usuarios {
		if ejecutar {
			err := bloquearUsuario(fecha, usuario.DN)
			if err != nil {
				fmt.Println(err, usuario.DN)
			}
		} else {
			fmt.Println(usuario.DN)
		}
	}
}

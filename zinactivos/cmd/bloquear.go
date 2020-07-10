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

func bloquearUsuario(fecha time.Time, DNUsuario string) error {

	dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	conexion, err := utils.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}
	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}

	marcaBorrado := fmt.Sprintf("ENCOLABORRADO : %s", fecha.Format(time.RFC822))
	modificacion := []utils.Reemplazo{
		{Clave: "description", Valor: marcaBorrado},
		{Clave: "zimbraMailStatus", Valor: "disabled"},
		{Clave: "zimbraAccountStatus", Valor: "closed"},
	}

	err = acceso.ModificarUsuario(DNUsuario, modificacion)
	return err
}

func bloquear() {

	dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	conexion, err := utils.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}
	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}
	filtro := "(&(zimbraLastLogonTimestamp=*)(!(description=ENCOLABORRADO*)))"
	atributos := []string{"uid", "zimbraLastLogonTimestamp"}
	acceso.ListarUsuarios(filtro, atributos, filtrarInactivosObjetos)

	for _, usuario := range acceso.Datos {
		fmt.Println(usuario.DN)
		//fecha := time.Now()
		//err := bloquearUsuario(fecha, usuario.DN)
		//if err != nil {
		//	fmt.Println(usuario.DN)
		//}
	}
}

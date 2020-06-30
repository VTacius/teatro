package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/ddesbloqueo/utils"
)

// desbloquearCmd represents the desbloquear command
var desbloquearCmd = &cobra.Command{
	Use:   "desbloquear",
	Short: "Desbloque un usuario",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		usuario := args[0]
		desbloquearUsuarios(usuario)
	},
}

func init() {
	rootCmd.AddCommand(desbloquearCmd)
}

func desbloquearUsuarios(usuario string) {
	cadenaConexion := utils.ObtenerCadenaConexion()

	conexion, err := utils.Conectar(cadenaConexion)

	if err != nil {
		utils.Salida("No se pudo conectar al servidor", err)
	}

	defer conexion.Close(context.Background())

	cadenaConsulta := `DELETE 
                        FROM intentos 
                        WHERE usuario=$1`

	resultado, err := conexion.Exec(context.Background(), cadenaConsulta, usuario)

	if err != nil {
		utils.Salida("No se pudo realizar tal operación", err)
	}

	usuario = fmt.Sprintf("\033[1m%s\033[0m", usuario)
	fueDesbloqueado := resultado.RowsAffected() >= 4

	var veredicto string
	if fueDesbloqueado {
		veredicto = fmt.Sprintf("El usuario %s ha sido desbloqueado", usuario)
	} else {
		veredicto = fmt.Sprintf("El usuario %s no estaba bloqueado", usuario)
	}

	fmt.Println(veredicto)
}

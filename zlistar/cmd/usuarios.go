package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/zlistar/base"
	"xibalba.com/vtacius/zlistar/utils"
)

// usuariosCmd represents the usuarios command
var usuariosCmd = &cobra.Command{
	Use:   "usuarios",
	Short: "Lista los usuarios",
	Run: func(cmd *cobra.Command, args []string) {
		salida := os.Stdout
		filtro := "(ObjectClass=zimbraAccount)"

		paraCSV, _ := cmd.Flags().GetBool("csv")
		dominio, _ := cmd.Flags().GetString("dominio")
		atributos, _ := cmd.Flags().GetStringArray("atributos")

		baseDN := utils.ConstruirBase(dominio)
		url, usuario, contrasenia := utils.ParametrosAccesoLdap()

		conexion, err := base.Conectar(url, usuario, contrasenia)
		if err != nil {
			utils.Ruptura("Error al conectarse", err)
		}

		usuarios := base.Acceso{Base: baseDN, Cliente: conexion}
		usuarios.Buscar(filtro, atributos).Listar()
		if usuarios.Err != nil {
			utils.Ruptura("Error al listar usuarios", usuarios.Err)
		}

		// Imprime el resultado en pantalla en el formato requerido
		if paraCSV {
			usuarios.ParaCSV(salida)
		} else {
			usuarios.Imprimir(salida)
		}
	},
}

func init() {
	rootCmd.AddCommand(usuariosCmd)
	usuariosCmd.Flags().Bool("csv", false, "Muestra el resultado como CSV")
	usuariosCmd.Flags().StringP("dominio", "d", "sv", "Dominio sobre el cual buscar")
	usuariosCmd.Flags().StringArrayP("atributos", "a", []string{"uid", "displayName"}, "Atributos a buscar")
}

package cmd

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
	"xibalba.com/vtacius/ddesbloqueo/utils"
)

// listarCmd represents the listar command
var listarCmd = &cobra.Command{
	Use:   "listar",
	Short: "Listar usuarios bloqueados",
	Run: func(cmd *cobra.Command, args []string) {
		usuario, _ := cmd.Flags().GetString("usuario")
		listarUsuarios(usuario)
	},
}

func init() {
	listarCmd.Flags().String("usuario", "ALL", "Especificar un usuario a verificar")
	rootCmd.AddCommand(listarCmd)
}

func obtenerAncho(conexion pgx.Conn) (int, error) {
	cadenaConsulta := `SELECT MAX(LENGTH(usuario))
                        FROM (
							SELECT usuario, count(usuario) AS intentos 
							FROM intentos 
							WHERE usuario !~ '\\s' 
							GROUP BY usuario) AS intentos 
                        WHERE intentos.intentos >= 4`
	contenido, err := conexion.Query(context.Background(), cadenaConsulta)
	if err != nil {
		return 0, err
	}

	var resultado int
	contenido.Next()
	err = contenido.Scan(&resultado)

	return resultado, err
}

func operacionBase(usuario string, conexion pgx.Conn) (pgx.Rows, error) {

	if usuario != "ALL" {
		cadenaConsulta := "SELECT usuario, count(usuario) FROM intentos WHERE usuario LIKE $1 GROUP BY usuario"
		contenido, err := conexion.Query(context.Background(), cadenaConsulta, usuario)

		return contenido, err
	}
	cadenaConsulta := `SELECT usuario, intentos
                        FROM (
							SELECT usuario, count(usuario) AS intentos 
							FROM intentos 
							WHERE usuario !~ '\\s' 
							GROUP BY usuario) AS intentos 
                        WHERE intentos.intentos >= 4`
	contenido, err := conexion.Query(context.Background(), cadenaConsulta)
	return contenido, err
}

func listarUsuarios(usuarioBusqueda string) {
	cadenaConexion := utils.ObtenerCadenaConexion()
	conexion, err := utils.Conectar(cadenaConexion)
	if err != nil {
		utils.Salida("No se pudo conectar al servidor", err)
	}

	ancho, err := obtenerAncho(conexion)
	if err != nil {
		utils.Salida("No se pudo conectar al servidor", err)
	}

	conexion, err = utils.Conectar(cadenaConexion)
	if err != nil {
		utils.Salida("No se pudo conectar al servidor", err)
	}
	contenido, err := operacionBase(usuarioBusqueda, conexion)
	if err != nil {
		utils.Salida("No se pudo conectar al servidor", err)
	}

	defer conexion.Close(context.Background())

	var usuario string
	var intentos int
	for contenido.Next() {
		err = contenido.Scan(&usuario, &intentos)

		if err != nil {
			utils.Salida("Consulta fallida", err)
		}

		fmt.Printf("| %-*s | %5d |\n", ancho, usuario, intentos)
	}

}

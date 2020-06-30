package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

// Salida Se asegura de salir con dignidad
func Salida(mensaje string, err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", mensaje, err)

	os.Exit(1)
}

// Conectar : Realiza la conexión al servidor de base de datos
func Conectar(cadenaConexion string) (pgx.Conn, error) {

	conexion, err := pgx.Connect(context.Background(), cadenaConexion)

	if err != nil {
		return *conexion, err
	}

	return *conexion, nil

}

package utils

import (
	"strings"
	"os"
	"fmt"
)

// ConstruirBase : Contruye una base DSN a partir de un dominio
func ConstruirBase(dominio string) string {
	componentes := strings.Split(dominio, ".")
	var resultado strings.Builder
	
	for _, dc := range componentes {
		unidad := "dc=" + dc + ","
		resultado.WriteString(unidad)
	}
	
	return strings.Trim(resultado.String(), ",") 
}

// Ruptura Se asegura de salir con dignidad
func Ruptura(mensaje string, err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", mensaje, err)

	os.Exit(1)
}
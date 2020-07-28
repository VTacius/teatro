package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	tm "github.com/buger/goterm"
)

type IntentoAcceso struct {
	usuario   string
	hostname  string
	ipaddress string
}

type Resultado struct {
	usuario   string
	hostname  string
	ipaddress string
	intento   int
}

func (resultado *Resultado) Mostrar(longitudes Longitudes) string {
	var cadena strings.Builder
	cadena.WriteString(fmt.Sprintf("%*d ", 8, resultado.intento))
	cadena.WriteString(fmt.Sprintf("%-*s ", longitudes.usuario, resultado.usuario))
	cadena.WriteString(fmt.Sprintf("%-*s ", longitudes.hostname, resultado.hostname))
	cadena.WriteString(fmt.Sprintf("%15s ", resultado.ipaddress))
	cadena.WriteString("\n")
	return cadena.String()

}

type Longitudes struct {
	usuario  int
	hostname int
}

func seleccionarCargaUtil(raw string) []string {
	indice := strings.Index(raw, "client=")
	indice += 7
	return strings.Fields(raw[indice:])
}

func extraerUsuario(raw string) string {
	indice := strings.Index(raw, "=")
	indice++
	return raw[indice:]
}
func parsearDireccion(raw string) (string, string) {
	inicio := strings.Index(raw, "[")
	direccion := raw[:inicio]
	inicio++
	//final := strings.Index(raw[inicio:], "]")
	ipaddress := strings.Trim(raw[inicio:], "],")
	return direccion, ipaddress
}
func mostrar(salida io.Writer) {
	fichero, err := os.Open("/var/log/zimbra.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", "Error leyendo fichero", err)
	}
	defer fichero.Close()

	scanner := bufio.NewScanner(fichero)
	intentos := make(map[IntentoAcceso]int)
	for scanner.Scan() {
		texto := scanner.Text()
		if indice := strings.Index(texto, "sasl_username"); indice > 0 {
			carga := seleccionarCargaUtil(texto)

			usuario := extraerUsuario(carga[2])
			hostname, ipaddress := parsearDireccion(carga[0])
			intentos[IntentoAcceso{usuario, hostname, ipaddress}]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", "Error leyendo fichero", err)
	}

	var resultados []Resultado
	var longitudes Longitudes
	for acceso, intentos := range intentos {
		item := Resultado{acceso.usuario, acceso.hostname, acceso.ipaddress, intentos}
		if len(item.usuario) > longitudes.usuario {
			longitudes.usuario = len(item.usuario)
		}
		if len(item.hostname) > longitudes.hostname {
			longitudes.hostname = len(item.hostname)
		}
		resultados = append(resultados, item)
	}
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].intento > resultados[j].intento
	})

	for _, intentos := range resultados {
		fmt.Fprintf(salida, intentos.Mostrar(longitudes))
	}

}

func main() {
	tm.Clear()
	for {

		tm.MoveCursor(1, 1)
		// Create Box with 30% width of current screen, and height of 20 lines
		box := tm.NewBox(tm.Width(), tm.Height()-1, 0)

		mostrar(box)

		// Mostramos todo
		tm.Print(box.String())

		tm.Flush()
		time.Sleep(5 * time.Second)
	}
}

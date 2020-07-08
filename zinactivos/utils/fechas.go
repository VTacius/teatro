package utils

import (
	"strconv"
	"time"
)

// RevisarIntervalo : ¿Las fechas tiene un intervalo mayor al requerido en segundos?
func RevisarIntervalo(intervalo float64, hoy time.Time, fecha time.Time) uint8 {
	lapso := hoy.Sub(fecha)
	if lapso < 0 {
		lapso *= -1
	}

	if lapso.Seconds() >= intervalo {
		return 1
	}

	return 0
}

func obtieneElemento(cadena string, inicio int, final int) int {
	resultado := 0
	if len(cadena) >= final {
		resultado, err := strconv.Atoi(cadena[inicio:final])
		if err != nil {
			resultado = 0
		}
		return resultado
	}
	return resultado

}

// Fechador Crea una fecha a partir de una cadena de texto
func Fechador(timestamp string) time.Time {
	anio := obtieneElemento(timestamp, 0, 4)
	mes := obtieneElemento(timestamp, 4, 6)
	dia := obtieneElemento(timestamp, 6, 8)
	hor := obtieneElemento(timestamp, 8, 10)
	min := obtieneElemento(timestamp, 10, 12)
	seg := obtieneElemento(timestamp, 12, 14)

	fecha := time.Date(anio, time.Month(mes), dia, hor, min, seg, 0, time.UTC)

	return fecha
}

package utils_test

import (
	"testing"
	"time"

	"xibalba.com/vtacius/zinactivos/utils"
)

func TestObtieneElemento(t *testing.T) {
	for _, caso := range parsearFechas {
		resultado := utils.ObtieneElemento(caso.timestamp, caso.inicio, caso.final)
		if resultado != caso.esperado {
			t.Errorf("Error en %s - Resultado: %d - Esperado: %d", caso.descripcion, resultado, caso.esperado)
		}
	}
}
func TestRevisarIntervalo(t *testing.T) {
	for _, caso := range intervalosDeFecha {
		resultado := utils.RevisarIntervalo(caso.intervalo, caso.hoy, caso.fecha)
		if resultado != caso.esperado {
			t.Errorf("Error en %s - Resultado: %d - Esperado %d", caso.descripcion, resultado, caso.esperado)
		}
	}
}

var parsearFechas = []struct {
	descripcion string
	timestamp   string
	inicio      int
	final       int
	esperado    int
}{
	{
		"Caso base",
		"2002",
		0,
		4,
		2002,
	},
	{
		"Obtención del mes",
		"200205",
		4,
		6,
		5,
	},
	{
		"Valor por defecto cuando la cadena sea más corta",
		"200205",
		6,
		8,
		0,
	},
	{
		"Valor por defecto cuando el valor no sea procesable como entero",
		"200205AB",
		6,
		8,
		0,
	},
}

var intervalosDeFecha = []struct {
	descripcion string
	intervalo   float64
	hoy         time.Time
	fecha       time.Time
	esperado    uint8
}{
	{
		"Fecha evaluada 8 minutos después de creada",
		86400,
		time.Date(2020, time.Month(10), 5, 7, 32, 0, 0, time.UTC),
		time.Date(2020, time.Month(10), 5, 7, 40, 0, 0, time.UTC),
		0,
	},
	{
		"Fecha evaluada 2 horas, 8 minutos después de creada",
		86400,
		time.Date(2020, time.Month(10), 5, 7, 32, 0, 0, time.UTC),
		time.Date(2020, time.Month(10), 5, 9, 40, 0, 0, time.UTC),
		0,
	},
	{
		"Fecha evaluada 23 horas, 59 minutos después de creada",
		86400,
		time.Date(2020, time.Month(10), 5, 7, 32, 0, 0, time.UTC),
		time.Date(2020, time.Month(10), 5, 8, 31, 0, 0, time.UTC),
		0,
	},
	{
		"Fecha evaluada 2 dias, 8 minutos después de creada",
		86400,
		time.Date(2020, time.Month(10), 5, 7, 32, 0, 0, time.UTC),
		time.Date(2020, time.Month(10), 7, 7, 40, 0, 0, time.UTC),
		1,
	},
}

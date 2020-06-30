BASE=$(PWD)
BIN=$(BASE)/bin

PCCD:=$(file < private/cadena_conexion_directorio) 

.PHONY: prueba ddesbloqueo ddesbloqueo_dep

all: prueba 

prueba:
	go build -o ${BIN}/prueba -ldflags '-X xibalba.com/teatro/prueba/main.version=${PCCD}' xibalba.com/teatro/prueba

ddesbloqueo_dep:
	go get ./ddesbloqueo

ddesbloqueo: ddesbloqueo_dep
	go build -o ./bin/ -ldflags "-X xibalba.com/vtacius/ddesbloqueo/utils.Conexion=$(shell cat private/cadena_conexion_directorio)" ./ddesbloqueo/main.go
	#go build -o ${BIN}/ddesbloqueo -ldflags "-X xibalba.com/teatro/ddesbloqueo/utils.Conexion=${PCCD}" xibalba.com/teatro/ddesbloqueo

BASE=$(PWD)
BIN=$(BASE)/bin

PCCD:=$(file < private/cadena_conexion_directorio) 

.PHONY: dep build 

all: build 

dep:
	go get 

build: dep 
	go build -o ./bin/ -ldflags "-X xibalba.com/vtacius/ddesbloqueo/utils.Conexion=$(shell cat private/cadena_conexion_directorio)" 

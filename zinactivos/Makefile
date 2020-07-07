BASE=$(PWD)
BIN=$(BASE)/bin

PCCD:=$(file < private/cadena_conexion_directorio) 

.PHONY: dep build 

all: build 

dep:
	go get 

build: dep 
	go build -o ./bin/ -ldflags "-X xibalba.com/vtacius/zinactivos/utils.Credenciales=$(shell cat private/credenciales_conexion_zimbra)" 

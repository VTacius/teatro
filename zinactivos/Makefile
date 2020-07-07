BASE=$(PWD)
BIN=$(BASE)/bin

ZPL:=$(file < private/zimbra_password_ldap) 
ZUL:=$(file < private/zimbra_usuario_ldap) 
ZSL:=$(file < private/zimbra_servidor_ldap) 
ZPC:=$(file < private/zimbra_password_correo) 
ZUC:=$(file < private/zimbra_usuario_correo) 
ZSC:=$(file < private/zimbra_servidor_correo) 

LDFLAGS:=-ldflags 
LDFLAGS+="-X xibalba.com/vtacius/zinactivos/utils.PasswordLdap=${ZPL}
LDFLAGS+=-X xibalba.com/vtacius/zinactivos/utils.UsuarioLdap=${ZUL}
LDFLAGS+=-X xibalba.com/vtacius/zinactivos/utils.ServidorLdap=${ZSL}
LDFLAGS+=-X xibalba.com/vtacius/zinactivos/utils.PasswordCorreo=${ZPC}
LDFLAGS+=-X xibalba.com/vtacius/zinactivos/utils.UsuarioCorreo=${ZUC}
LDFLAGS+=-X xibalba.com/vtacius/zinactivos/utils.ServidorCorreo=${ZSC}"
.PHONY: dep build 

all: build 

dep:
	go get 

build: dep 
	go build -o ./bin/ ${LDFLAGS} 

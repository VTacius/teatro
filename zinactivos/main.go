package main
import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"xibalba.com/vtacius/zinactivos/utils"
)

func revisarIntervalo(intervalo float64, hoy time.Time, fecha time.Time) bool {
	lapso := hoy.Sub(fecha)
	return lapso.Seconds() >= intervalo 
}

func filtrarInactivos(respuesta *ldap.SearchResult) (resultado []Usuario, longitud int) {
	hoy := time.Now()
	
	for _, item := range(respuesta.Entries) {
		zimbraLastLogonTimestamp := fechador(item.GetAttributeValue("zimbraLastLogonTimestamp"))
		nombre := strings.TrimSpace(item.GetAttributeValue("cn"))
		
		if revisarIntervalo(15552000.0, hoy, zimbraLastLogonTimestamp){
			resultado = append(resultado, Usuario{nombre, zimbraLastLogonTimestamp})	
			largo := len([]rune(nombre))
			if largo > longitud {
				longitud = largo
			}	
		}
	
	}

	return
}

func obtieneElemento(cadena string, inicio int, final int) int {
	resultado := 0
	if len(cadena) >= final {
		resultado, err := strconv.Atoi(cadena[inicio:final])
		if err != nil{
			resultado = 0
		}	
		return resultado
	}
	return resultado 
	
}

func fechador(timestamp string) time.Time {
	anio := obtieneElemento(timestamp, 0, 4)
	mes := obtieneElemento(timestamp, 4, 6)
	dia := obtieneElemento(timestamp, 6, 8)
	hor := obtieneElemento(timestamp, 8, 10)
	min := obtieneElemento(timestamp, 10, 12)
	seg := obtieneElemento(timestamp, 12, 14)
	
	fecha := time.Date(anio, time.Month(mes), dia, hor, min, seg, 0, time.UTC)
	
	return fecha
}

func conectar(url string, usuario string, contrasenia string) (*ldap.Conn, error) {
	// Conexión inicial
	config := ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}) 
	conexion, err := ldap.DialURL(url, config)
	if err != nil {
		return conexion, err
	}
	
	// Nos autenticamos
	err = conexion.Bind(usuario, contrasenia)
	if err != nil {
		return conexion, err
	}

	return conexion, nil
}

func peticionar(conexion ldap.Conn, base string, filtro string, atributos []string) (*ldap.SearchResult, error) {

	peticion := ldap.NewSearchRequest(base, 
		ldap.ScopeWholeSubtree, 
		ldap.NeverDerefAliases, 
		0, 0, false, filtro, atributos, nil)

	respuesta, err := conexion.Search(peticion)
	if err != nil {
		return respuesta, err
	}

	return respuesta, nil
}

// Usuario Contiene la definición actual de usuario
type Usuario struct {
	nombre string;
	zimbraLastLogonTimestamp time.Time
}

func main() {
	url := "ldap://10.10.20.2:389"	
	usuario, contrasenia := utils.ObtenerCredenciales()
	
	conexion, err := conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error en la conexión", err)
	}
	
	base := "ou=people,dc=salud,dc=gob,dc=sv"
	filtro := "(zimbraLastLogonTimestamp=*)"
	atributos :=[]string{"cn", "uid", "zimbraLastLogonTimestamp"}
	
	respuesta, err := peticionar(*conexion, base, filtro, atributos)
	if err != nil {
		utils.Salida("Error en la peticion", err)
	}
	
	usuariosInactivos, longitud := filtrarInactivos(respuesta)
	for _, usuario := range(usuariosInactivos){
		fecha := fmt.Sprintf("%02d/%02d/%d", 
								usuario.zimbraLastLogonTimestamp.Day(), 
								usuario.zimbraLastLogonTimestamp.Month(),
								usuario.zimbraLastLogonTimestamp.Year())
		fmt.Printf("| %-*s | %s |\n", longitud, usuario.nombre, fecha)
	}
}
package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/zinactivos/base"
	"xibalba.com/vtacius/zinactivos/utils"
)

// reporteCmd represents the reporte command
var reporteCmd = &cobra.Command{
	Use:   "reporte",
	Short: "Envia por correo un reporte del trabajo realizado sobre los usuarios bloqueados",
	Run: func(cmd *cobra.Command, args []string) {
		dnBase := "ou=people,dc=salud,dc=gob,dc=sv"
		usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()

		reporte(dnBase, usuario, contrasenia, url)
	},
}

func init() {
	rootCmd.AddCommand(reporteCmd)
}

var mensaje string = `
<html>

<head>
  	<meta charset="UTF-8">
  	<style>
  		tr:nth-child(even) {
			background-color: #f2f2f2;
  		}
  	</style>
</head>

{{if gt (len .Proximos) 0}}
<body>
	<table style="border-collapse: collapse">
		<tr>
			  <th colspan="4"><h3>Usuarios proximos a ser borrados<h3></th>
		<tr>
		<tr>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Usuario</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Nombre</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Última Conexión</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Estado</th>
		</tr>
		{{range .Proximos}}
		<tr>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Username}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Nombre}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{fecha .ZimbraLastLogonTimestamp}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Description}}</td>
		</tr>
		{{end}}
	</table>
</body>
{{end}}
{{if gt (len .Bloqueados) 0}}
<body>
	<table style="border-collapse: collapse">
		<tr>
			  <th colspan="4"><h3>Usuarios bloqueados hoy<h3></th>
		<tr>
		<tr>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Usuario</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Nombre</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Última Conexión</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Estado</th>
		</tr>
		{{range .Bloqueados}}
		<tr>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Username}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Nombre}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{fecha .ZimbraLastLogonTimestamp}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Description}}</td>
		</tr>
		{{end}}
	</table>
</body>
{{end}}
`

type datos struct {
	Bloqueados []base.Objeto
	Proximos   []base.Objeto
}

func contenido(subject string, data datos) (string, error) {
	var salida bytes.Buffer

	template, err := template.New("reporte").Funcs(template.FuncMap{
		"fecha": func(fecha time.Time) string {
			return fecha.Format(time.RFC822)
		},
	}).Parse(mensaje)
	if err != nil {
		return "", err
	}
	err = template.Execute(&salida, data)
	if err != nil {
		return "", err
	}

	return salida.String(), nil

}

func bloqueadoHoy(hoy time.Time, fecha time.Time) (resultado bool) {
	// No ha pasado más de un día de segundos
	return utils.RevisarIntervalo(86400, hoy, fecha) == 0
}

func proximoBorrado(hoy time.Time, fecha time.Time) (resultado bool) {
	// Han pasado más de un par de dos meses y 25 días
	return utils.RevisarIntervalo(7344000, hoy, fecha) == 1
}

func clasificarEntradas(entradas []base.Objeto) (bloqueadosHoy []base.Objeto, proximosABorrar []base.Objeto) {
	hoy := time.Now()
	for _, item := range entradas {
		timestamp := strings.SplitN(item.Description, ":", 2)[1]
		fecha, error := time.Parse(time.RFC822, strings.TrimSpace(timestamp))
		if error != nil {
			fmt.Printf("Fecha invalida %s\n", error)
		}
		if bloqueadoHoy(hoy, fecha) {
			bloqueadosHoy = append(bloqueadosHoy, item)
		} else if proximoBorrado(hoy, fecha) {
			proximosABorrar = append(proximosABorrar, item)
		}
	}
	return
}

func reporte(dnBase string, usuario string, contrasenia string, url string) {

	conexion, err := base.Conectar(url, usuario, contrasenia)
	if err != nil {
		utils.Salida("Error al conectarse", err)
	}

	acceso := base.ObjetoAcceso{Cliente: conexion, Base: dnBase}
	filtro := "(&(zimbraLastLogonTimestamp=*)(description=ENCOLABORRADO*))"
	atributos := []string{"cn", "uid", "zimbraLastLogonTimestamp", "description"}

	err = acceso.ObtenerUsuarios(filtro, atributos)
	if err != nil {
		utils.Salida("Error al ejecutar búsqueda", err)
	}

	acceso.ListarBusqueda()
	bloqueados, proximos := clasificarEntradas(acceso.Datos)

	// Si no hay datos que mostrar, pues que no se envia nada

	if len(bloqueados) == 0 && len(proximos) == 0 {
		fmt.Println("No envia reporte porque no hay nada que reportar")
		os.Exit(0)
	}

	// Empiezo con el envío de correo
	usuarioCorreo, passwordCorreo, servidorCorreo := utils.ConfiguracionEnvioCorreo()
	subject := "Subject: Reporte de usuarios inactivos\n"
	rcpt := "To: <alortiz@salud.gob.sv>, <rnajarro@salud.gob.sv>, <ajgonzalez@salud.gob.sv>\n"

	data := datos{Bloqueados: bloqueados, Proximos: proximos}

	mensaje, err := contenido(subject, data)
	if err != nil {
		utils.Salida("Error crear cuerpo del correo", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + rcpt + mime + mensaje)

	auth := smtp.PlainAuth("", usuarioCorreo, passwordCorreo, servidorCorreo)

	err = smtp.SendMail(servidorCorreo+":587", auth, usuarioCorreo, []string{"alortiz@salud.gob.sv"}, msg)
	if err != nil {
		utils.Salida("Error al enviar correo", err)
	}
}

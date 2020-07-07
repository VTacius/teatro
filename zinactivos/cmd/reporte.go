package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/zinactivos/base"
	"xibalba.com/vtacius/zinactivos/utils"
)

// reporteCmd represents the reporte command
var reporteCmd = &cobra.Command{
	Use:   "reporte",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reporte()
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

<body>
	<table style="border-collapse: collapse">
		<tr>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Usuario</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Nombre</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Última Conexión</th>
			<th style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">Estado</th>
		</tr>
		{{range .Usuarios}}
		<tr>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Username}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Nombre}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.ZimbraLastLogonTimestamp}}</td>
			<td style="padding: 0.5rem; text-align: left; border-bottom: 1px solid #ddd;">{{.Description}}</td>
		</tr>
		{{end}}
	</table>
</body>
`

type datos struct {
	Usuarios []base.Usuario
}

func contenido(subject string, data datos) (string, error) {
	var salida bytes.Buffer

	template, err := template.New("reporte").Parse(mensaje)
	if err != nil {
		return "", err
	}

	err = template.Execute(&salida, data)
	if err != nil {
		return "", err
	}

	return salida.String(), nil

}

func reporte() {

	usuario, contrasenia, url := utils.ConfiguracionAccesoLdap()
	acceso := base.NewAccesoUsuario(url, usuario, contrasenia)

	filtro := "(zimbraLastLogonTimestamp=*)"
	atributos := []string{"cn", "uid", "zimbraLastLogonTimestamp", "description"}
	acceso.ListarUsuarios(filtro, atributos, filtrarInactivos)

	// Empiezo con el envío de correo
	usuarioCorreo, passwordCorreo, servidorCorreo := utils.ConfiguracionEnvioCorreo()
	subject := "Subject: Reporte de usuarios inactivos\n"
	rcpt := "To: <alortiz@salud.gob.sv>\n"

	data := datos{Usuarios: acceso.Datos}

	mensaje, err := contenido(subject, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + rcpt + mime + mensaje)

	auth := smtp.PlainAuth("", usuarioCorreo, passwordCorreo, servidorCorreo)

	err = smtp.SendMail(servidorCorreo+":587", auth, usuarioCorreo, []string{"alortiz@salud.gob.sv"}, msg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

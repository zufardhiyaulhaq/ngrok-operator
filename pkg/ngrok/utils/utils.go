package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/models"
)

const NGROK_URL_REGEX = `https://.[^"]+`
const TMPL = `web_addr: 0.0.0.0:4040
{{if .Spec.AuthToken }}
authtoken: {{ .Spec.AuthToken }}
{{end}}
{{if .Spec.Region }}
region: {{ .Spec.Region }}
{{end}}
tunnels:
  app:
    proto: {{ .Spec.Protocol }}
    addr: {{ .Spec.Service }}:{{ .Spec.Port }}
    {{if eq .Spec.Protocol "http"}}
    inspect: {{ .Spec.Inspect }}
    {{if .Spec.Auth }}
    auth: {{ .Spec.Auth }}
    {{ end }}
    {{if .Spec.AuthToken }}{{if .Spec.Hostname }}
    hostname: {{ .Spec.Hostname }}
    {{end}}{{end}}
    {{end}}
    {{if eq .Spec.Protocol "tcp"}}
    {{if .Spec.AuthToken }}{{if .Spec.Hostname }}
    remote_addr: {{ .Spec.Hostname }}
    {{end}}{{end}}
    {{end}}
`

func GetNgrokURL(adminAPI string) (string, error) {
	response, err := http.Get(adminAPI)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var configuration models.TunnelsConfiguration
	err = json.Unmarshal(body, &configuration)
	if err != nil {
		return "", err
	}

	return configuration.Tunnels[0].PublicURL, nil
}

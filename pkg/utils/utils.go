package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/models"
)

const NGROK_URL_REGEX = `https://.[^"]+`
const TMPL = `web_addr: 0.0.0.0:4040
{{if .AuthToken }}
authtoken: {{ .AuthToken }}
{{end}}
{{if .Region }}
region: {{ .Region }}
{{end}}
tunnels:
  app:
    proto: {{ .Protocol }}
    addr: {{ .Service }}:{{ .Port }}
    {{if eq .Protocol "http"}}
    inspect: {{ .Inspect }}
    {{if .Auth }}
    auth: {{ .Auth }}
    {{ end }}
    {{if .AuthToken }}{{if .Hostname }}
    hostname: {{ .Hostname }}
    {{end}}{{end}}
    {{end}}
    {{if eq .Protocol "tcp"}}
    {{if .AuthToken }}{{if .RemoteAddr }}
    remote_addr: {{ .RemoteAddr }}
    {{end}}{{end}}
    {{end}}
`

func GetNgrokURL(api string) (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(api)
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

	if len(configuration.Tunnels) == 0 {
		return "", fmt.Errorf("configuration tunnels empty")
	}

	return configuration.Tunnels[0].PublicURL, nil
}

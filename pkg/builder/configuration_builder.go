package builder

import (
	"bytes"
	"strings"
	"text/template"

	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/utils"
)

type NgrokConfigurationBuilder struct {
	Spec *ngrokcomv1alpha1.NgrokSpec
}

func NewNgrokConfigurationBuilder() *NgrokConfigurationBuilder {
	return &NgrokConfigurationBuilder{}
}

func (n *NgrokConfigurationBuilder) SetSpec(spec *ngrokcomv1alpha1.NgrokSpec) *NgrokConfigurationBuilder {
	n.Spec = spec
	return n
}

func (n *NgrokConfigurationBuilder) Build() (string, error) {
	err := n.Spec.Validate()
	if err != nil {
		return "", err
	}

	var configurationBuffer bytes.Buffer

	templateEngine, err := template.New("ngrok").Parse(utils.TMPL)
	if err != nil {
		return "", err
	}

	err = templateEngine.Execute(&configurationBuffer, n.Spec)
	if err != nil {
		return "", err
	}

	var configuration []string
	for _, data := range strings.Split(configurationBuffer.String(), "\n") {
		if len(strings.TrimSpace(data)) != 0 {
			configuration = append(configuration, data)
		}
	}

	return strings.Join(configuration, "\n"), nil
}

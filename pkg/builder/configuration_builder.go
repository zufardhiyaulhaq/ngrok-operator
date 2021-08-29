package builder

import (
	"bytes"
	"html/template"

	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NgrokConfigurationBuilder struct {
	Spec *ngrokcomv1alpha1.NgrokSpec
}

func NewNgrokConfigurationBuilder(client client.Client) *NgrokConfigurationBuilder {
	return &NgrokConfigurationBuilder{}
}

func (n *NgrokConfigurationBuilder) SetSpec(spec *ngrokcomv1alpha1.NgrokSpec) *NgrokConfigurationBuilder {
	n.Spec = spec
	return n
}

func (n *NgrokConfigurationBuilder) Build() (string, error) {
	var configuration bytes.Buffer

	templateEngine, err := template.New("ngrok").Parse(utils.TMPL)
	if err != nil {
		return "", err
	}

	err = templateEngine.Execute(&configuration, n.Spec)
	if err != nil {
		return "", err
	}

	return configuration.String(), nil
}

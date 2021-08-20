package builder

import (
	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NgrokConfigMapBuilder struct {
	config *ngrokcomv1alpha1.Ngrok
}

func NewNgrokConfigMapBuilder() *NgrokConfigMapBuilder {
	return &NgrokConfigMapBuilder{}
}

func (n *NgrokConfigMapBuilder) SetConfig(config *ngrokcomv1alpha1.Ngrok) *NgrokConfigMapBuilder {
	n.config = config
	return n
}

func (n *NgrokConfigMapBuilder) Build() (*corev1.ConfigMap, error) {
	data := make(map[string]string)

	config, err := n.config.GenerateConfiguration()
	if err != nil {
		return nil, err
	}

	data["ngrok.conf"] = config
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.config.Name + "-cm-ngrok",
			Namespace: n.config.Namespace,
			Labels: map[string]string{
				"app":       n.config.Name,
				"generated": "ngrok-operator",
			},
		},
		Data: data,
	}

	return configMap, nil
}

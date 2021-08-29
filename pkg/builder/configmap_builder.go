package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NgrokConfigMapBuilder struct {
	Name      string
	Namespace string
	Config    string
}

func NewNgrokConfigMapBuilder() *NgrokConfigMapBuilder {
	return &NgrokConfigMapBuilder{}
}

func (n *NgrokConfigMapBuilder) SetConfig(config string) *NgrokConfigMapBuilder {
	n.Config = config
	return n
}

func (n *NgrokConfigMapBuilder) SetName(name string) *NgrokConfigMapBuilder {
	n.Name = name
	return n
}

func (n *NgrokConfigMapBuilder) SetNamespace(namespace string) *NgrokConfigMapBuilder {
	n.Namespace = namespace
	return n
}

func (n *NgrokConfigMapBuilder) Build() (*corev1.ConfigMap, error) {
	data := make(map[string]string)

	data["ngrok.conf"] = n.Config
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.Name + "-cm-ngrok",
			Namespace: n.Namespace,
			Labels: map[string]string{
				"app":       n.Name,
				"generated": "ngrok-operator",
			},
		},
		Data: data,
	}

	return configMap, nil
}

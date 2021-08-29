package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type NgrokServiceBuilder struct {
	Name      string
	Namespace string
}

func NewNgrokServiceBuilder() *NgrokServiceBuilder {
	return &NgrokServiceBuilder{}
}

func (n *NgrokServiceBuilder) SetName(name string) *NgrokServiceBuilder {
	n.Name = name
	return n
}

func (n *NgrokServiceBuilder) SetNamespace(namespace string) *NgrokServiceBuilder {
	n.Namespace = namespace
	return n
}

func (n *NgrokServiceBuilder) Build() (*corev1.Service, error) {
	Service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.Name + "-ngrok",
			Namespace: n.Namespace,
			Labels: map[string]string{
				"app":       n.Name,
				"generated": "ngrok-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":       n.Name,
				"generated": "ngrok-operator",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http-api",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(4040),
				},
			},
			Type: "ClusterIP",
		},
	}

	return Service, nil
}

package builder

import (
	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type NgrokServiceBuilder struct {
	config *ngrokcomv1alpha1.Ngrok
}

func NewNgrokServiceBuilder() *NgrokServiceBuilder {
	return &NgrokServiceBuilder{}
}

func (n *NgrokServiceBuilder) SetConfig(config *ngrokcomv1alpha1.Ngrok) *NgrokServiceBuilder {
	n.config = config
	return n
}

func (n *NgrokServiceBuilder) Build() (*corev1.Service, error) {
	Service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.config.Name + "-ngrok",
			Namespace: n.config.Namespace,
			Labels: map[string]string{
				"app":       n.config.Name,
				"generated": "ngrok-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":       n.config.Name,
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

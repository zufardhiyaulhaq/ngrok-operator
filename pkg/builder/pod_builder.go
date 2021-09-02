package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NgrokPodBuilder struct {
	Name      string
	Namespace string
	Image     string
}

func NewNgrokPodBuilder() *NgrokPodBuilder {
	return &NgrokPodBuilder{}
}

func (n *NgrokPodBuilder) SetName(name string) *NgrokPodBuilder {
	n.Name = name
	return n
}

func (n *NgrokPodBuilder) SetNamespace(namespace string) *NgrokPodBuilder {
	n.Namespace = namespace
	return n
}

func (n *NgrokPodBuilder) SetImage(image string) *NgrokPodBuilder {
	n.Image = image
	return n
}

func (n *NgrokPodBuilder) Build() (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.Name + "-ngrok",
			Namespace: n.Namespace,
			Labels: map[string]string{
				"app":       n.Name,
				"generated": "ngrok-operator",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "ngrok",
					Image:   n.Image,
					Command: []string{"ngrok", "start", "--config", "/ngrok/ngrok.conf", "--all"},
					Ports: []corev1.ContainerPort{
						{ContainerPort: int32(4040)},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      n.Name + "-cm-ngrok",
							MountPath: "/ngrok",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: n.Name + "-cm-ngrok",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: n.Name + "-cm-ngrok",
							},
						},
					},
				},
			},
		},
	}

	return pod, nil
}

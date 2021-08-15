package builder

import (
	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NgrokPodBuilder struct {
	config *ngrokcomv1alpha1.Ngrok
}

func NewNgrokPodBuilder() *NgrokPodBuilder {
	return &NgrokPodBuilder{}
}

func (n *NgrokPodBuilder) SetConfig(config *ngrokcomv1alpha1.Ngrok) *NgrokPodBuilder {
	n.config = config
	return n
}

func (n *NgrokPodBuilder) Build() (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.config.Name + "-ngrok",
			Namespace: n.config.Namespace,
			Labels: map[string]string{
				"app": n.config.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "ngrok",
					Image:   "wernight/ngrok",
					Command: []string{"ngrok", "start", "--config", "/ngrok/ngrok.conf", "--all"},
					Ports: []corev1.ContainerPort{
						{ContainerPort: int32(4040)},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      n.config.Name + "-cm-ngrok",
							MountPath: "/ngrok",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: n.config.Name + "-cm-ngrok",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: n.config.Name + "-cm-ngrok",
							},
						},
					},
				},
			},
		},
	}

	return pod, nil
}

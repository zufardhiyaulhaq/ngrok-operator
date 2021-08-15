/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"
	builder "github.com/zufardhiyaulhaq/ngrok-operator/pkg/ngrok/builder"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/ngrok/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const NGROK_STATUS_CREATING = "Creating"
const NGROK_STATUS_CREATED = "Created"
const NGROK_STATUS_URL_FETCHING = "Fetching"

// NgrokReconciler reconciles a Ngrok object
type NgrokReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ngrok.com,resources=ngroks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ngrok.com,resources=ngroks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ngrok.com,resources=ngroks/finalizers,verbs=update

func (r *NgrokReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	ngrok := &ngrokcomv1alpha1.Ngrok{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, ngrok)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.SetStatus(ngrok, NGROK_STATUS_CREATING, NGROK_STATUS_URL_FETCHING)
	if err != nil {
		return ctrl.Result{}, err
	}

	configmap, err := builder.NewNgrokConfigMapBuilder().
		SetConfig(ngrok).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	pod, err := builder.NewNgrokPodBuilder().
		SetConfig(ngrok).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(ngrok, configmap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(ngrok, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	createdConfigMap := &corev1.ConfigMap{}
	createdPod := &corev1.Pod{}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, createdConfigMap)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), configmap)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, createdPod)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	time.Sleep(30 * time.Second)
	var ngrokURL string

	if createdPod.Status.PodIP != "" {
		adminAPI := "http://" + createdPod.Status.PodIP + ":4040/api/tunnels"
		ngrokURL, err = utils.GetNgrokURL(adminAPI)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		return ctrl.Result{Requeue: true}, nil
	}

	err = r.SetStatus(ngrok, NGROK_STATUS_CREATED, ngrokURL)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NgrokReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ngrokcomv1alpha1.Ngrok{}).
		Complete(r)
}

func (r *NgrokReconciler) SetStatus(ngrok *ngrokcomv1alpha1.Ngrok, status string, url string) error {
	ngrok.Status.Status = status
	ngrok.Status.URL = url

	return r.Client.Status().Update(context.TODO(), ngrok)
}

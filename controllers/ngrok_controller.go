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
	builder "github.com/zufardhiyaulhaq/ngrok-operator/pkg/builder"
	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/handler"
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
	Scheme        *runtime.Scheme
	StatusHandler handler.StatusHandler
}

//+kubebuilder:rbac:groups=ngrok.com,resources=ngroks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ngrok.com,resources=ngroks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ngrok.com,resources=ngroks/finalizers,verbs=update

func (r *NgrokReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Start Ngrok Reconciler")

	ngrok := &ngrokcomv1alpha1.Ngrok{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, ngrok)
	if err != nil {
		return ctrl.Result{}, nil
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

	service, err := builder.NewNgrokServiceBuilder().
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

	if err := controllerutil.SetControllerReference(ngrok, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	createdConfigMap := &corev1.ConfigMap{}
	createdPod := &corev1.Pod{}
	createdService := &corev1.Service{}

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

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, createdService)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), service)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	time.Sleep(60 * time.Second)

	if createdPod.Status.Phase != corev1.PodRunning {
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	url, err := utils.GetNgrokURL("http://" + service.Name + "." + service.Namespace + ".svc" + "/api/tunnels")
	if err != nil {
		return ctrl.Result{}, err
	}

	// get the status from ngrok URL
	// if it's not running, recreate the pod
	status, err := r.StatusHandler.Running(url)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !status {
		err := r.Client.Delete(context.TODO(), createdPod)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	ngrok.Status.Status = NGROK_STATUS_CREATED
	ngrok.Status.URL = url

	err = r.Client.Status().Update(context.TODO(), ngrok)
	if err != nil {
		return ctrl.Result{}, err
	}

	// rather than finished the process and reconcile when object changed
	// force to reconcile every 30 seconds
	return ctrl.Result{RequeueAfter: time.Second * 30}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NgrokReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ngrokcomv1alpha1.Ngrok{}).
		Complete(r)
}

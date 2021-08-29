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
	"fmt"
	"time"

	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"
	builder "github.com/zufardhiyaulhaq/ngrok-operator/pkg/builder"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/handler"
	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	AUTH_TOKEN_TYPE_PLAIN  = "plain"
	AUTH_TOKEN_TYPE_SECRET = "secret"

	NGROK_STATUS_CREATING     = "Creating"
	NGROK_STATUS_CREATED      = "Created"
	NGROK_STATUS_URL_FETCHING = "Fetching"
)

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

	name := ngrok.Name
	namespace := ngrok.Namespace
	spec := ngrok.Spec.DeepCopy()

	if spec.AuthTokenType == AUTH_TOKEN_TYPE_SECRET {
		secret := &corev1.Secret{}

		err := r.Client.Get(context.TODO(), types.NamespacedName{Name: spec.AuthToken, Namespace: namespace}, secret)
		if err != nil && errors.IsNotFound(err) {
			log.Error(err, "cannot find Secret")
			return ctrl.Result{}, err
		} else if err != nil {
			return ctrl.Result{}, err
		}

		base64AuthToken, ok := secret.Data["authToken"]
		if !ok {
			log.Error(fmt.Errorf("cannot find key"), "cannot find key authToken")
			return ctrl.Result{}, err
		}

		spec.AuthToken = string(base64AuthToken)
	}

	log.Info("Build configuration")
	configuration, err := builder.NewNgrokConfigurationBuilder(r.Client).
		SetSpec(spec).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Build config map")
	configmap, err := builder.NewNgrokConfigMapBuilder().
		SetConfig(configuration).
		SetName(name).
		SetNamespace(namespace).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Build pod")
	pod, err := builder.NewNgrokPodBuilder().
		SetName(name).
		SetNamespace(namespace).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Build service")
	service, err := builder.NewNgrokServiceBuilder().
		SetName(name).
		SetNamespace(namespace).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference config map")
	if err := controllerutil.SetControllerReference(ngrok, configmap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference pod")
	if err := controllerutil.SetControllerReference(ngrok, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference service")
	if err := controllerutil.SetControllerReference(ngrok, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	createdConfigMap := &corev1.ConfigMap{}
	createdPod := &corev1.Pod{}
	createdService := &corev1.Service{}

	log.Info("get config map")
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, createdConfigMap)
	if err != nil && errors.IsNotFound(err) {
		log.Info("create config map")
		err = r.Client.Create(context.TODO(), configmap)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("get pod")
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, createdPod)
	if err != nil && errors.IsNotFound(err) {
		log.Info("create pod")
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("get service")
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, createdService)
	if err != nil && errors.IsNotFound(err) {
		log.Info("create service")
		err = r.Client.Create(context.TODO(), service)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("check pod running")
	if createdPod.Status.Phase != corev1.PodRunning {
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	log.Info("get ngrok url")
	url, err := utils.GetNgrokURL("http://" + service.Name + "." + service.Namespace + ".svc" + "/api/tunnels")
	if err != nil {
		return ctrl.Result{}, err
	}

	// get the status from ngrok URL
	// if it's not running, recreate the pod
	log.Info("get status ngrok")
	status, err := r.StatusHandler.Running(url)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !status {
		log.Info("delete ngrok pod to restart session")
		err := r.Client.Delete(context.TODO(), createdPod)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	log.Info("get ngrok")
	ngrok = &ngrokcomv1alpha1.Ngrok{}
	err = r.Client.Get(context.TODO(), req.NamespacedName, ngrok)
	if err != nil {
		return ctrl.Result{}, nil
	}

	ngrok.Status.Status = NGROK_STATUS_CREATED
	ngrok.Status.URL = url

	log.Info("update ngrok status")
	err = r.Client.Status().Update(context.TODO(), ngrok)
	if err != nil {
		return ctrl.Result{}, err
	}

	// rather than finished the process and reconcile when object changed
	// force to reconcile every 60 seconds
	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NgrokReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ngrokcomv1alpha1.Ngrok{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}

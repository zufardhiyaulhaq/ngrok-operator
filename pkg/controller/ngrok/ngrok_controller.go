package ngrok

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	ngrokv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/pkg/apis/ngrok/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ngrok")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Ngrok Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNgrok{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ngrok-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Ngrok
	err = c.Watch(&source.Kind{Type: &ngrokv1alpha1.Ngrok{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Ngrok

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ngrokv1alpha1.Ngrok{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ngrokv1alpha1.Ngrok{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNgrok implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNgrok{}

// ReconcileNgrok reconciles a Ngrok object
type ReconcileNgrok struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Ngrok object and makes changes based on the status send
// and what is in the Meetup.Spec
func (r *ReconcileNgrok) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Ngrok")

	// Fetch the Ngrok instance
	ngrok := &ngrokv1alpha1.Ngrok{}

	// Populate Ngrok instance
	err := r.client.Get(context.TODO(), request.NamespacedName, ngrok)

	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Ngrok resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get Ngrok.")
		return reconcile.Result{}, err
	}

	// Define a new ngrok Pod configuration
	configmap := newNgrokConfigMap(ngrok)
	pod := newNgrokPod(ngrok)

	// Set Ngrok instance as the owner and controller
	if err := controllerutil.SetControllerReference(ngrok, configmap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Ngrok instance as the owner and controller
	if err := controllerutil.SetControllerReference(ngrok, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this ConfigMap already exists
	foundConfigMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, foundConfigMap)

	// if the ConfigMap is not found
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ngrok configMap", "Pod.Namespace", configmap.Namespace, "Pod.Name", configmap.Name)
		err = r.client.Create(context.TODO(), configmap)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	foundPod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundPod)

	// if the pod is not found
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ngrok Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			ngrok.Status.Status = "error"
			err := r.client.Status().Update(context.TODO(), ngrok)

			return reconcile.Result{}, err
		}

		ngrok.Status.Status = "created"
		ngrok.Status.URL = "fetching"
		err = r.client.Status().Update(context.TODO(), ngrok)

		return reconcile.Result{}, nil
	} else if err != nil {
		ngrok.Status.Status = "error"
		err := r.client.Status().Update(context.TODO(), ngrok)

		return reconcile.Result{}, err
	}

	time.Sleep(60 * time.Second)

	if foundPod.Status.PodIP != "" {
		fmt.Println(foundPod.Status.PodIP)
		fmt.Println("http://" + foundPod.Status.PodIP + ":4040/api/tunnels")
		response, _ := http.Get("http://" + foundPod.Status.PodIP + ":4040/api/tunnels")
		fmt.Println(response)
	}

	ngrok.Status.URL = "http://xx"
	err = r.client.Status().Update(context.TODO(), ngrok)

	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", foundPod.Namespace, "Pod.Name", foundPod.Name)
	return reconcile.Result{}, nil
}

func newNgrokConfigMap(cr *ngrokv1alpha1.Ngrok) *corev1.ConfigMap {
	configMapData := make(map[string]string, 0)
	ngrokProperties := `
web_addr: 0.0.0.0:4040`
	configMapData["ngrok.conf"] = ngrokProperties

	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-cm-ngrok",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Data: configMapData,
	}
}

func newNgrokPod(cr *ngrokv1alpha1.Ngrok) *corev1.Pod {
	ngrokPort := int32(4040)
	labels := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-ngrok",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "ngrok",
					Image:   "wernight/ngrok",
					Command: []string{"ngrok", "http", "--config", "/ngrok/ngrok.conf", cr.Spec.Service + ":" + strconv.FormatInt(int64(cr.Spec.Port), 10)},
					Ports: []corev1.ContainerPort{
						{ContainerPort: ngrokPort},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      cr.Name + "-cm-ngrok",
							MountPath: "/ngrok",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: cr.Name + "-cm-ngrok",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: cr.Name + "-cm-ngrok",
							},
						},
					},
				},
			},
		},
	}
}

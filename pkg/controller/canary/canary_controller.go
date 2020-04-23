/*
Copyright 2018 Wantedly, Inc..

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

package canary

import (
	"context"

	canaryv1beta1 "github.com/wantedly/canary-controller/pkg/apis/canary/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Canary Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCanary{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("canary-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Canary
	err = c.Watch(&source.Kind{Type: &canaryv1beta1.Canary{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by Canary - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &canaryv1beta1.Canary{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCanary{}

// ReconcileCanary reconciles a Canary object
type ReconcileCanary struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Canary object and makes changes based on the state read
// and what is in the Canary.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=canary.k8s.wantedly.com,resources=canaries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=canary.k8s.wantedly.com,resources=canaries/status,verbs=get;update;patch
func (r *ReconcileCanary) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Canary instance
	instance := &canaryv1beta1.Canary{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// TODO(munisystem): Set a status into the Canary resource if the target deployment doesn't exist
	target, err := r.getDeployment(instance.Spec.TargetDeploymentName, instance.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	copied := target.DeepCopy()

	// Inject data into Canary's Deployment
	labels := make(map[string]string, len(copied.GetLabels())+1)
	labels["canary"] = "true"
	for key, value := range copied.GetLabels() {
		labels[key] = value
	}

	spec := copied.Spec
	spec.Template.Spec.Hostname = "canary"
	spec.Selector.MatchLabels["canary"] = "true"
	spec.Template.Labels["canary"] = "true"

	containers := make(map[string]canaryv1beta1.CanaryContainer, 0)
	for _, container := range instance.Spec.TargetContainers {
		containers[container.Name] = container
	}
	for i := range spec.Template.Spec.Containers {
		if container, ok := containers[spec.Template.Spec.Containers[i].Name]; ok {
			spec.Template.Spec.Containers[i].Image = container.Image
		}
		spec.Template.Spec.Containers[i].Env = append(spec.Template.Spec.Containers[i].Env, corev1.EnvVar{
			Name:  "CANARY_ENABLED",
			Value: "1",
		})
	}
	canary := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        copied.ObjectMeta.Name + "-canary",
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: copied.ObjectMeta.Annotations,
		},
		Spec: spec,
	}

	if err := controllerutil.SetControllerReference(instance, canary, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	_, err = r.getDeployment(canary.Name, canary.Namespace)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment for Canary server", "namespace", canary.Namespace, "name", canary.Name)
		err = r.Create(context.TODO(), canary)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCanary) getDeployment(name, namespace string) (*appsv1.Deployment, error) {
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, found)
	return found, err
}

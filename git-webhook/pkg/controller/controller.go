/*
Copyright 2019 Pressinfra SRL.

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"

	"github.com/presslabs/stack/git-webhook/pkg/notifier"
	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

const controllerName = "git-webhook-controller"

// Add creates a new Git Webhook Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSite{
		Client:   mgr.GetClient(),
		Log:      logf.Log.WithName(controllerName),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder(controllerName),
		notifier: notifier.DefaultNotifier(),
	}
}

// getKey returns a string that represents the key under which cluster is registered
func getKey(meta metav1.Object) client.ObjectKey {
	return client.ObjectKey{
		Namespace: meta.GetNamespace(),
		Name:      meta.GetName(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, re reconcile.Reconciler) error {
	r, ok := re.(*ReconcileSite)
	if !ok {
		return fmt.Errorf("%T is not of type *controller.ReconcileSite", re)
	}

	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: re})
	if err != nil {
		return err
	}

	// Watch for changes to Wordpress
	err = c.Watch(&source.Kind{Type: &wordpressv1alpha1.Wordpress{}}, &handler.Funcs{
		CreateFunc: func(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
			if evt.Meta == nil {
				r.Log.Error(fmt.Errorf("CreateEvent received with no metadata"), "CreateEvent", evt)
				return
			}
			key := getKey(evt.Meta)
			r.notifier.RegisterSite(key, "github.com/calind/site", "master")
			r.Log.V(1).Info("successfully registred", "key", key)
		},
		DeleteFunc: func(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
			if evt.Meta == nil {
				r.Log.Error(fmt.Errorf("CreateEvent received with no metadata"), "CreateEvent", evt)
				return
			}
			key := getKey(evt.Meta)
			r.notifier.UnregisterSite(key)
			r.Log.V(1).Info("successfully unregisted", "key", key)
		},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSite{}

// ReconcileSite reconciles a Wordpress object
type ReconcileSite struct {
	client.Client
	Log      logr.Logger
	scheme   *runtime.Scheme
	recorder record.EventRecorder
	notifier *notifier.Notifier
}

// Reconcile reads that state of the cluster for a Wordpress object and makes changes based on the state read
// and what is in the Wordpress.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=wordpress.presslabs.org,resources=wordpresses,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileSite) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Site instance
	wp := &wordpressv1alpha1.Wordpress{}

	err := r.Get(context.TODO(), request.NamespacedName, wp)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

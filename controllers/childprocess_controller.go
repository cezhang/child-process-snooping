/*
Copyright 2022.

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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	childprocessv1 "cezhang/childprocess/api/v1"
)

// ChildprocessReconciler reconciles a Childprocess object
type ChildprocessReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=childprocess.cezhang,resources=childprocesses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=childprocess.cezhang,resources=childprocesses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=childprocess.cezhang,resources=childprocesses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Childprocess object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ChildprocessReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	var cp childprocessv1.Childprocess
	if err := r.Get(ctx, req.NamespacedName, &cp); err != nil {
		log.Error(err, "unable to fetch childprocess")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get the target pod
	var tpod v1.Pod
	key := client.ObjectKey{Namespace: req.Namespace, Name: cp.Spec.Tpod}
	if err := r.Get(ctx, key, &tpod); err != nil {
		log.Error(err, fmt.Sprintf("unable to fetch tpod %s", key.String()))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	tpodScheduleNode := tpod.Spec.NodeName
	fmt.Println(time.Now())
	fmt.Println(tpodScheduleNode)
	cp.Status.Tpod = tpod.Status
	fmt.Println("mpod:" + cp.Status.Mpod)

	if len(cp.Status.Mpod) == 0 {
		if err := r.createMpod(ctx, &cp, &tpod); err != nil {
			log.Error(err, "unable to create mpod ")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	} else {
		var mpod v1.Pod
		key := client.ObjectKey{Namespace: req.Namespace, Name: cp.Status.Mpod}
		if err := r.Get(ctx, key, &mpod); err != nil {
			log.Error(err, fmt.Sprintf("unable to fetch tpod %s", key.String()))

			if err := r.createMpod(ctx, &cp, &tpod); err != nil {
				log.Error(err, "unable to create mpod ")
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
		} else {
			//do nothing
		}
	}

	if err := r.Status().Update(ctx, &cp); err != nil {
		log.Error(err, "unable to update CronJob status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChildprocessReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&childprocessv1.Childprocess{}).
		Complete(r)
}

func (r *ChildprocessReconciler) createMpod(ctx context.Context, cp *childprocessv1.Childprocess, tpod *v1.Pod) error {
	cp.Spec.Mpod.NodeName = tpod.Spec.NodeName
	mpodNmae := fmt.Sprintf("%s-mpod", tpod.Name)
	mpod := v1.Pod{
		TypeMeta: ctrl.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: ctrl.ObjectMeta{
			Name:      mpodNmae,
			Namespace: tpod.Namespace,
		},
		Spec: cp.Spec.Mpod,
	}
	if err := r.Create(ctx, &mpod); err != nil {
		return err
	}

	cp.Status.Mpod = mpodNmae
	return nil
}

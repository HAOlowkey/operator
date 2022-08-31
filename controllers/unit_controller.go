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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	applicationv1 "github.com/HAOlowkey/operator/api/v1"
)

// UnitReconciler reconciles a Unit object
type UnitReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=application.haolowkey.github.io,resources=units,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=application.haolowkey.github.io,resources=units/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=application.haolowkey.github.io,resources=units/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Unit object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *UnitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here
	logger.Info("fetching Unit resource")
	unit := &applicationv1.Unit{}
	if err := r.Get(ctx, req.NamespacedName, unit); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Unit resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Unit resource.")
		return ctrl.Result{}, err
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      unit.Name,
			Namespace: unit.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  unit.Spec.MainContainer.Name,
					Image: unit.Spec.MainContainer.Repository + ":" + unit.Spec.MainContainer.Tag,
					Ports: []v1.ContainerPort{
						{
							Name:          unit.Spec.MainContainer.Ports[0].Name,
							ContainerPort: unit.Spec.MainContainer.Ports[0].Port,
							Protocol:      v1.Protocol(unit.Spec.MainContainer.Ports[0].Protocol),
						},
					},
				},
				{
					Name:  unit.Spec.SideCarContainer[0].Name,
					Image: unit.Spec.SideCarContainer[0].Repository + ":" + unit.Spec.SideCarContainer[0].Tag,
					Ports: []v1.ContainerPort{
						{
							Name:          unit.Spec.SideCarContainer[0].Ports[0].Name,
							ContainerPort: unit.Spec.SideCarContainer[0].Ports[0].Port,
							Protocol:      v1.Protocol(unit.Spec.SideCarContainer[0].Ports[0].Protocol),
						},
					},
				},
				{
					Name:  unit.Spec.SideCarContainer[1].Name,
					Image: unit.Spec.SideCarContainer[1].Repository + ":" + unit.Spec.SideCarContainer[1].Tag,
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(unit, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	//查找同名deployment
	d := &v1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, d); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, pod); err != nil {
				logger.Error(err, "create Pod resource failed")
				return ctrl.Result{}, err
			}
			r.Recorder.Eventf(unit, v1.EventTypeNormal, "Created", "create pod %s", pod.Name)
			logger.Info("created Pod resource for Unit")
		}
	} else {
		logger.Info("existing Pod resource already exists for Unit, checking Image\")")
		r.Recorder.Eventf(d, v1.EventTypeNormal, "Created", "test %s", pod.Name)
		if unit.Spec.MainContainer.Repository+":"+unit.Spec.MainContainer.Tag != d.Spec.Containers[0].Image {
			d.Spec.Containers[0].Image = unit.Spec.MainContainer.Repository + ":" + unit.Spec.MainContainer.Tag
			if err = r.Update(ctx, d); err != nil {
				logger.Error(err, "update main Container image failed!")
				return ctrl.Result{}, err
			}
			r.Recorder.Eventf(unit, v1.EventTypeNormal, "Updated", "Change unit %s log sidecar image to %s", unit.Name, d.Spec.Containers[0].Image)

		} else if unit.Spec.SideCarContainer[0].Repository+":"+unit.Spec.SideCarContainer[0].Tag != d.Spec.Containers[1].Image {
			d.Spec.Containers[1].Image = unit.Spec.SideCarContainer[0].Repository + ":" + unit.Spec.SideCarContainer[0].Tag
			if err = r.Update(ctx, d); err != nil {
				logger.Error(err, "update main Container image failed!")
				return ctrl.Result{}, err
			}
			r.Recorder.Eventf(unit, v1.EventTypeNormal, "Updated", "Change unit %s monitor sidecar image to %s", unit.Name, d.Spec.Containers[1].Image)

		} else if unit.Spec.SideCarContainer[1].Repository+":"+unit.Spec.SideCarContainer[1].Tag != d.Spec.Containers[2].Image {
			d.Spec.Containers[2].Image = unit.Spec.SideCarContainer[1].Repository + ":" + unit.Spec.SideCarContainer[1].Tag
			if err = r.Update(ctx, d); err != nil {
				logger.Error(err, "update main Container image failed!")
				return ctrl.Result{}, err
			}
			r.Recorder.Eventf(unit, v1.EventTypeNormal, "Updated", "Change unit %s log sidecar image to %s", unit.Name, d.Spec.Containers[0].Image)
		} else {
			logger.Info("Pod has no change no need to update")
		}
	}

	unit.Status.ContainerCount = int64(len(unit.Spec.SideCarContainer) + 1)
	r.Status().Update(ctx, unit)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UnitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&applicationv1.Unit{}).
		Owns(&v1.Pod{}).
		Complete(r)
}

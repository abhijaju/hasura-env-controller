package reconciler

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	parser "github.com/abhijaju/hasura-env-controller/operator/configmapparser"
)

type EnvReconciler struct {
	client.Client
	log logr.Logger
}

func NewEnvReconciler(c client.Client, log logr.Logger) *EnvReconciler {
	r := EnvReconciler{
		Client: c,
		log:    log,
	}
	return &r
}

func (r *EnvReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if !(req.Name == "environment" && req.Namespace == "env-controller") {
		return ctrl.Result{}, nil
	}
	r.log.Info("Reconciling for config map", "name", req.Name, "namespace", req.Namespace)

	// get configmap
	cm := corev1.ConfigMap{}
	if err := r.Get(ctx, req.NamespacedName, &cm); err != nil {
		if errors.IsNotFound(err) {
			r.log.Info("configmap not found, ignoring for now")
			return ctrl.Result{}, nil
		}
		r.log.Error(err, "error getting configmap")
		return ctrl.Result{}, err
	}

	// parse configmap
	cmData, err := parser.NewConfigMapData(cm.Data["data"])
	if err != nil {
		r.log.Error(err, "error parsing configmap")
		return ctrl.Result{}, err
	}
	r.log.Info("configmap contents", "cmData", cmData)

	// get deployment
	dep := appsv1.Deployment{}
	key := types.NamespacedName{
		Name:      cmData.DeploymentName,
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, key, &dep); err != nil {
		if errors.IsNotFound(err) {
			r.log.Info("deployment not found, ignoring for now")
			return ctrl.Result{}, nil
		}
		r.log.Error(err, "error getting deployment")
		return ctrl.Result{}, err
	}

	// updating deployment
	r.updateDeploymentEnvVars(&dep, cmData.EnvVars)
	if err = r.Update(ctx, &dep); err != nil {
		r.log.Error(err, "error updating deployment")
		return ctrl.Result{}, err
	}

	r.log.Info("Reconciled successfully for configmap", "name", req.Name, "namespace", req.Namespace)
	return ctrl.Result{}, nil
}

func (r *EnvReconciler) updateDeploymentEnvVars(dep *appsv1.Deployment, envVars map[string]string) {
	r.log.Info("updating deployment env variables", "newVars", envVars, "existingVars", dep.Spec.Template.Spec.Containers[0].Env)

	finalEnv := make([]corev1.EnvVar, 0)
	// update existing Env Variables
	for _, existingEnvVar := range dep.Spec.Template.Spec.Containers[0].Env {
		val, present := envVars[existingEnvVar.Name]
		if present {
			finalEnv = append(finalEnv, corev1.EnvVar{Name: existingEnvVar.Name, Value: val})
			delete(envVars, existingEnvVar.Name)
			continue
		}
		finalEnv = append(finalEnv, existingEnvVar)
	}

	for key, val := range envVars {
		finalEnv = append(finalEnv, corev1.EnvVar{Name: key, Value: val})
	}

	dep.Spec.Template.Spec.Containers[0].Env = finalEnv
	r.log.Info("updated deployment env variables", "finalVars", dep.Spec.Template.Spec.Containers[0].Env)
}

package main

import (
	"flag"
	"os"

	corev1 "k8s.io/api/core/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	reconciler "github.com/abhijaju/hasura-env-controller/operator/reconciler"
)

func main() {
	zapOpts := zap.Options{
		Development: false,
	}
	zapOpts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zapOpts)))

	log := ctrl.Log.WithName("env-controller")
	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	r := reconciler.NewEnvReconciler(manager.GetClient(), log)
	if err = ctrl.NewControllerManagedBy(manager).For(&corev1.ConfigMap{}).Complete(r); err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}

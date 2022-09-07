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

package main

import (
	"flag"
	"fmt"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	examplecomv1alpha1 "github.com/mresvanis/he-sample-operator/api/v1alpha1"
	"github.com/mresvanis/he-sample-operator/controllers"
	"github.com/mresvanis/he-sample-operator/internal/conditions"
	"github.com/mresvanis/he-sample-operator/internal/finalizers"
	"github.com/mresvanis/he-sample-operator/internal/module"
	"github.com/mresvanis/he-sample-operator/internal/nodeselector"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(examplecomv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	klog.InitFlags(flag.CommandLine)

	flag.Parse()

	logger := klogr.New()

	ctrl.SetLogger(logger)

	setupLogger := logger.WithName("setup")

	watchNamespace, err := getWatchNamespace()
	if err != nil {
		setupLogger.Error(err, "unable to get WatchNamespace, "+
			"the manager will watch and manage resources in all namespaces")
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "e436e148.example.com",
		Namespace:              watchNamespace,
	})
	if err != nil {
		setupLogger.Error(err, "unable to start manager")
		os.Exit(1)
	}

	c := mgr.GetClient()
	s := mgr.GetScheme()

	mr := module.NewReconciler(c, s)
	fu := finalizers.NewUpdater(c)
	cu := conditions.NewUpdater(c)
	nsv := nodeselector.NewValidator(c)
	dcc := controllers.NewDeviceConfigReconciler(c, s, mgr.GetEventRecorderFor("deviceconfig-controller"), mr, fu, cu, nsv)

	if err = dcc.SetupWithManager(mgr); err != nil {
		setupLogger.Error(err, "unable to create controller", "controller", "DeviceConfig")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLogger.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLogger.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLogger.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLogger.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func getWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar si the contant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	var watchNamespaceEnvVar = "WATCH_NAMESPACE"

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}

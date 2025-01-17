// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package componenttest

import (
	"github.com/nokia/industrial-application-framework/application-lib/pkg/config"
	"github.com/nokia/industrial-application-framework/application-lib/pkg/handlers"
	"github.com/nokia/industrial-application-framework/application-lib/pkg/k8sdynamic"
	"github.com/nokia/industrial-application-framework/application-lib/pkg/kubelib"
	ctenv "github.com/nokia/industrial-application-framework/componenttest-lib/pkg/env"
	appdacnokiacomv1alpha1 "github.com/nokia/industrial-application-framework/consul-operator/api/v1alpha1"
	"github.com/nokia/industrial-application-framework/consul-operator/controllers"
	"github.com/nokia/industrial-application-framework/consul-operator/pkg/licenceexpired"
	"github.com/nokia/industrial-application-framework/consul-operator/pkg/monitoring"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"path/filepath"
	"runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

const ItBinaryRelativePath = "/../componenttest/resources"

var _ = BeforeSuite(func() {
	ctenv.TearUpTestEnv(getTestBinaryPath(ItBinaryRelativePath))

	k8sdynamic.Config = ctenv.Cfg
	kubelib.Config = ctenv.Cfg

	k8sManager, err := ctrl.NewManager(ctenv.Cfg, ctrl.Options{
		MetricsBindAddress: ":8383",
		Scheme:             ourScheme,
	})
	Expect(err).ToNot(HaveOccurred())

	k8sClient = k8sManager.GetClient()

	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	operatorConfiguration := config.OperatorConfig{
		ApplicationName:             "Consul",
		RuntimeDeploymentPath:       "../deployment",
		AppDeploymentDirName:        "app-deployment",
		RuntimeResReqPath:           "../deployment/resource-reqs-generated",
		ResReqDirName:               "resource-reqs",
		KubernetesAppDeploymentName: consulStatefulSetName,
		AppPnaName:                  "private-network-for-consul",
		Template: config.TemplateConfig{
			LeftDelimiter:  "[[",
			RightDelimiter: "]]",
		}}

	reconciler := controllers.AppSpecificReconciler{
		Common: handlers.OperatorReconciler{
			Client:        k8sManager.GetClient(),
			Scheme:        k8sManager.GetScheme(),
			Configuration: operatorConfiguration,
			Functions: handlers.ReconcilerHookFunctions{
				CreateAppCr:                 appdacnokiacomv1alpha1.CreateAppInstance,
				CreateAppStatusMonitor:      monitoring.CreateAppStatusMonitor,
				CreateLicenceExpiredHandler: licenceexpired.CreateLicenseExpiredHandler,
			},
		},
	}

	err = controllers.SetupWithManager(k8sManager, &reconciler)

	Expect(err).ToNot(HaveOccurred())

	log.Info("Starting the Operator")
	go func() {
		defer GinkgoRecover()
		Expect(k8sManager.Start(ctrl.SetupSignalHandler())).NotTo(HaveOccurred())
	}()

})

var _ = AfterSuite(func() {
	ctenv.ResetEtcd()

	ctenv.TearDownTestEnv()
})

var ourScheme = k8sruntime.NewScheme()

func init() {

	utilruntime.Must(clientgoscheme.AddToScheme(ourScheme))

	utilruntime.Must(appdacnokiacomv1alpha1.AddToScheme(ourScheme))
	//+kubebuilder:scaffold:scheme
}

func getTestBinaryPath(testBinariesRelativePath string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return basepath + testBinariesRelativePath
}

func TestConsulOperator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consul Operator Component Test Suite")

}

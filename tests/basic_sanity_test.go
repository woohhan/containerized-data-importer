package tests_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerized-data-importer/tests"
	"kubevirt.io/containerized-data-importer/tests/framework"
)

var _ = Describe("[rfe_id:1347][crit:high][vendor:cnv-qe@redhat.com][level:component]Basic Sanity", func() {
	f := framework.NewFrameworkOrDie("sanity", framework.Config{
		SkipNamespaceCreation: true,
	})

	Context("[test_id:1348]CDI service account should exist", func() {
		It("Should succeed", func() {
			result, err := tests.RunKubectlCommand(f, "get", "sa", "cdi-sa", "-n", f.CdiInstallNs)
			Expect(err).To(BeNil())
			Expect(result).To(ContainSubstring("cdi-sa"))
		})
	})

	Context("[test_id:1349]CDI Cluster role should exist", func() {
		It("Should succeed", func() {
			result, err := tests.RunKubectlCommand(f, "get", "clusterrole", "cdi")
			Expect(err).To(BeNil())
			Expect(result).To(ContainSubstring("cdi"))
		})
	})

	Context("[test_id:1350]CDI Cluster role binding should exist", func() {
		It("Should succeed", func() {
			result, err := tests.RunKubectlCommand(f, "get", "clusterrolebinding", "cdi-sa")
			Expect(err).To(BeNil())
			Expect(result).To(ContainSubstring("cdi-sa"))
		})
	})

	Context("CDI deployment should exist", func() {
		It("[test_id:1351]Should succeed", func() {
			result, err := tests.RunKubectlCommand(f, "get", "deployment", "cdi-deployment", "-n", f.CdiInstallNs)
			Expect(err).To(BeNil())
			Expect(result).To(ContainSubstring("cdi-deployment"))
		})
		It("[test_id:1352]There should be 1 replica", func() {
			result, err := tests.RunKubectlCommand(f, "get", "deployment", "cdi-deployment", "-o", "jsonpath={.spec.replicas}", "-n", f.CdiInstallNs)
			Expect(err).To(BeNil())
			Expect(result).To(ContainSubstring("1"))
		})
	})

	Context("cdi-sa RBAC rules are correct", func() {
		It("[test_id:1353]rules should match expectation", func() {
			sa := fmt.Sprintf("system:serviceaccount:" + f.CdiInstallNs + ":cdi-sa")

			eventExpectedResult := make(map[string]string)
			eventExpectedResult["get"] = "no"
			eventExpectedResult["list"] = "no"
			eventExpectedResult["watch"] = "no"
			eventExpectedResult["delete"] = "no"
			eventExpectedResult["create"] = "yes"
			eventExpectedResult["update"] = "no"
			eventExpectedResult["patch"] = "no"
			eventExpectedResult["deletecollection"] = "no"
			ValidateRBACForResource(f, eventExpectedResult, "events", sa)

			pvcExpectedResult := make(map[string]string)
			pvcExpectedResult["get"] = "yes"
			pvcExpectedResult["list"] = "yes"
			pvcExpectedResult["watch"] = "yes"
			pvcExpectedResult["delete"] = "yes"
			pvcExpectedResult["create"] = "yes"
			pvcExpectedResult["update"] = "yes"
			pvcExpectedResult["patch"] = "no"
			pvcExpectedResult["deletecollection"] = "no"
			ValidateRBACForResource(f, pvcExpectedResult, "persistentvolumeclaims", sa)
			ValidateRBACForResource(f, pvcExpectedResult, "persistentvolumeclaims/finalizers", sa)

			podExpectedResult := make(map[string]string)
			podExpectedResult["get"] = "yes"
			podExpectedResult["list"] = "yes"
			podExpectedResult["watch"] = "yes"
			podExpectedResult["delete"] = "yes"
			podExpectedResult["create"] = "yes"
			podExpectedResult["update"] = "no"
			podExpectedResult["patch"] = "no"
			podExpectedResult["deletecollection"] = "no"
			ValidateRBACForResource(f, podExpectedResult, "pods", sa)
			ValidateRBACForResource(f, podExpectedResult, "pods/finalizers", sa)

			secretsExpectedResult := make(map[string]string)
			secretsExpectedResult["get"] = "yes"
			secretsExpectedResult["list"] = "yes"
			secretsExpectedResult["watch"] = "yes"
			secretsExpectedResult["delete"] = "no"
			secretsExpectedResult["create"] = "yes"
			secretsExpectedResult["update"] = "no"
			secretsExpectedResult["patch"] = "no"
			secretsExpectedResult["deletecollection"] = "no"
			ValidateRBACForResource(f, secretsExpectedResult, "secrets", sa)
		})
	})
})

func ValidateRBACForResource(f *framework.Framework, expectedResults map[string]string, resource string, sa string) {
	for verb, expectedRes := range expectedResults {
		By(fmt.Sprintf("verifying cdi-sa "+resource+" rules, for verb %s", verb))
		result, err := tests.RunKubectlCommand(f, "auth", "can-i", "--as", sa, verb, resource)
		if expectedRes != "no" {
			Expect(err).ToNot(HaveOccurred())
		}
		Expect(result).To(ContainSubstring(expectedRes))
	}
}

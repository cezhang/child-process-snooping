/*
Ideally, we should have one `<kind>_controller_test.go` for each controller scaffolded and called in the `suite_test.go`.
So, let's write our example test for the CronJob controller (`cronjob_controller_test.go.`)
*/

/*

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
// +kubebuilder:docs-gen:collapse=Apache License

/*
As usual, we start with the necessary imports. We also define some utility variables.
*/
package controllers

import (
	childprocessv1 "cezhang/childprocess/api/v1"
	"context"
	"fmt"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

)

// +kubebuilder:docs-gen:collapse=Imports

/*
The first step to writing a simple integration test is to actually create an instance of CronJob you can run tests against.
Note that to create a CronJob, you’ll need to create a stub CronJob struct that contains your CronJob’s specifications.

Note that when we create a stub CronJob, the CronJob also needs stubs of its required downstream objects.
Without the stubbed Job template spec and the Pod template spec below, the Kubernetes API will not be able to
create the CronJob.
*/
var _ = Describe("Childprocess controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		CPName      = "test-cp"
		CPNamespace = "default"
		TPOD 		= "tpod"
		timeout  = time.Second * 60
		duration = time.Second * 10
		interval = time.Second * 3
	)

	Context("When updating CP Status", func() {
		It("Should update status tpod when target pod is set", func() {
			By("By creating a new Childprocess")
			ctx := context.Background()

			// create target pod first
			tpod := v1.Pod{
				TypeMeta: ctrl.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
				ObjectMeta: ctrl.ObjectMeta{
					Name:      TPOD,
					Namespace: CPNamespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "busybox",
							Image: "busybox",
							Command: []string{"sleep", "3000"},
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, &tpod)).Should(Succeed())
			Eventually(func() bool {
				etpod := &v1.Pod{}
				key := client.ObjectKey{Namespace: CPNamespace, Name: TPOD}
				err := k8sClient.Get(ctx, key, etpod)
				if err != nil {
					return false
				}
				fmt.Println("etpod:", etpod.Status.Phase)
				if etpod.Status.Phase == v1.PodRunning {
					return true
				}else{
					return false
				}
			}, timeout, interval).Should(BeTrue())



			// create CRD
			cp := &childprocessv1.Childprocess{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "childprocess.cezhang/v1",
					Kind:       "Childprocess",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      CPName,
					Namespace: CPNamespace,
				},
				Spec: childprocessv1.ChildprocessSpec{
					Tpod: TPOD,
					Mpod: v1.PodSpec {
						Containers: []v1.Container{
							{
								Name: "mpod",
								Image: "nginx",
							},
						},
					},

				},
			}
			Expect(k8sClient.Create(ctx, cp)).Should(Succeed())

			/*
				After creating this CronJob, let's check that the CronJob's Spec fields match what we passed in.
				Note that, because the k8s apiserver may not have finished creating a CronJob after our `Create()` call from earlier, we will use Gomega’s Eventually() testing function instead of Expect() to give the apiserver an opportunity to finish creating our CronJob.

				`Eventually()` will repeatedly run the function provided as an argument every interval seconds until
				(a) the function’s output matches what’s expected in the subsequent `Should()` call, or
				(b) the number of attempts * interval period exceed the provided timeout value.

				In the examples below, timeout and interval are Go Duration values of our choosing.
			*/

			cpLookupKey := types.NamespacedName{Name: CPName, Namespace: CPNamespace}
			createdCP := &childprocessv1.Childprocess{}
			fmt.Println("debug1:", createdCP.Spec.Tpod)
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cpLookupKey, createdCP)
				fmt.Println("debug:", createdCP.Spec.Tpod)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(cp.Spec.Tpod).Should(Equal(TPOD))


			createdCP1 := &childprocessv1.Childprocess{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cpLookupKey, createdCP1)
				fmt.Println("createdCP1", createdCP1)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			//Expect(cp.Status.Mpod).Should(Equal(TPOD+"-mpod"))

		})
	})

})

/*
	After writing all this code, you can run `go test ./...` in your `controllers/` directory again to run your new test!
*/

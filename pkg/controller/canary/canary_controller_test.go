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
	"testing"
	"time"

	"github.com/onsi/gomega"
	canaryv1beta1 "github.com/wantedly/canary-controller/pkg/apis/canary/v1beta1"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	c      client.Client
	target = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default", Labels: map[string]string{"app": "nginx"}},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment": "foo-deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"deployment": "foo-deployment"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
	targetKey       = types.NamespacedName{Name: "foo", Namespace: "default"}
	expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo-canary", Namespace: "default"}}
	depKey          = types.NamespacedName{Name: "foo-canary", Namespace: "default"}
)

const timeout = time.Second * 5

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &canaryv1beta1.Canary{
		ObjectMeta: metav1.ObjectMeta{Name: "foo-canary", Namespace: "default"},
		Spec: canaryv1beta1.CanarySpec{
			TargetDeploymentName: "foo",
			TargetContainers: []canaryv1beta1.CanaryContainer{
				{
					Name:  "nginx",
					Image: "nginx:latest",
				},
			},
		},
	}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the target Deployment for testing
	err = c.Create(context.TODO(), target)
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), target)
	g.Eventually(func() error { return c.Get(context.TODO(), targetKey, &appsv1.Deployment{}) }, timeout).
		Should(gomega.Succeed())

	// Create the Canary object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	deploy := &appsv1.Deployment{}
	g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
		Should(gomega.Succeed())

	// Delete the Deployment and expect Reconcile to be called for Deployment deletion
	g.Expect(c.Delete(context.TODO(), deploy)).NotTo(gomega.HaveOccurred())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
		Should(gomega.Succeed())
}

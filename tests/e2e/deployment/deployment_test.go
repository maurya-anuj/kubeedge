/*
Copyright 2019 The KubeEdge Authors.

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

package deployment

import (
	"bytes"
	"net/http"
	"time"

	"github.com/kubeedge/kubeedge/tests/e2e/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	AppHandler        = "/api/v1/namespaces/default/pods"
	NodeHandler       = "/api/v1/nodes"
	DeploymentHandler = "/apis/apps/v1/namespaces/default/deployments"
	ServiceHandler    = "/api/v1/namespaces/default/services"
	EndpointHandler   = "/api/v1/namespaces/default/endpoints"
)

var DeploymentTestTimerGroup *utils.TestTimerGroup = utils.NewTestTimerGroup()

//Run Test cases
var _ = Describe("Application deployment test in E2E scenario", func() {
	var UID string
	var UIDServer string
	var UIDClient string
	var testTimer *utils.TestTimer
	var testDescription GinkgoTestDescription
	Context("Test application deployment and delete deployment using deployment spec", func() {
		BeforeEach(func() {
			// Get current test description
			testDescription = CurrentGinkgoTestDescription()
			// Start test timer
			testTimer = DeploymentTestTimerGroup.NewTestTimer(testDescription.TestText)
		})
		AfterEach(func() {
			// End test timer
			testTimer.End()
			// Print result
			testTimer.PrintResult()
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler)
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UID {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					StatusCode := utils.DeleteDeployment(ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, deployment.Name)
					Expect(StatusCode).Should(Equal(http.StatusOK))
				}
			}
			utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			utils.PrintTestcaseNameandStatus()
		})

		It("E2E_APP_DEPLOYMENT_1: Create deployment and check the pods are coming up correctly", func() {
			var deploymentList v1.DeploymentList
			var podlist metav1.PodList
			replica := 1
			//Generate the random string and assign as a UID
			UID = "edgecore-depl-app-" + utils.GetRandomString(5)
			IsAppDeployed := utils.HandleDeployment(false, false, http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UID, ctx.Cfg.AppImageUrl[1], nodeSelector, "", replica)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler)
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UID {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					break
				}
			}
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
		})
		It("E2E_APP_DEPLOYMENT_2: Create deployment with replicas and check the pods are coming up correctly", func() {
			var deploymentList v1.DeploymentList
			var podlist metav1.PodList
			replica := 3
			//Generate the random string and assign as a UID
			UID = "edgecore-depl-app-" + utils.GetRandomString(5)
			IsAppDeployed := utils.HandleDeployment(false, false, http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UID, ctx.Cfg.AppImageUrl[1], nodeSelector, "", replica)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler)
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UID {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					break
				}
			}
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
		})

		It("E2E_APP_DEPLOYMENT_3: Create deployment and check deployment ctrler re-creating pods when user deletes the pods manually", func() {
			var deploymentList v1.DeploymentList
			var podlist metav1.PodList
			replica := 3
			//Generate the random string and assign as a UID
			UID = "edgecore-depl-app-" + utils.GetRandomString(5)
			IsAppDeployed := utils.HandleDeployment(false, false, http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UID, ctx.Cfg.AppImageUrl[1], nodeSelector, "", replica)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler)
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UID {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					break
				}
			}
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
			for _, pod := range podlist.Items {
				_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
				Expect(StatusCode).Should(Equal(http.StatusOK))
			}
			utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UID {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					break
				}
			}
			Expect(len(podlist.Items)).Should(Equal(replica))
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
		})

	})
	Context("Test application deployment using Pod spec", func() {
		BeforeEach(func() {
			// Get current test description
			testDescription = CurrentGinkgoTestDescription()
			// Start test timer
			testTimer = DeploymentTestTimerGroup.NewTestTimer(testDescription.TestText)
		})
		AfterEach(func() {
			// End test timer
			testTimer.End()
			// Print result
			testTimer.PrintResult()
			var podlist metav1.PodList
			label := nodeName
			podlist, err := utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
			Expect(err).To(BeNil())
			for _, pod := range podlist.Items {
				_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
				Expect(StatusCode).Should(Equal(http.StatusOK))
			}
			utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			utils.PrintTestcaseNameandStatus()
		})

		It("E2E_POD_DEPLOYMENT_1: Create a pod and check the pod is coming up correclty", func() {
			var podlist metav1.PodList
			//Generate the random string and assign as a UID
			UID = "pod-app-" + utils.GetRandomString(5)
			IsAppDeployed := utils.HandlePod(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+AppHandler, UID, ctx.Cfg.AppImageUrl[0], nodeSelector)
			Expect(IsAppDeployed).Should(BeTrue())
			label := nodeName
			podlist, err := utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
			Expect(err).To(BeNil())
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
		})

		It("E2E_POD_DEPLOYMENT_2: Create the pod and delete pod happening successfully", func() {
			var podlist metav1.PodList
			//Generate the random string and assign as a UID
			UID = "pod-app-" + utils.GetRandomString(5)
			IsAppDeployed := utils.HandlePod(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+AppHandler, UID, ctx.Cfg.AppImageUrl[0], nodeSelector)
			Expect(IsAppDeployed).Should(BeTrue())
			label := nodeName
			podlist, err := utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
			Expect(err).To(BeNil())
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
			for _, pod := range podlist.Items {
				_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
				Expect(StatusCode).Should(Equal(http.StatusOK))
			}
			utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
		})
		It("E2E_POD_DEPLOYMENT_3: Create pod and delete the pod successfully, and delete already deleted pod and check the behaviour", func() {
			var podlist metav1.PodList
			//Generate the random string and assign as a UID
			UID = "pod-app-" + utils.GetRandomString(5)
			IsAppDeployed := utils.HandlePod(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+AppHandler, UID, ctx.Cfg.AppImageUrl[0], nodeSelector)
			Expect(IsAppDeployed).Should(BeTrue())
			label := nodeName
			podlist, err := utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
			Expect(err).To(BeNil())
			utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
			for _, pod := range podlist.Items {
				_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
				Expect(StatusCode).Should(Equal(http.StatusOK))
			}
			utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + UID)
			Expect(StatusCode).Should(Equal(http.StatusNotFound))
		})
		It("E2E_POD_DEPLOYMENT_4: Create and delete pod multiple times and check all the Pod created and deleted successfully", func() {
			//Generate the random string and assign as a UID
			for i := 0; i < 10; i++ {
				UID = "pod-app-" + utils.GetRandomString(5)
				IsAppDeployed := utils.HandlePod(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+AppHandler, UID, ctx.Cfg.AppImageUrl[0], nodeSelector)
				Expect(IsAppDeployed).Should(BeTrue())
				label := nodeName
				podlist, err := utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
				Expect(err).To(BeNil())
				utils.WaitforPodsRunning(ctx.Cfg.K8SMasterForKubeEdge, podlist, 240*time.Second)
				for _, pod := range podlist.Items {
					_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
					Expect(StatusCode).Should(Equal(http.StatusOK))
				}
				utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			}
		})
	})
	Context("Test Pod communication with edgeMesh", func() {
		BeforeEach(func() {
			// Get current test description
			testDescription = CurrentGinkgoTestDescription()
			// Start test timer
			testTimer = DeploymentTestTimerGroup.NewTestTimer(testDescription.TestText)
		})
		AfterEach(func() {
			// End test timer
			testTimer.End()
			// Print result
			testTimer.PrintResult()
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				label := nodeName
				podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
				Expect(err).To(BeNil())
				StatusCode := utils.DeleteDeployment(ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, deployment.Name)
				Expect(StatusCode).Should(Equal(http.StatusOK))
			}
			utils.CheckPodDeleteState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)

			var serviceList metav1.ServiceList
			err = utils.GetServices(&serviceList, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler+utils.LabelSelector+"service%3Dtest")
			Expect(err).To(BeNil())
			for _, service := range serviceList.Items {
				StatusCode := utils.DeleteSvc(ctx.Cfg.K8SMasterForKubeEdge + ServiceHandler + "/" + service.Name)
				Expect(StatusCode).Should(Equal(http.StatusOK))
			}
			utils.PrintTestcaseNameandStatus()
		})

		FIt("E2E_SERVICE_EDGEMESH_1: Create two pods and check the pods are communicating or not: POSITIVE", func() {
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			var servicelist metav1.ServiceList
			//Generate the random string and assign as a UID
			UIDServer = "pod-app-server" + utils.GetRandomString(5)
			depobj := utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed := utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					break
				}
			}
			utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			utils.Info("\n Server app deployed \n")

			serviceName := "pod-app-server"
			// Deploy service over the server pod
			err = utils.ExposePodService(serviceName, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler, 80, intstr.FromInt(8000), metav1.ServiceTypeClusterIP)
			Expect(err).To(BeNil())
			err = utils.GetServices(&servicelist, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler)
			Expect(err).To(BeNil())

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep := utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))

			UIDClient = "pod-app-client" + utils.GetRandomString(5)
			depobj = utils.CreateDeployment(UIDClient, ctx.Cfg.AppImageUrl[3], nodeSelector, 1, map[string]string{"app": "client"}, 81)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDClient, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					// label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler+utils.LabelSelector+"app%3Dclient", "")
					Expect(err).To(BeNil())
					break
				}
			}
			utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
			// Check weather the name variable is changed in server
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))
		})

		FIt("E2E_SERVICE_EDGEMESH_2: Client pod restart: POSITIVE", func() {
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			var servicelist metav1.ServiceList
			//Generate the random string and assign as a UID
			UIDServer = "pod-app-server" + utils.GetRandomString(5)
			depobj := utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed := utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
				}
			}
			utils.Info("\n Server app deployed \n")

			serviceName := "pod-app-server"
			// Deploy service over the server pod
			err = utils.ExposePodService(serviceName, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler, 80, intstr.FromInt(8000), metav1.ServiceTypeClusterIP)
			Expect(err).To(BeNil())
			err = utils.GetServices(&servicelist, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler)
			Expect(err).To(BeNil())

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep := utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))

			//deploy client deployment
			UIDClient = "pod-app-client" + utils.GetRandomString(5)
			depobj = utils.CreateDeployment(UIDClient, ctx.Cfg.AppImageUrl[3], nodeSelector, 1, map[string]string{"app": "client"}, 81)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDClient, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
					break
				}
			}
			utils.Info("\n Client app deployed \n")

			//check name changed(communication happened)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))

			//delete client
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())

			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					label := deployment.Spec.Selector
					for key, value := range label.MatchLabels {
						podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler+utils.LabelSelector+key+"%3D"+value, "")
						Expect(err).To(BeNil())
						for _, pod := range podlist.Items {
							_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
							Expect(StatusCode).Should(Equal(http.StatusOK))
						}
					}
				}
			}

			//change the name back to default again
			var jsonStr = []byte("Default")
			_, err = http.Post(ep, "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				panic(err)
			}

			//deployment will restart it check again pod is there
			podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler+utils.LabelSelector+"app"+"%3D"+"client", "")
			Expect(err).To(BeNil())
			utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge, podlist)

			//check the name is changed of not
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))
		})

		FIt("E2E_SERVICE_EDGEMESH_3: Server pod restart: POSITIVE", func() {
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			var servicelist metav1.ServiceList
			//Generate the random string and assign as a UID
			UIDServer = "pod-app-server" + utils.GetRandomString(5)
			depobj := utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed := utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
				}
			}
			utils.Info("\n Server app deployed \n")

			serviceName := "pod-app-server"
			// Deploy service over the server pod
			err = utils.ExposePodService(serviceName, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler, 80, intstr.FromInt(8000), metav1.ServiceTypeClusterIP)
			Expect(err).To(BeNil())
			err = utils.GetServices(&servicelist, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler)
			Expect(err).To(BeNil())

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep := utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))

			//deploy client deployment
			UIDClient = "pod-app-client" + utils.GetRandomString(5)
			depobj = utils.CreateDeployment(UIDClient, ctx.Cfg.AppImageUrl[3], nodeSelector, 1, map[string]string{"app": "client"}, 81)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDClient, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
					break
				}
			}
			utils.Info("\n Client app deployed \n")

			//check name changed(communication happened)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))

			///delete server pod
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := deployment.Spec.Selector
					for key, value := range label.MatchLabels {
						podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler+utils.LabelSelector+key+"%3D"+value, "")
						Expect(err).To(BeNil())
						for _, pod := range podlist.Items {
							_, StatusCode := utils.DeletePods(ctx.Cfg.K8SMasterForKubeEdge + AppHandler + "/" + pod.Name)
							Expect(StatusCode).Should(Equal(http.StatusOK))
						}
					}
				}
			}

			//deployment will restart it check again pod is there
			podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler+utils.LabelSelector+"app"+"%3D"+"server", "")
			Expect(err).To(BeNil())
			utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge, podlist)

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep = utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)

			//check the name is changed of not
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))
		})

		FIt("E2E_SERVICE_EDGEMESH_4: Server deployment gets deleted: FAILURE", func() {
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			var servicelist metav1.ServiceList
			//Generate the random string and assign as a UID
			UIDServer = "pod-app-server" + utils.GetRandomString(5)
			depobj := utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed := utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
				}
			}
			utils.Info("\n Server app deployed \n")

			serviceName := "pod-app-server"
			// Deploy service over the server pod
			err = utils.ExposePodService(serviceName, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler, 80, intstr.FromInt(8000), metav1.ServiceTypeClusterIP)
			Expect(err).To(BeNil())
			err = utils.GetServices(&servicelist, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler)
			Expect(err).To(BeNil())

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep := utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))

			//deploy client deployment
			UIDClient = "pod-app-client" + utils.GetRandomString(5)
			depobj = utils.CreateDeployment(UIDClient, ctx.Cfg.AppImageUrl[3], nodeSelector, 1, map[string]string{"app": "client"}, 81)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDClient, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
					break
				}
			}
			utils.Info("\n Client app deployed \n")

			//check name changed(communication happened)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))

			//delete server deployment
			StatusCode := utils.DeleteDeployment(ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer)
			Expect(StatusCode).Should(Equal(http.StatusOK))

			//deploy again with same deployment name
			depobj = utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
				}
			}
			utils.Info("\n Server app deployed again\n")

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep = utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)

			//check the name is should not have been changed
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))
		})

		FIt("E2E_SERVICE_EDGEMESH_5: delete service : FAILURE", func() {
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			var servicelist metav1.ServiceList
			//Generate the random string and assign as a UID
			UIDServer = "pod-app-server" + utils.GetRandomString(5)
			depobj := utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed := utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
				}
			}
			utils.Info("\n Server app deployed \n")

			serviceName := "pod-app-server"
			// Deploy service over the server pod
			err = utils.ExposePodService(serviceName, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler, 80, intstr.FromInt(8000), metav1.ServiceTypeClusterIP)
			Expect(err).To(BeNil())
			err = utils.GetServices(&servicelist, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler)
			Expect(err).To(BeNil())

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep := utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))

			//deploy client deployment
			UIDClient = "pod-app-client" + utils.GetRandomString(5)
			depobj = utils.CreateDeployment(UIDClient, ctx.Cfg.AppImageUrl[3], nodeSelector, 1, map[string]string{"app": "client"}, 81)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDClient, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
					break
				}
			}
			utils.Info("\n Client app deployed \n")

			//check name changed(communication happened)
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Changed"))

			//delete service
			StatusCode := utils.DeleteSvc(ctx.Cfg.K8SMasterForKubeEdge + ServiceHandler + "/" + UIDServer)
			Expect(StatusCode).Should(Equal(http.StatusOK))

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep = utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)

			//change the name back to default again
			var jsonStr = []byte("Default")
			_, err = http.Post(ep, "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				panic(err)
			}

			//check the name should not have been changed
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))
		})

		FIt("E2E_SERVICE_EDGEMESH_6: create Loadbalancer service : FAILURE", func() {
			var podlist metav1.PodList
			var deploymentList v1.DeploymentList
			var servicelist metav1.ServiceList
			//Generate the random string and assign as a UID
			UIDServer = "pod-app-server" + utils.GetRandomString(5)
			depobj := utils.CreateDeployment(UIDServer, ctx.Cfg.AppImageUrl[2], nodeSelector, 1, map[string]string{"app": "server"}, 80)
			IsAppDeployed := utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDServer, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			err := utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDServer {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
				}
			}
			utils.Info("\n Server app deployed \n")

			serviceName := "pod-app-server"
			// Deploy service over the server pod
			err = utils.ExposePodService(serviceName, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler, 80, intstr.FromInt(8000), metav1.ServiceTypeLoadBalancer)
			Expect(err).To(BeNil())
			err = utils.GetServices(&servicelist, ctx.Cfg.K8SMasterForKubeEdge+ServiceHandler)
			Expect(err).To(BeNil())

			//deploy client deployment
			UIDClient = "pod-app-client" + utils.GetRandomString(5)
			depobj = utils.CreateDeployment(UIDClient, ctx.Cfg.AppImageUrl[3], nodeSelector, 1, map[string]string{"app": "client"}, 81)
			IsAppDeployed = utils.HandleRequest(http.MethodPost, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler, UIDClient, depobj)
			Expect(IsAppDeployed).Should(BeTrue())
			time.Sleep(time.Second * 1)
			err = utils.GetDeployments(&deploymentList, ctx.Cfg.K8SMasterForKubeEdge+DeploymentHandler+utils.LabelSelector+"app%3Dkubeedge")
			Expect(err).To(BeNil())
			for _, deployment := range deploymentList.Items {
				if deployment.Name == UIDClient {
					label := nodeName
					podlist, err = utils.GetPods(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, label)
					Expect(err).To(BeNil())
					utils.CheckPodRunningState(ctx.Cfg.K8SMasterForKubeEdge+AppHandler, podlist)
					break
				}
			}
			utils.Info("\n Client app deployed \n")

			// Check endpoints created
			utils.CheckEndpointCreated(ctx.Cfg.K8SMasterForKubeEdge+EndpointHandler, serviceName)

			_, ep := utils.GetServiceEndpoint(ctx.Cfg.K8SMasterForKubeEdge + EndpointHandler + "/" + serviceName)
			utils.Info("ep %s", ep)

			//check the name is should not have been changed
			Expect(utils.Getname(ep)).To(BeEquivalentTo("Default"))
		})
	})
})

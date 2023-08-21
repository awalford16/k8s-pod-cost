package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func AnnotatePod(clientset *kubernetes.Clientset, namespace, podName string, annotations map[string]string) error {
	retryErr := retry.OnError(
		retry.DefaultBackoff,
		func(err error) bool {
			return true
		},
		func() error {
			pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			pod.Annotations = annotations

			_, updateErr := clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
			return updateErr
		},
	)

	return retryErr
}

func main() {
	// kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to using KUBECONFIG path if available
		config, err = clientcmd.BuildConfigFromFlags("", "/home/adamw/.kube/config")
		if err != nil {
			panic(err.Error())
		}
	}
	// Get the metrics client
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	metricsCS, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default" // leave empty to get data from all namespaces
	podMetricsList, err := metricsCS.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, v := range podMetricsList.Items {
		pod := v.GetName()
		memory := v.Containers[0].Usage.Memory().Value() / (1024 * 1024)

		annotations := map[string]string{
			"controller.kubernetes.io/pod-deletion-cost": strconv.Itoa(int(memory * -1)),
		}

		err = AnnotatePod(cs, namespace, pod, annotations)
		if err != nil {
			log.Fatalf("Error annotating pod: %v", err)
		}

		fmt.Println("Pod annotated successfully!")
	}
}
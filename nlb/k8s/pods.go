package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var Client *kubernetes.Clientset

func NewClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil

}
func GetPodDetails(podName string) (string, error) {

	pod, err := Client.CoreV1().Pods("default").Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pod.Status.PodIP, nil
}

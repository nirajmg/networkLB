package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Client *kubernetes.Clientset

func NewClient() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
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

func ListPod() ([]string, error) {

	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=server",
	}

	podList, err := Client.CoreV1().Pods("default").List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", podList.Items[0].Status.PodIP)
	podIPs := make([]string, 0)
	for _, v := range podList.Items {
		podIPs = append(podIPs, v.Status.PodIP)
	}

	return podIPs, nil
}

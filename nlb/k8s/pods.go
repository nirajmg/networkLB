package k8s

import (
	"context"
	"flag"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Client *kubernetes.Clientset

type PodDetails struct {
	IP     string
	Memory float64
}

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

func ListPod() ([]*PodDetails, error) {

	listOptions := metav1.ListOptions{
		LabelSelector: "app.server/filter=server",
	}

	podList, err := Client.CoreV1().Pods("default").List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}

	podIPs := make([]*PodDetails, 0)
	for _, v := range podList.Items {
		if v.Status.Phase == "Running" {
			p := &PodDetails{
				IP:     v.Status.PodIP,
				Memory: v.Spec.Containers[0].Resources.Requests.Memory().AsApproximateFloat64(),
			}
			podIPs = append(podIPs, p)
		}
	}

	return podIPs, nil
}

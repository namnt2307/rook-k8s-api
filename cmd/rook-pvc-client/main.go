package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"      // k8s api
	"k8s.io/client-go/tools/clientcmd" // load config
)

func initClient() *kubernetes.Clientset {
	// import kubeconfig
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println("Using kubeconfig ", kubeconfig)

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	//Load clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

func listPods(clientset *kubernetes.Clientset, namespace string) {
	apiVersion := clientset.CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "",
		FieldSelector: "",
	}
	pods, err := apiVersion.Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal(err)
	}
	file, _ := json.MarshalIndent(pods, "", " ")

	_ = ioutil.WriteFile("pods.json", file, 0644)
	for _, pod := range pods.Items {
		fmt.Printf("%v \n", pod.Name)
	}
}

func listPVC(clientset *kubernetes.Clientset, namespace string) {
	apiVersion := clientset.CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "",
		FieldSelector: "",
	}
	pvcs, err := apiVersion.PersistentVolumeClaims(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal(err)
	}
	file, _ := json.MarshalIndent(pvcs, "", " ")

	_ = ioutil.WriteFile("pvcs.json", file, 0644)
	fmt.Printf("%-16s %-64s \n", "NAMESPACE", "PVC")
	for _, pvc := range pvcs.Items {
		fmt.Printf("%-16s %-64s  \n", namespace, pvc.Name)
		// os.Mkdir(pvc.Name, 0777)
		// os.Chmod(pvc.Name, 0777)
	}
}

func main() {
	clientset := initClient()
	// listPods(clientset, "dev")
	listPVC(clientset, "kube-logging")
}

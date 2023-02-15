package currentversions

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// example command 'kubectl get deployments,statefulsets,daemonsets -A -o go-template --template '{{range .items}}{{index .metadata.labels "app.kubernetes.io/name"}} {{index .metadata.labels "app.kubernetes.io/version"}} {{"\n"}}{{end}}''

func GetClientSet(cluster bool, kubeconfig string) *kubernetes.Clientset {
	if cluster == false {

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return clientset
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

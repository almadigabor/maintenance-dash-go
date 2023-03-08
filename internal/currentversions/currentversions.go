package currentversions

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/almadigabor/maintenance-dash-go/internal/data"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// example command 'kubectl get deployments,statefulsets,daemonsets -A -o go-template --template '{{range .items}}{{index .metadata.labels "app.kubernetes.io/name"}} {{index .metadata.labels "app.kubernetes.io/version"}} {{"\n"}}{{end}}”
func NewClientSet(cluster bool, kubeconfig string) *kubernetes.Clientset {
	if !cluster {

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		cs, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		return cs
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return cs
}

func GetSvcsToScan(ctx context.Context, clientSet *kubernetes.Clientset) {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"maintenance/scan": "true"}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	services, _ := clientSet.CoreV1().Services("").List(ctx, listOptions)
	for _, svc := range services.Items {
		fmt.Printf("%v chart: %v", svc.Name, svc.ObjectMeta.Labels["helm.sh/chart"])
	}
}

func AddK8sNodeVersionInfo(ctx context.Context, clientSet *kubernetes.Clientset, appsVersionInfo *[]data.AppVersionInfo) {
	nodes, _ := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	for _, node := range nodes.Items {
		v, err := semver.NewVersion(node.Status.NodeInfo.KubeletVersion)
		if err != nil {
			log.Warnf("Unable to parse kubernetes node version: %v", node.Status.NodeInfo.KubeletVersion)
		}
		*appsVersionInfo = append(*appsVersionInfo, data.AppVersionInfo{
			AppName:        "kubernetes",
			CurrentVersion: v,
		})
	}
}

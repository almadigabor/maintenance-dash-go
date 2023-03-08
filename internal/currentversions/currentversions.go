package currentversions

import (
	"context"
	"fmt"

	"github.com/almadigabor/maintenance-dash-go/internal/data"
	"github.com/almadigabor/maintenance-dash-go/internal/parseversion"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Returns the versions of nodes and labeled services
func GetCurrentVersions(ctx context.Context, cluster bool, kubeconfig string) []*data.AppVersionInfo {
	clientSet := newClientSet(cluster, kubeconfig)
	services := getSvcsToScan(ctx, clientSet)
	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("Unable to get node info: %v", err)
	}

	var result []*data.AppVersionInfo

	for _, node := range nodes.Items {
		semverVersion := parseversion.ToSemver(node.Status.NodeInfo.KubeletVersion)
		result = append(result, &data.AppVersionInfo{
			CurrentVersion:  semverVersion,
			NewReleasesName: node.ObjectMeta.Annotations["maintenance/releasename"],
		})
	}

	for _, service := range services.Items {
		semverVersion := parseversion.ToSemver(service.ObjectMeta.Labels["app.kubernetes.io/version"])
		result = append(result, &data.AppVersionInfo{
			CurrentVersion:  semverVersion,
			NewReleasesName: service.ObjectMeta.Annotations["maintenance/releasename"],
		})
	}

	return result
}

// Initializes new ClientSet either based on kubeconfig or in-cluster
func newClientSet(cluster bool, kubeconfig string) *kubernetes.Clientset {
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

// Get services annotated with maintenance/scan=true
func getSvcsToScan(ctx context.Context, clientSet *kubernetes.Clientset) *corev1.ServiceList {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"maintenance/scan": "true"}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	services, err := clientSet.CoreV1().Services("").List(ctx, listOptions)

	if err != nil {
		log.Errorf("Unable to get services to scan: %v", err)
	}
	for _, svc := range services.Items {
		fmt.Printf("%v chart: %v", svc.Name, svc.ObjectMeta.Labels["helm.sh/chart"])
	}

	return services
}

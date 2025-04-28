package clusterclient

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// BuildClusterClient creates a kubernetes client for the destination cluster.
func BuildClusterClient(server, token string) (kubernetes.Interface, error) {
	config := &rest.Config{
		Host:        server,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	return kubernetes.NewForConfig(config)
}

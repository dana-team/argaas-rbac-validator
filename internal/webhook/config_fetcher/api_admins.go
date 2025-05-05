package config_fetcher

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FetchAPIAdmins reads the ConfigMap api-admins inside Application namespace.
func FetchAPIAdmins(ctx context.Context, k8sClient client.Client, namespace string) ([]string, error) {
	var cm corev1.ConfigMap
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: "api-admins"}, &cm)
	if err != nil {
		return nil, fmt.Errorf("failed to get api-admins ConfigMap: %w", err)
	}

	usersData, ok := cm.Data["users"]
	if !ok {
		return nil, fmt.Errorf("api-admins ConfigMap missing 'users' key")
	}

	users := strings.Split(strings.TrimSpace(usersData), "\n")
	return users, nil
}

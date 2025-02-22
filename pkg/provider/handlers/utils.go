package handlers

import (
	"context"
	"net/http"
	"path"

	"github.com/containerd/containerd"

	"github.com/openfaas/faasd/pkg"
	faasd "github.com/openfaas/faasd/pkg"
)

func getRequestNamespace(namespace string) string {

	if len(namespace) > 0 {
		return namespace
	}
	return faasd.DefaultFunctionNamespace
}

func readNamespaceFromQuery(r *http.Request) string {
	q := r.URL.Query()
	return q.Get("namespace")
}

func getNamespaceSecretMountPath(userSecretPath string, namespace string) string {
	return path.Join(userSecretPath, namespace)
}

// validNamespace indicates whether the namespace is eligable to be
// used for OpenFaaS functions.
func validNamespace(client *containerd.Client, namespace string) (bool, error) {
	if namespace == faasd.DefaultFunctionNamespace {
		return true, nil
	}

	store := client.NamespaceService()
	labels, err := store.Labels(context.Background(), namespace)
	if err != nil {
		return false, err
	}

	if value, found := labels[pkg.NamespaceLabel]; found && value == "true" {
		return true, nil
	}

	return false, nil
}

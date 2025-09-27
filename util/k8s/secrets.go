package k8s

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetSecrets retrieves a secret value.
func GetSecrets(ctx context.Context, clientset kubernetes.Interface, namespace, name, key string) ([]byte, error) {
	secretsIf := clientset.CoreV1().Secrets(namespace)

	secret, err := secretsIf.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return []byte{}, errors.Wrapf(err, "secret %s not found", name)
	}
	val, ok := secret.Data[key]
	if !ok {
		return []byte{}, errors.Wrapf(err, "secret %s does not have the key %s", name, key)
	}
	return val, nil
}

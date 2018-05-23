package client

import (
	yaml "gopkg.in/yaml.v2"

	"github.com/dchest/uniuri"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
)

const (
	tillerHost = "tiller-deploy.kube-system.svc.cluster.local:44134"
	chartPath  = "/minio-chart.tgz"
)

// Install creates a new minio chart release
func Install(releaseName, namespace string) error {
	vals, err := yaml.Marshal(map[string]interface{}{
		"accessKey": uniuri.NewLen(20),
		"secretKey": uniuri.NewLen(40),
		"minioConfig": map[string]interface{}{
			"webhook": map[string]interface{}{
				"enable": true,
				"endpoint": "http://service-adapter.fission/minio/" + namespace,
			},
		},
	})
	if err != nil {
		return err
	}
	helmClient := helm.NewClient(helm.Host(tillerHost))
	_, err = helmClient.InstallRelease(chartPath, namespace, helm.ReleaseName(releaseName), helm.ValueOverrides(vals))
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a minio chart release
func Delete(releaseName string) error {
	helmClient := helm.NewClient(helm.Host(tillerHost))
	if _, err := helmClient.DeleteRelease(releaseName); err != nil {
		return err
	}
	return nil
}

// GetPassword returns the minio access/secret key for a chart release
func GetPassword(releaseName, namespace string) (string, string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", "", err
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(releaseName + "-minio-user")
	if err != nil {
		return "", "", err
	}
	return string(secret.Data["accesskey"]), string(secret.Data["secretkey"]), nil
}

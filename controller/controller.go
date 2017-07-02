package controller

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/controller"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"github.com/yqf3139/minio-broker/client"
)

type errNoSuchInstance struct {
	instanceID string
}

func (e errNoSuchInstance) Error() string {
	return fmt.Sprintf("no such instance with ID %s", e.instanceID)
}

type minioController struct {
}

// CreateController creates an instance of a User Provided service broker controller.
func CreateController() controller.Controller {
	return &minioController{}
}

func (c *minioController) Catalog() (*brokerapi.Catalog, error) {
	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				Name:        "minio",
				ID:          "10faf2e0-5e61-11e7-9c17-fb844ec31790",
				Description: "minio database",
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          "18413dde-5e61-11e7-93e3-179ef79e4573",
						Description: "minio database",
						Free:        true,
					},
				},
				Bindable: true,
			},
		},
	}, nil
}

func (c *minioController) CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (*brokerapi.CreateServiceInstanceResponse, error) {
	if err := client.Install(releaseName(id), id); err != nil {
		return nil, err
	}
	glog.Infof("Created minio Service Instance:\n%v\n", id)
	return &brokerapi.CreateServiceInstanceResponse{}, nil
}

func (c *minioController) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *minioController) RemoveServiceInstance(id string) (*brokerapi.DeleteServiceInstanceResponse, error) {
	if err := client.Delete(releaseName(id)); err != nil {
		return nil, err
	}
	return &brokerapi.DeleteServiceInstanceResponse{}, nil
}

func (c *minioController) Bind(instanceID, bindingID string, req *brokerapi.BindingRequest) (*brokerapi.CreateServiceBindingResponse, error) {
	host := releaseName(instanceID) + "-minio-svc." + instanceID + ".svc.cluster.local"
	accessKey, secretKey, err := client.GetPassword(releaseName(instanceID), instanceID)
	if err != nil {
		return nil, err
	}
	return &brokerapi.CreateServiceBindingResponse{
		Credentials: brokerapi.Credential{
			"host":           host,
			"port":           "9000",
			"keys.access":    accessKey,
			"keys.secret":    secretKey,
		},
	}, nil
}

func (c *minioController) UnBind(instanceID string, bindingID string) error {
	// Since we don't persist the binding, there's nothing to do here.
	return nil
}

func releaseName(id string) string {
	return "i-" + id
}

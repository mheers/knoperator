package deployment

import (
	"encoding/json"
	"sync"

	"github.com/mheers/knoperator/api/deployment/models"
	"github.com/mheers/knoperator/config"
	"github.com/mheers/knoperator/helpers"
	"github.com/mheers/knoperator/integrations/deployment"
	"github.com/mheers/knoperator/services"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// DeploymentService is a struct that contains the connection to the immudb
type DeploymentService struct {
	config                *config.Config
	once                  sync.Once
	deploymentIntegration *deployment.DeploymentIntegration
}

// NewDeploymentService creates a new deploymentservice
func NewDeploymentService(
	config *config.Config,
	deploymentIntegration *deployment.DeploymentIntegration,
) (*DeploymentService, error) {
	deploymentservice := &DeploymentService{
		config:                config,
		deploymentIntegration: deploymentIntegration,
	}
	return deploymentservice, nil
}

// Start starts the deploymentservice
func (ds *DeploymentService) Start() error {
	var err error
	ds.once.Do(func() {
		err = ds.start()
	})
	return err
}

func (ds *DeploymentService) start() error {
	logrus.Info("knoperator: start: deploymentservice.start()")

	nc := services.MQClient().Connection

	go func() {
		err := ds.deploymentIntegration.WatchJobs(nc)
		if err != nil {
			logrus.Errorf("knoperator: start: deploymentservice.start: WatchJobs: %s", err.Error())
		}
	}()

	nc.Subscribe("knoperator.pods.get", func(m *nats.Msg) {
		pods, err := ds.deploymentIntegration.GetPods()
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		responseJSON, _ := json.Marshal(pods)
		msg := &nats.Msg{
			Data: responseJSON,
		}
		m.RespondMsg(msg)
	})

	nc.Subscribe("knoperator.pods.delete", func(m *nats.Msg) {
		name := string(m.Data)
		err := ds.deploymentIntegration.DeletePod(name)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	nc.Subscribe("knoperator.deployments.get", func(m *nats.Msg) {
		pods, err := ds.deploymentIntegration.GetDeployments()
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		responseJSON, _ := json.Marshal(pods)
		msg := &nats.Msg{
			Data: responseJSON,
		}
		m.RespondMsg(msg)
	})

	nc.Subscribe("knoperator.deployments.delete", func(m *nats.Msg) {
		name := string(m.Data)
		err := ds.deploymentIntegration.DeleteDeployment(name)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	nc.Subscribe("knoperator.deployments.create", func(m *nats.Msg) {
		var request models.DeploymentCreateRequest
		err := json.Unmarshal(m.Data, &request)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}

		err = ds.deploymentIntegration.CreateDeployment(request.Name, request.Image, request.Command, request.Args, request.Env)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	nc.Subscribe("knoperator.deployments.update", func(m *nats.Msg) {
		type DeploymentUpdateRequest struct {
			Name    string
			Image   string
			Command []string
			Args    []string
			Env     map[string]string
		}
		var request DeploymentUpdateRequest
		err := json.Unmarshal(m.Data, &request)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}

		err = ds.deploymentIntegration.UpdateDeployment(request.Name, request.Image, request.Command, request.Args, request.Env)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	nc.Subscribe("knoperator.deployments.scale", func(m *nats.Msg) {
		type DeploymentScaleRequest struct {
			Name      string
			NReplicas int
		}
		var request DeploymentScaleRequest
		err := json.Unmarshal(m.Data, &request)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}

		err = ds.deploymentIntegration.ScaleDeployment(request.Name, request.NReplicas)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	nc.Subscribe("knoperator.jobs.get", func(m *nats.Msg) {
		pods, err := ds.deploymentIntegration.GetJobs()
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		responseJSON, _ := json.Marshal(pods)
		msg := &nats.Msg{
			Data: responseJSON,
		}
		m.RespondMsg(msg)
	})

	nc.Subscribe("knoperator.jobs.delete", func(m *nats.Msg) {
		name := string(m.Data)
		err := ds.deploymentIntegration.DeleteJob(name)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	nc.Subscribe("knoperator.jobs.create", func(m *nats.Msg) {
		var request models.JobCreateRequest
		err := json.Unmarshal(m.Data, &request)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}

		err = ds.deploymentIntegration.CreateJob(request.Name, request.Image, request.Command, request.Args, request.Env, request.MountPoints)
		if err != nil {
			helpers.HandleMQError(m, err)
			return
		}
		helpers.HandleMQOK(m)
	})

	return nil
}

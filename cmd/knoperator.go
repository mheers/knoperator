package cmd

import (
	"errors"

	deploymentService "github.com/mheers/knoperator/api/deployment"
	"github.com/mheers/knoperator/config"
	"github.com/mheers/knoperator/helpers"
	deploymentIntegration "github.com/mheers/knoperator/integrations/deployment"
	"github.com/mheers/knoperator/k8sclient"
	mqmodels "github.com/mheers/knoperator/mqclient/models"
	"github.com/mheers/knoperator/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	knoperatorCmd = &cobra.Command{
		Use:   "start",
		Short: "starts the knoperator",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := config.GetConfig(false)

			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			return startServer(config)
		},
	}
)

func startServer(config *config.Config) error {

	// create the k8s client connection
	err := initK8s(config)
	if err != nil {
		return err
	}

	// create the mq client connection
	err = initMQ(config)
	if err != nil {
		return err
	}

	// Start Deployment Integration
	logrus.Info("Starting Deployment Integration")
	deploymentI, err := deploymentIntegration.NewAPI(config, services.K8sClient)
	if err != nil {
		return err
	}

	// Create the deployment service
	logrus.Info("Starting deployment service")
	deploymentS, err := deploymentService.NewDeploymentService(config, deploymentI)
	if err != nil {
		return err
	}
	err = deploymentS.Start()
	if err != nil {
		return err
	}

	// loop forever
	select {}
}

func initK8s(config *config.Config) error {
	logrus.Debug("enabling feature 'K8s'")
	k8sClient, err := k8sclient.Init(config)
	if err != nil {
		return err
	}
	services.K8sClient = k8sClient

	return nil
}

func initMQ(config *config.Config) error {
	logrus.Debug("enabling feature 'Message Queue'")
	if config.MQURI == "" {
		return errors.New("no MQURI found")
	}

	if config.MQCredsPath == "" && (config.MQUSeed == "" || config.MQJWT == "") {
		return errors.New("no MQCredsPath or MQUSeed/MQJWT found")
	}

	mqClient, err := mqmodels.NewMQClient(config)
	if err != nil {
		return err
	}
	services.SetMQClient(mqClient)

	return nil
}

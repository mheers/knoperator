package services

import mqmodels "github.com/mheers/knoperator/mqclient/models"

// iMQClient holds the connection to the Message Queue
var iMQClient *mqmodels.MQClient

func SetMQClient(client *mqmodels.MQClient) {
	iMQClient = client
}

func MQClient() *mqmodels.MQClient {
	if iMQClient == nil {
		panic("MQClient is not initialized")
	}
	return iMQClient
}

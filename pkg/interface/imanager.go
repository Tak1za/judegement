package imanager

import "github.com/Tak1za/judgement/pkg/models"

type IManager interface {
	BroadcastMessage(message models.Message)
	UnregisterSubscriber(client *models.Client)
	RegisterSubscriber(client *models.Client)
	Send(message []byte, ignoreClient *models.Client)
	Start()
}

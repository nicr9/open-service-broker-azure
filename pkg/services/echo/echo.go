package echo

import (
	"log"

	"time"

	"context"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

type module struct{}

// New returns a new instance of a type that fulfills the service.Module
// and provides an example of how such a module is implemented
func New() service.Module {
	return &module{}
}

func (m *module) GetName() string {
	return "example"
}

func (m *module) ValidateProvisioningParameters(
	provisioningParameters interface{},
) error {
	params := provisioningParameters.(*echoProvisioningParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) GetProvisioner() (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("generateMessageId", m.generateProvisioningMessageID),
		service.NewProvisioningStep("pause", m.pause),
		service.NewProvisioningStep("logMessage", m.logProvisioningMessage),
	)
}

func (m *module) generateProvisioningMessageID(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	log.Println("Executing generateProvisioningMessageID...")
	pc := provisioningContext.(*echoProvisioningContext)
	pc.MessageID = uuid.NewV4().String()
	return pc, nil
}

func (m *module) pause(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	log.Println("Executing pause...")
	select {
	case <-time.NewTimer(time.Minute).C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return provisioningContext, nil
}

func (m *module) logProvisioningMessage(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	log.Println("Executing logProvisioningMessage...")
	pc := provisioningContext.(*echoProvisioningContext)
	params := provisioningParameters.(*echoProvisioningParameters)
	pc.Message = params.Message
	log.Printf("Provisioning %s: %s", pc.MessageID, params.Message)
	return pc, nil
}

func (m *module) ValidateBindingParameters(
	bindingParameters interface{},
) error {
	params := bindingParameters.(*echoBindingParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) Bind(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, error) {
	log.Println("Executing Bind...")
	params := bindingParameters.(*echoBindingParameters)
	bc := &echoBindingContext{
		MessageID: uuid.NewV4().String(),
		Message:   params.Message,
	}
	log.Printf("Binding %s: %s", bc.MessageID, params.Message)
	return bc, nil
}

func (m *module) Unbind(
	provisioningContext interface{},
	bindingContext interface{},
) error {
	log.Println("Executing Unbind...")
	bc := bindingContext.(*echoBindingContext)
	log.Printf("Unbinding %s: %s", bc.MessageID, bc.Message)
	return nil
}

func (m *module) GetDeprovisioner() (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep(
			"logDeprovisioningMessage",
			m.logDeprovisioningMessage,
		),
	)
}

func (m *module) logDeprovisioningMessage(
	ctx context.Context,
	provisioningContext interface{},
) (interface{}, error) {
	log.Println("Executing logDeprovisioningMessage...")
	pc := provisioningContext.(*echoProvisioningContext)
	log.Printf("Deprovisioning %s: %s", pc.MessageID, pc.Message)
	return pc, nil
}

// services/startup_service.go
package services

import (
	"context"
	"log"

	"github.com/google/wire"
)

type StartupServiceInterface interface {
	Initialize(ctx context.Context) error
}

type StartupServiceStruct struct {
	roleService ReWriteRoleServiceInterface
	// Add other initialization services here
}

func NewStartupService(
	roleService ReWriteRoleServiceInterface,
) StartupServiceInterface {
	return &StartupServiceStruct{
		roleService: roleService,
	}
}

var _ StartupServiceInterface = (*StartupServiceStruct)(nil)

var StartupServiceSet = wire.NewSet(
	NewStartupService,
)

func (s *StartupServiceStruct) Initialize(ctx context.Context) error {
	log.Println("Starting application initialization...")
	
	// Initialize roles
	if err := s.roleService.InitializeReWriteRole(ctx); err != nil {
		return err
	}
	
	// Add other initialization tasks here
	
	log.Println("Application initialization completed successfully")
	return nil
}
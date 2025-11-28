package services

import (
	"context"
	models "ganjineh-auth/internal/models/entities"
	"ganjineh-auth/internal/repositories"
	"ganjineh-auth/internal/repositories/db"

	"github.com/google/wire"
)

type ReWriteRoleServiceInterface interface {
	InitializeReWriteRole(ctx context.Context) error
}


type RewriteRoleServiceStruct struct {
	iamRepo		*db.Queries 
	roleRepo 	repositories.RedisPermissionRepositoryInterface
}

func NewRewriteRoleService(
	postRepo *db.Queries,
	roleRepo repositories.RedisPermissionRepositoryInterface,
) ReWriteRoleServiceInterface {
    return &RewriteRoleServiceStruct{
   		iamRepo: postRepo,
    	roleRepo: roleRepo,
    }
}

var RewriteRoleServiceSet = wire.NewSet(
	NewRewriteRoleService,
	// wire.Bind(new(OTPServiceInterface), new(*OTPServiceStruct)),
)

var _ ReWriteRoleServiceInterface = (*RewriteRoleServiceStruct)(nil)

func (r *RewriteRoleServiceStruct) InitializeReWriteRole(ctx context.Context) error{
	if err:= r.roleRepo.ClearAllRoles(ctx); err!=nil{
		return err
	}
	roles, err := r.iamRepo.GetAllRoles(ctx)
	if err != nil{
		return err
	}
	for _, role := range roles{
		roleM := models.Role{
			Name: role.Name,
			PermissionCodes: role.PermissionCodes,
			IsSystemRole: role.IsSystemRole,
		}
		if err:=r.roleRepo.StorePermission(ctx, &roleM); err != nil{
			return err
		}
	}
	return  nil
}
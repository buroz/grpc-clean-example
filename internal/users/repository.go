package users

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/buroz/grpc-clean-example/pkg/domain"
)

type UsersRepository interface {
	Create(context.Context, *domain.User) (driver.DocumentID, error)
	FindByEmail(context.Context, string) ([]domain.User, error)
	SetNewTokens(context.Context, *domain.User) error
	SetNewPasswordResetRequest(context.Context, *domain.User) error
	FindByAccessToken(context.Context, string, bool) ([]domain.User, error)
	FindByConfirmationToken(context.Context, string) ([]domain.User, error)
	ConfirmUserWithToken(context.Context, string, string) (bool, error)
	SetNewPassword(context.Context, *domain.User) ([]domain.User, error)
	FindByTokens(context.Context, string, string, bool) ([]domain.User, error)
	FindById(context.Context, string) (domain.User, error)
	// FindExpiredToken(context.Context, string, string) ([]domain.User, error)
}

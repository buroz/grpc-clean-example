package interceptors

import (
	"context"
	"errors"

	"github.com/buroz/grpc-clean-example/internal/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var protectedMethods = []string{
	"/users.UsersService/Profile",
}

type Interceptors struct {
	usersService users.UsersService
}

func NewServiceInterceptors(us users.UsersService) *Interceptors {
	return &Interceptors{
		usersService: us,
	}
}

func (in *Interceptors) AuthorizationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		found := false

		for _, method := range protectedMethods {
			if method == info.FullMethod {
				found = true
				break
			}
		}

		if found {
			meta, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, errors.New("could not grab metadata from context")
			}

			var values []string
			var token string

			if ok {
				values = meta.Get("authorization")
			}

			if len(values) > 0 && len(values[0]) > 7 {
				token = values[0][7:]

				foundUser, err := in.usersService.ValidateWithToken(ctx, token)
				if err != nil {
					return nil, errors.New("token parse error")
				}

				if foundUser != nil {
					meta.Set("user", foundUser.Id.Key())
					newCtx := metadata.NewIncomingContext(ctx, meta)
					return handler(newCtx, req)
				}

				return nil, errors.New("user not found")
			}

			return nil, errors.New("no token provided")
		}

		return handler(ctx, req)
	}
}

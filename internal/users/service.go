package users

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/buroz/grpc-clean-example/pkg/domain"
	"github.com/buroz/grpc-clean-example/pkg/exceptions"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersService interface {
	Register(context.Context, *domain.UserRegisterDto) (*domain.User, *exceptions.ErrorResponse)
	Login(context.Context, *domain.UserLoginDto) (*domain.User, *exceptions.ErrorResponse)
	Profile(context.Context, string) (*domain.User, *exceptions.ErrorResponse)
	ConfirmRegistration(context.Context, string) (bool, *exceptions.ErrorResponse)
	PasswordResetRequest(context.Context, *domain.UserPasswordResetRequestDto) (bool, *exceptions.ErrorResponse)
	SetNewPassword(context.Context, string, *domain.UserSetNewPasswordDto) (bool, *exceptions.ErrorResponse)

	ValidateWithToken(context.Context, string) (*domain.User, *exceptions.ErrorResponse)
	RefreshToken(context.Context, string, string) (*domain.User, *exceptions.ErrorResponse)

	hashPassword(password string) (string, error)
	checkPasswordHash(string, string) bool
}

type usersService struct {
	ur UsersRepository
}

func NewUsersService(ur UsersRepository) UsersService {
	return &usersService{
		ur,
	}
}

func (us *usersService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (us *usersService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (us *usersService) Register(ctx context.Context, dto *domain.UserRegisterDto) (*domain.User, *exceptions.ErrorResponse) {
	user := new(domain.User)

	user.Email = dto.Email
	user.FirstName = dto.FirstName
	user.LastName = dto.LastName

	hashedPassword, err := us.hashPassword(dto.Password)
	if err != nil {
		return user, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
	}

	user.Password = hashedPassword
	user.ConfirmationToken = uuid.New().String()
	user.CreatedAt = time.Now()
	// user.Subscriptions = append(user.Subscriptions, dto.Subscription)

	id, err := us.ur.Create(ctx, user)
	if err != nil {
		if driver.IsArangoErrorWithErrorNum(err, driver.ErrArangoUniqueConstraintViolated) {
			return nil, exceptions.NewDuplicateError([]string{"email"})
		}

		return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
	}

	user.Id = id

	return user, nil
}

func (us *usersService) Login(ctx context.Context, dto *domain.UserLoginDto) (*domain.User, *exceptions.ErrorResponse) {
	foundUsers, err := us.ur.FindByEmail(ctx, dto.Email)
	if err != nil {
		return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
	}

	if len(foundUsers) > 0 {
		user := &foundUsers[0]

		if user.Password != "" && us.checkPasswordHash(dto.Password, user.Password) {
			user.LatestToken = uuid.New().String()
			user.RefreshToken = uuid.New().String()

			err = us.ur.SetNewTokens(ctx, user)
			if err != nil {
				return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
			}

			return user, nil
		}

		return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeWrongPassword}
	}

	return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeNotFound}
}

func (us *usersService) PasswordResetRequest(ctx context.Context, dto *domain.UserPasswordResetRequestDto) (bool, *exceptions.ErrorResponse) {
	foundUsers, err := us.ur.FindByEmail(ctx, dto.Email)
	if err != nil {
		return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
	}

	if len(foundUsers) > 0 {
		user := &foundUsers[0]

		user.PasswordResetToken = uuid.New().String()

		err = us.ur.SetNewPasswordResetRequest(ctx, user)
		if err != nil {
			return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
		}

		return len(foundUsers) > 0, nil
	}

	return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeNotFound}
}

func (us *usersService) ConfirmRegistration(ctx context.Context, token string) (bool, *exceptions.ErrorResponse) {
	foundUsers, err := us.ur.FindByConfirmationToken(ctx, token)
	if err != nil {
		return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
	}

	if len(foundUsers) > 0 {
		user := &foundUsers[0]

		confirmed, err := us.ur.ConfirmUserWithToken(ctx, user.Id.Key(), token)
		if err != nil {
			return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
		}

		return confirmed, nil
	}

	return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeNotFound}
}

func (us *usersService) SetNewPassword(ctx context.Context, token string, dto *domain.UserSetNewPasswordDto) (bool, *exceptions.ErrorResponse) {
	if dto.Password == dto.PasswordConfirm {
		user := new(domain.User)

		hashed, err := us.hashPassword(dto.Password)
		if err != nil {
			return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
		}

		user.Password = hashed
		user.PasswordResetToken = token

		foundUsers, err := us.ur.SetNewPassword(ctx, user)
		if err != nil {
			return false, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
		}

		if len(foundUsers) > 0 {
			return true, nil
		}

		return false, nil
	}

	return false, exceptions.NewDuplicateError([]string{"password_confirm"})
}

func (us *usersService) RefreshToken(ctx context.Context, token, refreshToken string) (*domain.User, *exceptions.ErrorResponse) {
	foundUsers, err := us.ur.FindByTokens(ctx, token, refreshToken, false)
	if err != nil {
		return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeNotFound}
	}

	if len(foundUsers) > 0 {
		user := &foundUsers[0]

		user.LatestToken = uuid.New().String()
		user.RefreshToken = uuid.New().String()

		err = us.ur.SetNewTokens(ctx, user)
		if err != nil {
			return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
		}

		return user, nil
	}

	return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeNotFound}
}

func (us *usersService) ValidateWithToken(ctx context.Context, token string) (*domain.User, *exceptions.ErrorResponse) {
	doc, err := us.ur.FindByAccessToken(ctx, token, true)
	if err != nil {
		return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeInternalServer}
	}

	if len(doc) > 0 {
		user := &doc[0]
		return user, nil
	}

	return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeUnauthorized}
}

func (us *usersService) Profile(ctx context.Context, id string) (*domain.User, *exceptions.ErrorResponse) {
	doc, err := us.ur.FindById(ctx, id)
	if err != nil {
		return nil, &exceptions.ErrorResponse{Code: exceptions.ErrCodeNotFound}
	}

	return &doc, nil
}

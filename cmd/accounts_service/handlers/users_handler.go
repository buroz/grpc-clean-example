package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/buroz/grpc-clean-example/internal/common"
	"github.com/buroz/grpc-clean-example/internal/users"
	AMQP "github.com/buroz/grpc-clean-example/pkg/amqp"
	"github.com/buroz/grpc-clean-example/pkg/config"
	"github.com/buroz/grpc-clean-example/pkg/domain"
	"github.com/buroz/grpc-clean-example/pkg/validator"
	"github.com/streadway/amqp"
)

var (
	emailsExchangeName = "EMAILS"
	// emailsQueueName    = "EMAIL_SEND"
)

type UsersHandler struct {
	usersService   users.UsersService
	emailsExchange *AMQP.AmqpClient
}

// TARTIŞMALI: amqp için buraya config çekmek mi gerekli?
func (handler *UsersHandler) Inject(usersService users.UsersService, amqpConfig *config.AmqpConfig) error {
	amqpClient := AMQP.NewAmqpClient(amqpConfig)

	err := amqpClient.Connect()
	if err != nil {
		return err
	}

	err = amqpClient.Channel.ExchangeDeclare(
		emailsExchangeName, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	handler.emailsExchange = &amqpClient

	handler.usersService = usersService

	return nil
}

func (handler *UsersHandler) Register(ctx context.Context, dto *users.RegisterRequest) (*users.TokenResponse, error) {
	registerDto := domain.UserRegisterDto{}

	registerDto.Email = dto.GetEmail()
	registerDto.FirstName = dto.GetFirstName()
	registerDto.LastName = dto.GetLastName()
	registerDto.Password = dto.GetPassword()

	/*
		// TODO: ...
		registerDto.Subscription = domain.Subscription{
			Status:    true,
			AppType:   domain.ServiceTypeSupplierProgram,
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(1, 0, 0),
		}
	*/

	resp := &users.TokenResponse{
		Success: false,
	}

	err := validator.ValidateStruct(&registerDto)
	if err != nil {
		resp.Errors = common.ParseErrors(err)
		return resp, nil
	}

	newUser, err := handler.usersService.Register(ctx, &registerDto)
	if err != nil {
		resp.Errors = common.ParseErrors(err)
		return resp, nil
	}

	emailMsg, jsonErr := json.Marshal(map[string]string{
		"to":      newUser.Email,
		"body":    newUser.ConfirmationToken,
		"subject": "Welcome to Example!",
	})
	if jsonErr != nil {
		// TODO: parse error
		fmt.Println("naber")
	}

	emailError := handler.emailsExchange.Channel.Publish(
		emailsExchangeName, // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        emailMsg,
		},
	)
	if emailError != nil {
		// TODO: parse error
		// Hatta user'ı silmemiz gerekli mi acaba?
		log.Println(emailError)
	}

	resp.Success = true

	return resp, nil
}

// TODO !!!
func (handler *UsersHandler) Login(ctx context.Context, dto *users.LoginRequest) (*users.TokenResponse, error) {
	loginDto := domain.UserLoginDto{}

	loginDto.Email = dto.GetEmail()
	loginDto.Password = dto.GetPassword()

	resp := &users.TokenResponse{
		Result:  &users.TokenResult{},
		Success: false,
	}

	err := validator.ValidateStruct(&loginDto)
	if err != nil {
		resp.Errors = common.ParseErrors(err)
		return resp, nil
	}

	loggedUser, err := handler.usersService.Login(ctx, &loginDto)
	if err != nil {
		resp.Errors = common.ParseErrors(err)
		return resp, nil
	}

	resp.Success = true
	resp.Result.Token = loggedUser.LatestToken
	resp.Result.RefreshToken = loggedUser.RefreshToken

	return resp, nil
}

func (handler *UsersHandler) ConfirmRegistration(ctx context.Context, dto *users.ConfirmRegistrationRequest) (*users.ConfirmRegistrationResponse, error) {
	token := dto.GetToken()

	resp := &users.ConfirmRegistrationResponse{
		Success: false,
	}

	success, err := handler.usersService.ConfirmRegistration(ctx, token)
	if err != nil {
		resp.Errors = common.ParseErrors(err)
		return resp, nil
	}

	resp.Success = success

	return resp, nil
}

// TODO !!!
func (handler *UsersHandler) Profile(ctx context.Context, _ *common.Empty) (*users.ProfileResponse, error) {
	var err error
	response := new(users.ProfileResponse)
	return response, err
}

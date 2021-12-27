package users

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/buroz/grpc-clean-example/pkg/domain"
	"github.com/buroz/grpc-clean-example/pkg/storage"
)

var collectionName = "users"

type usersRepository struct {
	col driver.Collection
}

func NewUsersRepository(ctx context.Context, arango storage.ArangoClient) (UsersRepository, error) {
	col, err := arango.Collection(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	_, _, err = col.EnsurePersistentIndex(ctx, []string{"email"}, &driver.EnsurePersistentIndexOptions{
		Unique: true,
	})
	if err != nil {
		return nil, err
	}

	return &usersRepository{col}, nil
}

func (ur *usersRepository) Create(ctx context.Context, user *domain.User) (driver.DocumentID, error) {
	doc, err := ur.col.CreateDocument(ctx, user)
	if err != nil {
		return "", err
	}

	return doc.ID, nil
}

func (ur *usersRepository) FindByEmail(ctx context.Context, email string) ([]domain.User, error) {
	foundUsers := []domain.User{}

	cursor, err := ur.col.Database().Query(ctx,
		`
			FOR doc IN users
			FILTER doc.email == @email
			FILTER doc.status == true
			RETURN doc
		`,
		map[string]interface{}{
			"email": email,
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		u := domain.User{}

		meta, err := cursor.ReadDocument(ctx, &u)
		if err != nil {
			return nil, err
		}

		u.Id = meta.ID

		foundUsers = append(foundUsers, u)
	}

	return foundUsers, err
}

func (ur *usersRepository) FindByAccessToken(ctx context.Context, token string, forceExpired bool) ([]domain.User, error) {
	foundUsers := []domain.User{}

	query := `
	FOR doc IN users
		FILTER doc.status == true
		FILTER doc.latest_token == @token
		RETURN doc
	`

	if forceExpired {
		query = `
		FOR doc IN users
			FILTER doc.status == true
			FILTER doc.latest_token == @token
			FILTER DATE_DIFF(DATE_NOW(), DATE_ADD(doc.last_login, 1, "days"), "s") > 0
			RETURN doc
		`
	}

	cursor, err := ur.col.Database().Query(ctx, query,
		map[string]interface{}{
			"token": token,
		},
	)
	if err != nil {
		return foundUsers, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		u := domain.User{}

		meta, err := cursor.ReadDocument(ctx, &u)
		if err != nil {
			return nil, err
		}

		u.Id = meta.ID

		foundUsers = append(foundUsers, u)
	}

	return foundUsers, err
}

func (ur *usersRepository) FindByConfirmationToken(ctx context.Context, token string) ([]domain.User, error) {
	foundUsers := []domain.User{}

	query := `
	FOR doc IN users
		FILTER doc.status != true
		FILTER doc.confirmation_token == @token
		RETURN doc
	`

	cursor, err := ur.col.Database().Query(ctx, query,
		map[string]interface{}{
			"token": token,
		},
	)
	if err != nil {
		return foundUsers, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		u := domain.User{}

		meta, err := cursor.ReadDocument(ctx, &u)
		if err != nil {
			return nil, err
		}

		u.Id = meta.ID

		foundUsers = append(foundUsers, u)
	}

	return foundUsers, err
}

func (ur *usersRepository) ConfirmUserWithToken(ctx context.Context, key, token string) (bool, error) {
	foundUsers := []domain.User{}

	query := `
	FOR u IN users
		FILTER u.status != true
		FILTER u.confirmation_token == @token
		UPDATE {
			_key: u._key,
			status: true,
			confirmation_token: null
		} IN users OPTIONS { keepNull: false }
		return u
	`

	cursor, err := ur.col.Database().Query(ctx, query,
		map[string]interface{}{
			"token": token,
		},
	)
	if err != nil {
		return len(foundUsers) > 0, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		u := domain.User{}

		meta, err := cursor.ReadDocument(ctx, &u)
		if err != nil {
			return len(foundUsers) > 0, err
		}

		u.Id = meta.ID

		foundUsers = append(foundUsers, u)
	}

	return len(foundUsers) > 0, err
}

func (ur *usersRepository) FindByTokens(ctx context.Context, token string, refreshToken string, forceExpired bool) ([]domain.User, error) {
	foundUsers := []domain.User{}

	query := `
	FOR doc IN users
		FILTER doc.latest_token == @token
		FILTER doc.refresh_token == @refresh_token
		RETURN doc
	`

	if forceExpired {
		query = `
		FOR doc IN users
			FILTER doc.latest_token == @token
			FILTER doc.refresh_token == @refresh_token
			FILTER DATE_DIFF(DATE_NOW(), DATE_ADD(doc.last_login, 1, "days"), "s") > 0
			RETURN doc
		`
	}

	cursor, err := ur.col.Database().Query(ctx, query,
		map[string]interface{}{
			"token":         token,
			"refresh_token": refreshToken,
		},
	)
	if err != nil {
		return foundUsers, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		u := domain.User{}

		meta, err := cursor.ReadDocument(ctx, &u)
		if err != nil {
			return nil, err
		}

		u.Id = meta.ID

		foundUsers = append(foundUsers, u)
	}

	return foundUsers, err
}

func (ur *usersRepository) SetNewTokens(ctx context.Context, user *domain.User) error {
	patch := map[string]interface{}{
		"latest_token":  user.LatestToken,
		"refresh_token": user.RefreshToken,
		"last_login":    time.Now(),
	}

	_, err := ur.col.UpdateDocument(ctx, user.Id.Key(), patch)
	if err != nil {
		return err
	}

	return nil
}

func (ur *usersRepository) SetNewPasswordResetRequest(ctx context.Context, user *domain.User) error {
	patch := map[string]interface{}{
		"password_reset_token":        user.PasswordResetToken,
		"password_reset_requested_at": time.Now(),
		"status":                      false,
	}

	_, err := ur.col.UpdateDocument(ctx, user.Id.Key(), patch)
	if err != nil {
		return err
	}

	return nil
}

func (ur *usersRepository) SetNewPassword(ctx context.Context, user *domain.User) ([]domain.User, error) {
	foundUsers := []domain.User{}

	query := `
	FOR u IN users
		FILTER u.password_reset_token == @token
		FILTER DATE_DIFF(DATE_NOW(), DATE_ADD(u.password_reset_requested_at, 1, "days"), "s") > 0
		UPDATE {
			_key: u._key,
			password: @password,
			password_reset_token: null,
			password_reset_requested_at: null,
			refresh_token: null,
			latest_token: null,
			status: true,
		} IN users OPTIONS { keepNull: false }
		RETURN u
	`

	cursor, err := ur.col.Database().Query(ctx, query,
		map[string]interface{}{
			"token":    user.PasswordResetToken,
			"password": user.Password,
		},
	)
	if err != nil {
		return foundUsers, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		u := domain.User{}

		meta, err := cursor.ReadDocument(ctx, &u)
		if err != nil {
			return foundUsers, err
		}

		u.Id = meta.ID

		foundUsers = append(foundUsers, u)
	}

	return foundUsers, err
}

func (ur *usersRepository) FindById(ctx context.Context, id string) (domain.User, error) {
	user := domain.User{}
	_, err := ur.col.ReadDocument(ctx, id, &user)
	return user, err
}

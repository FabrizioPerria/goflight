package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/handlers/middleware"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupAuthDb() (*testUserDb, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	userStore := db.NewMongoDbUserStore(client)
	mainStore := db.Store{User: userStore}
	return &testUserDb{Store: mainStore, Client: client}, nil
}

func teardownAuthDb(t *testing.T, db *testUserDb) {
	if err := db.Store.User.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := db.Client.Disconnect(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func authenticateUser(authHandler *AuthHandler, app *fiber.App, user UserAuthenticate) (*http.Response, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(userJson))
	req.Header.Set("Content-Type", "application/json")
	response, error := app.Test(req)
	return response, error
}

func TestAuthenticate(t *testing.T) {
	db, err := setupAuthDb()
	assert.NoError(t, err)
	defer teardownAuthDb(t, db)

	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, error := createUser(&userHandler, app, user)
	assert.NoError(t, error)
	assert.Equal(t, 201, response.StatusCode)
	authHandler := AuthHandler{store: db.Store}

	app.Post("/auth", authHandler.HandleAuthenticate)

	userToAuthenticate := UserAuthenticate{
		Email:    user.Email,
		Password: user.PlainPassword,
	}

	authResponse, err := authenticateUser(&authHandler, app, userToAuthenticate)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, authResponse.StatusCode)

	var authResponseDecoded AuthResponse
	err = json.NewDecoder(authResponse.Body).Decode(&authResponseDecoded)
	assert.NoError(t, err)
	assert.Empty(t, authResponseDecoded.User.EncryptedPassword)
	assert.NotEmpty(t, authResponseDecoded.Token)
	token := authResponseDecoded.Token

	assert.NotEmpty(t, token)
	app.Use(middleware.JWTAuthentication)

	app.Get("/api/v1/users", userHandler.HandleGetUsersv1)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("X-Api-Token", token)
	response, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	var users []types.User
	err = json.NewDecoder(response.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, user.Email, users[0].Email)
	assert.True(t, reflect.DeepEqual(authResponseDecoded.User, users[0]))
}

func TestAuthenticateInvalidCredentials(t *testing.T) {
	db, err := setupAuthDb()
	assert.NoError(t, err)
	defer teardownAuthDb(t, db)

	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, error := createUser(&userHandler, app, user)
	assert.NoError(t, error)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
	authHandler := AuthHandler{store: db.Store}

	app.Post("/auth", authHandler.HandleAuthenticate)

	userToAuthenticate := UserAuthenticate{
		Email:    user.Email,
		Password: "wrongpassword",
	}

	authResponse, err := authenticateUser(&authHandler, app, userToAuthenticate)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, authResponse.StatusCode)
}

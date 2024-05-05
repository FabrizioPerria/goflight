package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri    = "mongodb://localhost:27017"
	dbName = "goflight_test"
)

type testUserDb struct {
	Store  db.Store
	Client *mongo.Client
}

func setupUsersDb() (*testUserDb, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	userStore := db.NewMongoDbUserStore(client)
	mainStore := db.Store{User: userStore}
	return &testUserDb{Store: mainStore, Client: client}, nil
}

func teardownUsersDb(t *testing.T, db *testUserDb) {
	if err := db.Store.User.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := db.Client.Disconnect(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func getValidUser() types.CreateUserParams {
	return types.CreateUserParams{
		FirstName:     "Frank",
		LastName:      "Potato",
		Email:         "fp@test.com",
		Phone:         "123456789",
		PlainPassword: "password",
	}
}

func getInvalidUser() types.CreateUserParams {
	return types.CreateUserParams{
		FirstName:     "F",
		LastName:      "P",
		Email:         "fptest,com",
		Phone:         "123456789",
		PlainPassword: "pa",
	}
}

func createUser(userHandler *UserHandler, app *fiber.App, user types.CreateUserParams) (*http.Response, error) {
	app.Post("/api/v1/users", userHandler.HandlePostCreateUserv1)
	marshalUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(marshalUser))
	req.Header.Add("Content-Type", "application/json")
	return app.Test(req)
}

func getUsers(userHandler *UserHandler, app *fiber.App) (*http.Response, error) {
	app.Get("/api/v1/users", userHandler.HandleGetUsersv1)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Add("Content-Type", "application/json")
	return app.Test(req)
}

func TestPostCreateValidUser(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, error := createUser(&userHandler, app, user)
	assert.NoError(t, error)
	assert.Equal(t, 201, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyT.Id)
	assert.Equal(t, user.FirstName, bodyT.FirstName)
	assert.Equal(t, user.LastName, bodyT.LastName)
	assert.Equal(t, user.Email, bodyT.Email)
	assert.Equal(t, user.Phone, bodyT.Phone)
	assert.Equal(t, len(bodyT.EncryptedPassword), 0)
}

func TestPostCreateInvalidUser(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getInvalidUser()
	response, err := createUser(&userHandler, app, user)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := map[string]any{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Contains(t, bodyT, "errors")
	assert.Contains(t, bodyT["errors"], "first_name")
	assert.Contains(t, bodyT["errors"], "last_name")
	assert.Contains(t, bodyT["errors"], "email")

	errors := bodyT["errors"].(map[string]any)

	assert.Equal(t, "first name must be at least 3 characters", errors["first_name"])
	assert.Equal(t, "last name must be at least 3 characters", errors["last_name"])
	assert.Equal(t, "email is not valid", errors["email"])
}

func TestGetUsersEmpty(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	response, err := getUsers(&userHandler, app)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)

	bodyT := []types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(bodyT))
}

func TestGetUsersNotEmpty(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, err := createUser(&userHandler, app, user)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	response, err = getUsers(&userHandler, app)
	assert.NoError(t, err)
	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyT := []types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bodyT))

	assert.NotEmpty(t, bodyT[0].Id)
	assert.Equal(t, user.FirstName, bodyT[0].FirstName)
	assert.Equal(t, user.LastName, bodyT[0].LastName)
	assert.Equal(t, user.Email, bodyT[0].Email)
	assert.Equal(t, user.Phone, bodyT[0].Phone)
	assert.Equal(t, len(bodyT[0].EncryptedPassword), 0)
}

func TestGetUserById(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, err := createUser(&userHandler, app, user)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)

	id := bodyT.Id.Hex()

	req := httptest.NewRequest("GET", "/api/v1/users/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	app.Get("/api/v1/users/:uid", userHandler.HandleGetUserv1)
	response, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err = io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT = types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyT.Id)
	assert.Equal(t, user.FirstName, bodyT.FirstName)
	assert.Equal(t, user.LastName, bodyT.LastName)
	assert.Equal(t, user.Email, bodyT.Email)
	assert.Equal(t, user.Phone, bodyT.Phone)
	assert.Equal(t, len(bodyT.EncryptedPassword), 0)
}

func TestGetUserByIdNotFound(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	store := UserHandler{store: db.Store}

	app := fiber.New()

	req := httptest.NewRequest("GET", "/api/v1/users/16624e25e22069075acbb235", nil)
	req.Header.Add("Content-Type", "application/json")
	app.Get("/api/v1/users/:uid", store.HandleGetUserv1)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
}

func TestDeleteUserById(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, err := createUser(&userHandler, app, user)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)

	id := bodyT.Id.Hex()
	req := httptest.NewRequest("DELETE", "/api/v1/users/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	app.Delete("/api/v1/users/:uid", userHandler.HandleDeleteUserv1)
	response, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err = io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	assert.Equal(t, "User deleted: ", string(body))
}

func TestDeleteUserByIdNotFound(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	req := httptest.NewRequest("DELETE", "/api/v1/users/16624e25e22069075acbb235", nil)
	req.Header.Add("Content-Type", "application/json")
	app.Delete("/api/v1/users/:uid", userHandler.HandleDeleteUserv1)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
}

func TestDeleteAllUsers(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, err := createUser(&userHandler, app, user)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	response, err = getUsers(&userHandler, app)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := []types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bodyT))

	req := httptest.NewRequest("DELETE", "/api/v1/user", nil)
	req.Header.Add("Content-Type", "application/json")
	app.Delete("/api/v1/user", userHandler.HandleDeleteAllUsersv1)
	response, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	response, err = getUsers(&userHandler, app)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err = io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT = []types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(bodyT))
}

func TestDeleteAllUsersEmpty(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	req := httptest.NewRequest("DELETE", "/api/v1/user", nil)
	req.Header.Add("Content-Type", "application/json")
	app.Delete("/api/v1/user", userHandler.HandleDeleteAllUsersv1)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestPutUser(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	user := getValidUser()
	response, err := createUser(&userHandler, app, user)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	id := bodyT.Id.Hex()

	updateUser := types.UpdateUserParams{
		FirstName: "NewName",
		LastName:  "NewLastName",
	}

	marshalUser, err := json.Marshal(updateUser)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/users/"+id, bytes.NewReader(marshalUser))
	req.Header.Add("Content-Type", "application/json")
	app.Put("/api/v1/users/:uid", userHandler.HandlePutUserv1)
	response, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	req = httptest.NewRequest("GET", "/api/v1/users/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	app.Get("/api/v1/users/:uid", userHandler.HandleGetUserv1)
	response, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err = io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)

	bodyT = types.User{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyT.Id)
	assert.Equal(t, user.FirstName, bodyT.FirstName)
	assert.Equal(t, user.LastName, bodyT.LastName)
	assert.Equal(t, user.Email, bodyT.Email)
	assert.Equal(t, user.Phone, bodyT.Phone)
	assert.Equal(t, len(bodyT.EncryptedPassword), 0)
}

func TestPutUserNotFound(t *testing.T) {
	db, err := setupUsersDb()
	assert.NoError(t, err)
	defer teardownUsersDb(t, db)
	userHandler := UserHandler{store: db.Store}

	app := fiber.New()

	updateUser := types.UpdateUserParams{
		FirstName: "NewName",
		LastName:  "NewLastName",
	}

	marshalUser, err := json.Marshal(updateUser)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/users/16624e25e22069075acbb235", bytes.NewReader(marshalUser))
	req.Header.Add("Content-Type", "application/json")
	app.Put("/api/v1/users/:uid", userHandler.HandlePutUserv1)
	response, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
}

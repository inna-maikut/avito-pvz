//go:build integration

package integration

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/api"
)

func Test_RegisterLogin(t *testing.T) {
	setUp()

	email := strconv.Itoa(rand.Int()) + "email@gmail.com"

	resp := apiPost(t, "/login", "", api.PostLoginJSONBody{
		Email:    openapi_types.Email(email),
		Password: "password1",
	})
	assertStatus(t, resp, http.StatusUnauthorized)

	resp = apiPost(t, "/register", "", api.PostRegisterJSONBody{
		Email:    openapi_types.Email(email),
		Password: "password1",
		Role:     api.Moderator,
	})
	assertStatus(t, resp, http.StatusCreated)

	resp = apiPost(t, "/login", "", api.PostLoginJSONBody{
		Email:    openapi_types.Email(email),
		Password: "password1",
	})
	assertStatus(t, resp, http.StatusOK)
	moderatorToken := parseJSON[string](t, resp)
	require.NotEmpty(t, moderatorToken)

	resp = apiPost(t, "/pvz", moderatorToken, api.PostPvzJSONRequestBody{
		City: api.Москва,
	})
	assertStatus(t, resp, http.StatusCreated)
}

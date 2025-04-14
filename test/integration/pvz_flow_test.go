//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

func Test_PVZFlow_OK(t *testing.T) {
	setUp()

	moderatorToken := dummyLogin(t, model.UserRoleModerator)
	employeeToken := dummyLogin(t, model.UserRoleEmployee)

	resp := apiPost(t, "/pvz", moderatorToken, api.PostPvzJSONRequestBody{
		City: api.Москва,
	})
	assertStatus(t, resp, http.StatusCreated)

	pvz := parseJSON[api.PVZ](t, resp)
	require.NotNil(t, pvz.Id)

	resp = apiPost(t, "/receptions", employeeToken, api.PostReceptionsJSONBody{
		PvzId: *pvz.Id,
	})
	assertStatus(t, resp, http.StatusCreated)

	reception := parseJSON[api.Reception](t, resp)
	require.NotNil(t, reception.Id)

	for i := 0; i < 50; i++ {
		resp = apiPost(t, "/products", employeeToken, api.PostProductsJSONBody{
			PvzId: *pvz.Id,
			Type:  api.PostProductsJSONBodyTypeОбувь,
		})
		assertStatus(t, resp, http.StatusCreated)
	}

	resp = apiPost(t, "/pvz/"+pvz.Id.String()+"/close_last_reception", employeeToken, struct{}{})
	assertStatus(t, resp, http.StatusOK)
}

package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simple-golang/internal/adapter/inbound/echo/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignIn_Success(t *testing.T) {
	e, cleanup := InitTestApp()
	defer cleanup()

	// ===== REQUEST =====
	reqBody := request.SignInRequest{
		Email:    "superadmin@mail.com",
		Password: "12345678",
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/signin",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()

	// ===== HIT API =====
	e.ServeHTTP(rec, req)

	// ===== ASSERT STATUS =====
	assert.Equal(t, http.StatusOK, rec.Code)

	// ===== PARSE RESPONSE =====
	var resp struct {
		Message string `json:"message"`
		Data    struct {
			AccessToken string `json:"access_token"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			Phone       string `json:"phone"`
		} `json:"data"`
	}

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// ===== ASSERT RESPONSE =====
	assert.Equal(t, "Success", resp.Message)
	assert.NotEmpty(t, resp.Data.AccessToken)
	assert.Equal(t, reqBody.Email, resp.Data.Email)
	assert.NotEmpty(t, resp.Data.Name)
}

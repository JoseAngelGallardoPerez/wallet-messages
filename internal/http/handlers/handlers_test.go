package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Confialink/wallet-messages/internal/model"

	"github.com/Confialink/wallet-messages/internal/auth"
)

// performRequest
func performRequest(r http.Handler, method, path string, headers map[string]string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	return res
}

// getClientAuthHeader returns map containing clients bearer token
func getClientAuthHeader() map[string]string {
	headers := map[string]string{}
	headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ImNsaWVudCJ9.Rx5Vz-kwp7cJocuRx3xvrsgRCtf_D5fsF1IxZPPVCaw"
	return headers
}

// getAdminAuthHeader returns map containing admins bearer token
func getAdminAuthHeader() map[string]string {
	headers := map[string]string{}
	headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ImFkbWluIn0.oDr14tn99Ma3wUe-6FAlWgsYWrHe6YnbOL9VCmhdcTA"
	return headers
}

// createAuthService creates new auth service instance
func createAuthService() *auth.Service {
	return auth.NewService(auth.NewPermissionService())
}

// Repository is a fake repository to avoid real database requests
type Repository struct{}

// FindByUserAndParams retrieve the list of messages
func (repo *Repository) FindByUserAndParams(userId int64, params url.Values) ([]*model.Message, error) {
	return []*model.Message{}, nil
}

// FindAllUsers retrieve the list of users
func (repo *Repository) CountByUserAndParams(userId int64, params url.Values) (*int64, error) {
	count := int64(10)
	return &count, nil
}

// FindByUserAndParams retrieve the list of messages
func (repo *Repository) FindUnassigned(params url.Values) ([]*model.Message, error) {
	return []*model.Message{}, nil
}

// FindAllUsers retrieve the list of users
func (repo *Repository) CountUnassigned(params url.Values) (*int64, error) {
	count := int64(10)
	return &count, nil
}

// FindByID find user by id
func (repo *Repository) FindByID(id uint) (*model.Message, error) {
	idStr := string(id)
	private := model.MessagePrivate{SenderId: &idStr}
	return &model.Message{MessagePrivate: private}, nil
}

// Create creates new user
func (repo *Repository) Create(message *model.Message) (*model.Message, error) {
	return message, nil
}

// Update updates an existing user
func (repo *Repository) Update(message *model.Message) (*model.Message, error) {
	return message, nil
}

// Delete delete an existing user
func (repo *Repository) Delete(message *model.Message) error {
	return nil
}

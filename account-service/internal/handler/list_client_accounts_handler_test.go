package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/handler"
	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/models"
)

type mockListClientAccountsRepo struct {
	accounts []models.Account
	err      error
}

func (m *mockListClientAccountsRepo) ListByClientID(clientID uint) ([]models.Account, error) {
	return m.accounts, m.err
}

func TestListClientAccountsHTTP_ReturnsDnevnaMesecnaPotrosnja(t *testing.T) {
	accounts := []models.Account{
		{
			ID:                1,
			BrojRacuna:        "000100000000000001",
			Stanje:            5000,
			RaspolozivoStanje: 4500,
			DnevniLimit:       100000,
			MesecniLimit:      1000000,
			DnevnaPotrosnja:   300,
			MesecnaPotrosnja:  1500,
			Currency:          models.Currency{Kod: "RSD"},
			Status:            "aktivan",
		},
	}

	h := handler.NewListClientAccountsHTTPHandler(&mockListClientAccountsRepo{accounts: accounts})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/client/42", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 account, got %d", len(resp))
	}

	acc := resp[0]
	if _, ok := acc["dnevnaPotrosnja"]; !ok {
		t.Error("response missing dnevnaPotrosnja field")
	}
	if _, ok := acc["mesecnaPotrosnja"]; !ok {
		t.Error("response missing mesecnaPotrosnja field")
	}
	if acc["dnevnaPotrosnja"].(float64) != 300 {
		t.Errorf("expected dnevnaPotrosnja=300, got %v", acc["dnevnaPotrosnja"])
	}
	if acc["mesecnaPotrosnja"].(float64) != 1500 {
		t.Errorf("expected mesecnaPotrosnja=1500, got %v", acc["mesecnaPotrosnja"])
	}
}

func TestListClientAccountsHTTP_InvalidClientID_ReturnsBadRequest(t *testing.T) {
	h := handler.NewListClientAccountsHTTPHandler(&mockListClientAccountsRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/client/not-a-number", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

package models_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/account-service/internal/models"
)

func TestOvlascenoLice_HasRequiredFields(t *testing.T) {
	rt := reflect.TypeOf(models.OvlascenoLice{})
	required := []string{"ID", "Ime", "Prezime", "Email", "BrojTelefona", "FirmaID", "AccountID"}
	for _, name := range required {
		if _, ok := rt.FieldByName(name); !ok {
			t.Errorf("OvlascenoLice missing field: %s", name)
		}
	}
}

func TestOvlascenoLice_IDIsPrimaryKey(t *testing.T) {
	rt := reflect.TypeOf(models.OvlascenoLice{})
	f, ok := rt.FieldByName("ID")
	if !ok {
		t.Fatal("ID field not found")
	}
	tag := f.Tag.Get("gorm")
	if !strings.Contains(tag, "primaryKey") {
		t.Errorf("ID: expected gorm tag to contain primaryKey, got: %s", tag)
	}
}

func TestOvlascenoLice_JSONTags(t *testing.T) {
	rt := reflect.TypeOf(models.OvlascenoLice{})
	cases := map[string]string{
		"ID":           "id",
		"Ime":          "ime",
		"Prezime":      "prezime",
		"Email":        "email",
		"BrojTelefona": "broj_telefona",
		"FirmaID":      "firma_id",
		"AccountID":    "account_id",
	}
	for field, expectedJSON := range cases {
		f, ok := rt.FieldByName(field)
		if !ok {
			t.Errorf("field %s not found", field)
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag != expectedJSON {
			t.Errorf("%s: expected json:%q, got %q", field, expectedJSON, jsonTag)
		}
	}
}

func TestOvlascenoLice_CanInstantiate(t *testing.T) {
	ol := models.OvlascenoLice{
		Ime:          "Petar",
		Prezime:      "Petrovic",
		Email:        "petar@example.com",
		BrojTelefona: "+381601234567",
		FirmaID:      1,
		AccountID:    2,
	}
	if ol.Ime != "Petar" {
		t.Errorf("expected Ime=Petar, got %s", ol.Ime)
	}
	if ol.FirmaID != 1 {
		t.Errorf("expected FirmaID=1, got %d", ol.FirmaID)
	}
	if ol.AccountID != 2 {
		t.Errorf("expected AccountID=2, got %d", ol.AccountID)
	}
}

func TestCard_HasOvlascenoLiceID(t *testing.T) {
	rt := reflect.TypeOf(models.Card{})
	f, ok := rt.FieldByName("OvlascenoLiceID")
	if !ok {
		t.Fatal("Card missing field: OvlascenoLiceID")
	}
	// Should be a pointer (optional)
	if f.Type.Kind() != reflect.Ptr {
		t.Errorf("OvlascenoLiceID should be a pointer (optional), got %s", f.Type.Kind())
	}
}

package models

// OvlascenoLice represents an authorized person linked to a business (poslovni) account.
// It is informational — used to display authorized signatories in account/card details.
type OvlascenoLice struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Ime          string `json:"ime"`
	Prezime      string `json:"prezime"`
	Email        string `json:"email"`
	BrojTelefona string `json:"broj_telefona"`
	FirmaID      uint   `json:"firma_id"`
	AccountID    uint   `json:"account_id"` // poslovni račun
}

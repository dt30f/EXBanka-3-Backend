package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/exchange-service/internal/repository"
)

// TaxCollector handles monthly tax collection from users' RSD accounts.
type TaxCollector struct {
	taxSvc    *TaxService
	orderRepo *repository.OrderRepository
	taxRepo   *repository.TaxRepository
}

// NewTaxCollector creates a new TaxCollector.
func NewTaxCollector(taxSvc *TaxService, orderRepo *repository.OrderRepository, taxRepo *repository.TaxRepository) *TaxCollector {
	return &TaxCollector{
		taxSvc:    taxSvc,
		orderRepo: orderRepo,
		taxRepo:   taxRepo,
	}
}

// TaxDebt represents a failed collection for a user.
type TaxDebt struct {
	UserID   uint
	UserType string
	Amount   float64
	Reason   string
}

// CollectionResult summarises the outcome of a tax collection run.
type CollectionResult struct {
	Period         string
	UsersProcessed int
	TotalCollected float64
	Debts          []TaxDebt
}

// CollectForPeriod collects all unpaid capital-gains tax for the given period
// (format "YYYY-MM") by debiting users' RSD accounts and crediting the state
// treasury account.  Users who cannot pay are recorded as debts.
func (c *TaxCollector) CollectForPeriod(period string) CollectionResult {
	result := CollectionResult{Period: period}

	users, err := c.taxRepo.ListDistinctUsersWithUnpaidTax(period)
	if err != nil {
		slog.Error("tax_cron: failed to list users with unpaid tax", "period", period, "error", err)
		return result
	}

	treasuryID, err := c.orderRepo.GetStateTreasuryAccountID()
	if err != nil {
		slog.Error("tax_cron: failed to get state treasury account", "error", err)
		return result
	}
	if treasuryID == 0 {
		slog.Warn("tax_cron: state treasury account not found, skipping collection", "period", period)
		return result
	}

	for _, u := range users {
		totalOwed, err := c.taxRepo.SumUnpaidTaxForUser(u.UserID, u.UserType, period)
		if err != nil || totalOwed <= 0 {
			continue
		}

		accounts, err := c.orderRepo.GetUserRSDAccounts(u.UserID, u.UserType)
		if err != nil {
			slog.Warn("tax_cron: cannot get RSD accounts for user", "userID", u.UserID, "error", err)
			result.Debts = append(result.Debts, TaxDebt{
				UserID:   u.UserID,
				UserType: u.UserType,
				Amount:   totalOwed,
				Reason:   fmt.Sprintf("failed to retrieve RSD accounts: %v", err),
			})
			continue
		}

		// Attempt to debit the full tax amount from the first account with
		// sufficient balance.
		collected := false
		for _, acc := range accounts {
			if acc.RaspolozivoStanje < totalOwed {
				continue
			}
			if err := c.orderRepo.DebitAccount(acc.ID, totalOwed); err != nil {
				continue
			}
			// Credit the state treasury.
			if err := c.orderRepo.CreditAccount(treasuryID, totalOwed); err != nil {
				// Best-effort: log but don't roll back the debit (treasury credit
				// failure is an operational issue, not a user-facing one).
				slog.Error("tax_cron: failed to credit treasury", "userID", u.UserID, "amount", totalOwed, "error", err)
			}
			if err := c.taxRepo.MarkTaxRecordsPaid(u.UserID, u.UserType, period); err != nil {
				slog.Error("tax_cron: failed to mark records paid", "userID", u.UserID, "error", err)
			}
			result.TotalCollected += totalOwed
			collected = true
			break
		}

		if !collected {
			result.Debts = append(result.Debts, TaxDebt{
				UserID:   u.UserID,
				UserType: u.UserType,
				Amount:   totalOwed,
				Reason:   "insufficient balance in all RSD accounts",
			})
		}

		result.UsersProcessed++
	}

	slog.Info("tax_cron: collection complete",
		"period", period,
		"users_processed", result.UsersProcessed,
		"total_collected_rsd", result.TotalCollected,
		"debts", len(result.Debts),
	)
	return result
}

// PreviousMonthPeriod returns the "YYYY-MM" string for the calendar month
// preceding the current month.
func PreviousMonthPeriod() string {
	t := time.Now().UTC().AddDate(0, -1, 0)
	return fmt.Sprintf("%d-%02d", t.Year(), t.Month())
}

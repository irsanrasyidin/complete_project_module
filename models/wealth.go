package models

import "time"

type Wallet struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   string    `json:"balance"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IncomeCategory struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Code        string    `json:"code" gorm:"uniqueIndex"`
	Description string    `json:"description"`
	Example     string    `json:"example"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpenseCategory struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Code        string    `json:"code" gorm:"uniqueIndex"`
	Description string    `json:"description"`
	Example     string    `json:"example"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Rule struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MasterData struct {
	Wallets          []Wallet          `json:"wallets"`
	IncomeCategories  []IncomeCategory  `json:"income_categories"`
	ExpenseCategories []ExpenseCategory `json:"expense_categories"`
	Rules            []Rule            `json:"rules"`
}

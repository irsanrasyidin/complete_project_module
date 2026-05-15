package masterdata

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/irsanrasyidin/complete_project/module/models"
	"gorm.io/gorm"
)

type Repository interface {
	EnsureSchema(ctx context.Context) error
	ClearMasterData(ctx context.Context) error
	SeedDefaultMasterData(ctx context.Context) error
	GetMasterData(ctx context.Context) (models.MasterData, error)
	ListWallets(ctx context.Context) ([]models.Wallet, error)
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	ListIncomeCategories(ctx context.Context) ([]models.IncomeCategory, error)
	CreateIncomeCategory(ctx context.Context, category *models.IncomeCategory) error
	ListExpenseCategories(ctx context.Context) ([]models.ExpenseCategory, error)
	CreateExpenseCategory(ctx context.Context, category *models.ExpenseCategory) error
	ListRules(ctx context.Context) ([]models.Rule, error)
	CreateRule(ctx context.Context, rule *models.Rule) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) EnsureSchema(ctx context.Context) error {
	return r.db.WithContext(ctx).AutoMigrate(
		&models.Wallet{},
		&models.IncomeCategory{},
		&models.ExpenseCategory{},
		&models.Rule{},
	)
}

func (r *repository) ClearMasterData(ctx context.Context) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	tables := []any{
		&models.Wallet{},
		&models.IncomeCategory{},
		&models.ExpenseCategory{},
		&models.Rule{},
	}
	for _, table := range tables {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			_ = tx.Rollback().Error
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) SeedDefaultMasterData(ctx context.Context) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&[]models.Wallet{
		{ID: newID(), Name: "Cash Wallet", Type: "Cash", Balance: "Rp 12.500.000", Status: "Active", Note: "Uang tunai harian"},
		{ID: newID(), Name: "Bank BCA", Type: "Bank Account", Balance: "Rp 84.200.000", Status: "Active", Note: "Rekening utama"},
		{ID: newID(), Name: "E-Wallet", Type: "Digital Wallet", Balance: "Rp 3.150.000", Status: "Active", Note: "Pembayaran cepat"},
		{ID: newID(), Name: "Investment Wallet", Type: "Investment", Balance: "Rp 246.000.000", Status: "Active", Note: "Dana yang sedang diinvestasikan"},
	}).Error; err != nil {
		_ = tx.Rollback().Error
		return err
	}

	if err := tx.Create(&[]models.IncomeCategory{
		{ID: newID(), Name: "Salary", Code: "INC-SAL", Description: "Gaji tetap bulanan", Example: "Transfer payroll"},
		{ID: newID(), Name: "Bonus", Code: "INC-BON", Description: "Bonus tahunan atau insentif", Example: "Performance bonus"},
		{ID: newID(), Name: "Business", Code: "INC-BIZ", Description: "Pendapatan dari usaha", Example: "Invoice client"},
		{ID: newID(), Name: "Dividend", Code: "INC-DIV", Description: "Pendapatan pasif dari aset", Example: "ETF dividend"},
		{ID: newID(), Name: "Other income", Code: "INC-OTH", Description: "Pendapatan lain-lain", Example: "Refund / cashback"},
	}).Error; err != nil {
		_ = tx.Rollback().Error
		return err
	}

	if err := tx.Create(&[]models.ExpenseCategory{
		{ID: newID(), Name: "Rent", Code: "EXP-REN", Description: "Biaya tempat tinggal", Example: "Monthly rent"},
		{ID: newID(), Name: "Utilities", Code: "EXP-UTL", Description: "Listrik, air, internet", Example: "Telco bill"},
		{ID: newID(), Name: "Food", Code: "EXP-FOD", Description: "Makan dan kebutuhan harian", Example: "Groceries"},
		{ID: newID(), Name: "Transport", Code: "EXP-TRN", Description: "Transportasi dan bensin", Example: "Fuel / ride hailing"},
		{ID: newID(), Name: "Health", Code: "EXP-HLT", Description: "Kesehatan dan obat", Example: "Clinic visit"},
		{ID: newID(), Name: "Education", Code: "EXP-EDU", Description: "Kursus dan belajar", Example: "Online class"},
		{ID: newID(), Name: "Lifestyle", Code: "EXP-LIF", Description: "Hiburan dan personal care", Example: "Entertainment"},
	}).Error; err != nil {
		_ = tx.Rollback().Error
		return err
	}

	if err := tx.Create(&[]models.Rule{
		{ID: newID(), Title: "One wallet, one purpose", Description: "Setiap wallet punya fungsi jelas agar pencatatan rapi."},
		{ID: newID(), Title: "Income must be categorized", Description: "Semua pemasukan wajib masuk ke jenis pemasukan tertentu."},
		{ID: newID(), Title: "Expense must be tagged", Description: "Pengeluaran diberi kategori agar analisis lebih mudah."},
	}).Error; err != nil {
		_ = tx.Rollback().Error
		return err
	}

	return tx.Commit().Error
}

func (r *repository) GetMasterData(ctx context.Context) (models.MasterData, error) {
	wallets, err := r.ListWallets(ctx)
	if err != nil {
		return models.MasterData{}, err
	}
	income, err := r.ListIncomeCategories(ctx)
	if err != nil {
		return models.MasterData{}, err
	}
	expense, err := r.ListExpenseCategories(ctx)
	if err != nil {
		return models.MasterData{}, err
	}
	rules, err := r.ListRules(ctx)
	if err != nil {
		return models.MasterData{}, err
	}
	return models.MasterData{
		Wallets:          wallets,
		IncomeCategories: income,
		ExpenseCategories: expense,
		Rules:            rules,
	}, nil
}

func (r *repository) ListWallets(ctx context.Context) ([]models.Wallet, error) {
	var items []models.Wallet
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	wallet.ID = newID()
	now := time.Now().UTC()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *repository) ListIncomeCategories(ctx context.Context) ([]models.IncomeCategory, error) {
	var items []models.IncomeCategory
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) CreateIncomeCategory(ctx context.Context, category *models.IncomeCategory) error {
	category.ID = newID()
	now := time.Now().UTC()
	category.CreatedAt = now
	category.UpdatedAt = now
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *repository) ListExpenseCategories(ctx context.Context) ([]models.ExpenseCategory, error) {
	var items []models.ExpenseCategory
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) CreateExpenseCategory(ctx context.Context, category *models.ExpenseCategory) error {
	category.ID = newID()
	now := time.Now().UTC()
	category.CreatedAt = now
	category.UpdatedAt = now
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *repository) ListRules(ctx context.Context) ([]models.Rule, error) {
	var items []models.Rule
	if err := r.db.WithContext(ctx).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) CreateRule(ctx context.Context, rule *models.Rule) error {
	rule.ID = newID()
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	return r.db.WithContext(ctx).Create(rule).Error
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}

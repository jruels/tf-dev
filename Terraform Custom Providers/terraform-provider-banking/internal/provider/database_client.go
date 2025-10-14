package provider

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type BankingDBClient struct {
	DB *sql.DB
}

// func NewBankingDBClient(host string, port int64, user, password, dbname string) (*BankingDBClient, error) {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &BankingDBClient{DB: db}, nil
// }

// **Create a new customer account**
func (c *BankingDBClient) CreateCustomerAccount(firstName, lastName, email, accountType string, balance float64) (int, error) {
	var accountID int
	err := c.DB.QueryRow(
		`INSERT INTO customer_accounts (first_name, last_name, email, account_type, balance)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		firstName, lastName, email, accountType, balance,
	).Scan(&accountID)
	if err != nil {
		return 0, err
	}
	return accountID, nil
}

// **Retrieve customer account details**
func (c *BankingDBClient) GetCustomerAccount(email string) (*CustomerAccount, error) {
	var account CustomerAccount
	err := c.DB.QueryRow(
		`SELECT id, first_name, last_name, email, account_type, balance FROM customer_accounts WHERE email=$1`,
		email,
	).Scan(&account.ID, &account.FirstName, &account.LastName, &account.Email, &account.AccountType, &account.Balance)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// **Update an existing customer account**
func (c *BankingDBClient) UpdateCustomerAccount(id string, firstName, lastName, email, accountType string, balance float64) error {
	_, err := c.DB.Exec(
		`UPDATE customer_accounts 
		 SET first_name=$1, last_name=$2, email=$3, account_type=$4, balance=$5 
		 WHERE id=$6`,
		firstName, lastName, email, accountType, balance, id,
	)
	return err
}

// **Delete a customer account**
func (c *BankingDBClient) DeleteCustomerAccount(email string) error {
	result, err := c.DB.Exec("DELETE FROM customer_accounts WHERE email = $1", email)
	if err != nil {
		return err
	}

	// ✅ Ensure at least one row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account with email '%s' not found", email)
	}

	return nil
}

// **CustomerAccount Struct**
type CustomerAccount struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	AccountType string
	Balance     float64
}

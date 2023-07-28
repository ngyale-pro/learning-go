package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all function to execute db queries and transaction
// This structure is a composition which extends the capacity of the structure Query (similar to inheritence)
type Store struct {
	db *sql.DB
	*Queries
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		fmt.Println("Rollback...")
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("Transaction err: %v, Rollback err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// var transactionKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update account's balance within a single database
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//fmt.Println(transactionName, ": Creating transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(transactionName, ": Creating Entry FromAccount")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -(arg.Amount),
		})
		if err != nil {
			return err
		}

		//fmt.Println(transactionName, ": Creating Entry ToAccount")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	account1ID int64,
	account1Amount int64,
	account2ID int64,
	account2Amount int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: account1Amount,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: account2Amount,
	})

	return
}

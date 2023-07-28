package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {
	existed := make(map[int]bool)

	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fmt.Println("Before -> FromAccount: ", fromAccount.Balance, " | ToAccount: ", toAccount.Balance)

	args := TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        10,
	}

	n := 5

	err_channel := make(chan error)
	result_channel := make(chan TransferTxResults)

	for i := 0; i < n; i++ {
		// transactionName := fmt.Sprintf("Transaction: %d", i)
		go func() {
			ctx := context.Background()
			// ctx := context.WithValue(context.Background(), transactionKey, transactionName)
			result, err := store.TransferTx(ctx, args)

			err_channel <- err
			result_channel <- result
		}()
	}

	for i := 0; i < n; i++ {
		// fmt.Println("Checking errors..., i:", i)
		err := <-err_channel
		require.NoError(t, err)
		result := <-result_channel
		require.NotEmpty(t, result)

		// Check transfer

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, args.Amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check fromEntry

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -args.Amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Check toEntry

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, args.Amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check Account
		resFromAccount := result.FromAccount
		require.NotEmpty(t, resFromAccount)
		require.Equal(t, fromAccount.ID, resFromAccount.ID)

		resToAccount := result.ToAccount
		require.NotEmpty(t, resToAccount)
		require.Equal(t, toAccount.ID, resToAccount.ID)

		// Check Balance
		fmt.Println("Transaction Balance:", i, " -> FromAccount: ", resFromAccount.Balance, " | ToAccount: ", resToAccount.Balance)
		diff1 := fromAccount.Balance - resFromAccount.Balance // 100 - 90
		diff2 := resToAccount.Balance - toAccount.Balance     // 110 - 100

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%args.Amount == 0) // Diff1 is a multiple of args.amount

		number_of_transfer := int(diff1 / args.Amount)
		require.True(t, number_of_transfer >= 1 && number_of_transfer <= n) // The number of transfer goes from 1 to the number of goroutines
		require.NotContains(t, existed, number_of_transfer)
		existed[number_of_transfer] = true
	}

	// Check Account final balance
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println("After -> FromAccount: ", updatedFromAccount.Balance, " | ToAccount: ", updatedToAccount.Balance)
	require.Equal(t, updatedFromAccount.Balance, fromAccount.Balance-(int64(n)*args.Amount))
	require.Equal(t, updatedToAccount.Balance, toAccount.Balance+(int64(n)*args.Amount))

}

func TestTransferTransactionDeadlock(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fmt.Println("Before -> FromAccount: ", fromAccount.Balance, " | ToAccount: ", toAccount.Balance)

	n := 10

	err_channel := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := fromAccount.ID
		toAccountID := toAccount.ID

		if i%2 == 1 {
			fromAccountID = toAccount.ID
			toAccountID = fromAccount.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        10,
			})

			err_channel <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-err_channel
		require.NoError(t, err)
	}

	// Check Account final balance
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println("After -> FromAccount: ", updatedFromAccount.Balance, " | ToAccount: ", updatedToAccount.Balance)
	require.Equal(t, updatedFromAccount.Balance, fromAccount.Balance)
	require.Equal(t, updatedToAccount.Balance, toAccount.Balance)

}

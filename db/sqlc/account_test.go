package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ngyale-pro/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	res := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), res.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, res.ID, account.ID)
	require.Equal(t, res.Owner, account.Owner)
	require.Equal(t, res.Balance, account.Balance)
	require.Equal(t, res.Currency, account.Currency)
	require.WithinDuration(t, res.CreatedAt, account.CreatedAt, time.Second) // Created at similar time (delta in second)
}

func TestUpdateAccount(t *testing.T) {
	res := createRandomAccount(t)
	updateArgs := UpdateAccountParams{
		ID:      res.ID,
		Balance: util.RandomMoney(),
	}

	testQueries.UpdateAccount(context.Background(), updateArgs)

	account, err := testQueries.GetAccount(context.Background(), res.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, res.ID, account.ID)
	require.Equal(t, res.Owner, account.Owner)
	require.Equal(t, updateArgs.Balance, account.Balance)
	require.Equal(t, res.Currency, account.Currency)
	require.WithinDuration(t, res.CreatedAt, account.CreatedAt, time.Second) // Created at similar time (delta in second)
}

func TestDeleteAccount(t *testing.T) {
	res := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), res.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), res.ID)

	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5, // Skip the first five
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

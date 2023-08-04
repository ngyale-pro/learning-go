package db

import (
	"context"
	"testing"
	"time"

	"github.com/ngyale-pro/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)

	args := CreateEntryParams{
		AccountID: -1,
		Amount:    100,
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	t.Logf("Err: %v", err)
	if err == nil {
		t.Errorf("Error: %v", err)
	}

	if (entry != Entry{}) {
		t.Errorf("Error: %v", err)
	}
}

func TestGetEntry(t *testing.T) {
	res := createRandomEntry(t)
	entry, err := testQueries.GetEntry(context.Background(), res.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.ID, res.ID)
	require.Equal(t, entry.AccountID, res.AccountID)
	require.Equal(t, entry.Amount, res.Amount)

	require.WithinDuration(t, entry.CreatedAt, res.CreatedAt, time.Second) // Created at similar time (delta in second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		args := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomInteger(1, int64(i+1)),
		}

		entry, err := testQueries.CreateEntry(context.Background(), args)
		require.NoError(t, err)
		require.NotEmpty(t, entry)

		require.Equal(t, entry.AccountID, account.ID)
		require.Equal(t, entry.Amount, args.Amount)

		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5, // Skip the first five
	}

	entries, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func createRandomEntry(t *testing.T) Entry {

	account := createRandomAccount(t)

	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomInteger(1, account.Balance),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, args.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

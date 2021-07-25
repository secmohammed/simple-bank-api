package db

import (
	"context"
	"testing"

	"github.com/secmohammed/simple-bank-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount() (*Account, error, *CreateAccountParams) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	return &account, err, &arg
}

func TestCreateAccount(t *testing.T) {
	account, err, arg := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	created, _, _ := createRandomAccount()
	account, err := testQueries.GetAccount(context.Background(), created.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, created.ID, account.ID)
	require.Equal(t, created.Balance, account.Balance)
	require.Equal(t, created.Owner, account.Owner)
}

func TestUpdateAccount(t *testing.T) {
	created, _, _ := createRandomAccount()
	arg := UpdateAccountParams{
		ID:      created.ID,
		Balance: util.RandomMoney(),
	}
	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, created.ID, account.ID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, created.Owner, account.Owner)
}

func TestListAccounts(t *testing.T) {
	createRandomAccount()
	createRandomAccount()
	arg := ListAccountsParams{
		Limit:  2,
		Offset: 2,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 2)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestDeleteAccount(t *testing.T) {
	created, _, _ := createRandomAccount()
	err := testQueries.DeleteAccount(context.Background(), created.ID)
	require.NoError(t, err)
}

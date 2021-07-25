package db

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
    store := NewStore(testDB)
    account1, _, _ := createRandomAccount()
    account2, _, _ := createRandomAccount()
    n := 5
    amount := int64(10)
    errs := make(chan error)
    results := make(chan TransferTxResult)
    for i := 0; i < n; i++ {
        go func() {
            ctx := context.Background()

            result, err := store.TransferTx(ctx, TransferTxParams{
                FromAccountID: account1.ID,
                ToAccountID:   account2.ID,
                Amount:        amount,
            })
            errs <- err
            results <- result
        }()
    }
    existed := make(map[int]bool)
    // check results
    for i := 0; i < n; i++ {
        err := <-errs
        require.NoError(t, err)
        result := <-results
        require.NotEmpty(t, result)
        // check transfer
        transfer := result.Transfer
        require.NotEmpty(t, transfer)
        require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
        require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
        require.Equal(t, amount, transfer.Amount)
        require.NotZero(t, transfer.ID)
        require.NotZero(t, transfer.CreatedAt.Time)
        _, err = store.GetTransfer(context.Background(), transfer.ID)
        require.NoError(t, err)
        // check entries
        fromEntry := result.FromEntry
        require.NotEmpty(t, fromEntry)
        require.Equal(t, account1.ID, fromEntry.AccountID.Int64)
        require.Equal(t, -amount, fromEntry.Amount)
        require.NotZero(t, fromEntry.ID)
        require.NotZero(t, fromEntry.CreatedAt)

        _, err = store.GetEntry(context.Background(), fromEntry.ID)
        require.NoError(t, err)
        toEntry := result.ToEntry
        require.NotEmpty(t, toEntry)
        require.Equal(t, account2.ID, toEntry.AccountID.Int64)
        require.Equal(t, amount, toEntry.Amount)
        require.NotZero(t, toEntry.ID)
        require.NotZero(t, toEntry.CreatedAt.Time)

        _, err = store.GetEntry(context.Background(), toEntry.ID)
        require.NoError(t, err)
        // check accounts
        fromAccount := result.FromAccount
        require.NotEmpty(t, fromAccount)
        require.Equal(t, account1.ID, fromAccount.ID)

        toAccount := result.ToAccount
        require.NotEmpty(t, toAccount)
        require.Equal(t, account2.ID, toAccount.ID)

        diff1 := account1.Balance - fromAccount.Balance
        diff2 := toAccount.Balance - account2.Balance
        require.Equal(t, diff1, diff2)
        require.True(t, diff1 > 0)
        require.True(t, diff1%amount == 0) // amount, amount * 2 , 3 * amount, .... n* amount
        k := int(diff1 / amount)
        require.True(t, k >= 1 && k <= n)
        require.NotContains(t, existed, k)
        existed[k] = true
    }
    // check the final updated balances
    updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
    require.NoError(t, err)
    updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
    require.NoError(t, err)

    require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
    require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

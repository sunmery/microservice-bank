package query

import (
	"backend/internal/data/db/gen"
	"context"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
    conn := Connect()
	store := NewStore(conn)

	fake := faker.New()
	
	// 编写5个并发来测试事务是否一致
	account1,err := store.CreateAccount(context.Background(), &gen.CreateAccountParams{
		Owner:    fake.Person().Name(),
		Balance:  0,
		Currency: "USD",
	})
	if err != nil {
		t.Fatal(err)
	}
	account2,err2 := store.CreateAccount(context.Background(), &gen.CreateAccountParams{
		Owner:    fake.Person().Name(),
		Balance:  0,
		Currency: "USD",
	})
	if err2 != nil {
		t.Fatal(err)
	}

	n := 5
	var amount int64 = 10

	// 保证并发的安全性以及安全的共享
	errs := make(chan error)
	results := make(chan TransferResult)
	for i := 0; i < n; i++ {
		go func() {

		result, err := store.Transaction(context.Background(), TransferParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount: amount,
		})

		errs <- err
		results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.NotZero(t, transfer.ID)

		// 查找from创建转账记录
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID,account1.ID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// 查找to创建转账记录
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

	}
}

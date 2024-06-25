package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/jaswdr/faker/v2"

	"user/internal/biz/repository"

	"github.com/stretchr/testify/require"
)

var txKey = struct{}{} // 空对象

func TestTransferTx(t *testing.T) {
	conn := Connect()
	store := NewStore(conn)

	fake := faker.New()

	// 编写5个并发来测试事务是否一致
	mockUser1Name := fake.Person().Name()
	account1, err := store.CreateAccount(context.Background(), &repository.CreateAccountParams{
		Owner:    mockUser1Name,
		Balance:  20,
		Currency: "USD",
	})
	require.NoError(t, err)

	mockUser2Name := fake.Person().Name()
	account2, err2 := store.CreateAccount(context.Background(), &repository.CreateAccountParams{
		Owner:    mockUser2Name,
		Balance:  0,
		Currency: "USD",
	})
	require.NoError(t, err2)

	fmt.Printf("tx: before 1: %v 2:%v\n", account1.Balance, account2.Balance)

	n := 2
	var amount int64 = 10

	// 保证并发的安全性以及安全的共享
	errs := make(chan error)
	results := make(chan TransferResult)


	for i := 0; i < n; i++ {
		i := i
		go func() {
			txName := fmt.Sprintf("tx %d", i+1)
			result, err := store.Transaction(context.WithValue(context.Background(), txKey, txName), TransferParams{
				FromAccount: account1.ID,
				ToAccount:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	var existed = make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// 检查转账
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// 检查账单
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// 检查账户余额
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		fmt.Printf("交易后的余额: %v: %v, %v: %v\n", fromAccount.Owner, fromAccount.Balance, toAccount.Owner, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k:= int(diff1 / amount)
		// 次数
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed,k)
		existed[k] = true
	}
		// 检查账户
}

package data

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

var Db *Data

func TestNewDatabase(t *testing.T)  {
	db := Db.db
	testQueries = gen.New(conn)

	arg := gen.CreateAccountParams{
		Owner:    "tom2",
		Balance:  100,
		Currency: "USD",
	}

	account, err2 := testQueries.CreateAccount(context.Background(), &arg)
	require.NoError(t, err2)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

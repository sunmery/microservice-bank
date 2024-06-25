package query

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stretchr/testify/require"

	"user/internal/biz/repository"
)

var testQueries *repository.Queries

func Setup() *pgxpool.Pool {
	// # Example Keyword/Value
	// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
	// # Example URL
	// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	url := "postgres://root:msdnmm@103.71.69.244:5432/bank?sslmode=disable&pool_max_conns=10"
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		os.Exit(1)
	}

	return conn
}

func TestSqlc(t *testing.T) {
	conn := Setup()
	defer conn.Close()

	testQueries = repository.New(conn)

	arg := repository.CreateAccountParams{
		Owner:    "tom4",
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

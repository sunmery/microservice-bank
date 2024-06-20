package query

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/internal/data/db/gen"
)

var testQueries *gen.Queries

func Setup() *pgxpool.Pool {
	// # Example Keyword/Value
	// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
	// # Example URL
	// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	url := "postgres://root:msdnmm@103.71.69.244:5432/bank?sslmode=disable&pool_max_conns=10"
	// url := "user=root password=msdnmm host=103.71.69.244 port=5432 dbname=bank sslmode=disable pool_max_conns=10"
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		os.Exit(1)
	}

	return conn
}

func TestSqlc(t *testing.T) {
	conn := Setup()
    defer conn.Close()

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

package query

import (
	"backend/internal/data/db/gen"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

// Store 扩展db, 因为sqlc生成的查询不支持事务
type Store struct {
	*gen.Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store  {
	return &Store{
		db: db,
		Queries: gen.New(db),
	}
}

func Connect() *pgxpool.Pool {
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

// TransferTx 通用的事务
func TransferTx(ctx context.Context, fn func(*gen.Queries) error) error  {
	conn := Connect()
	defer conn.Close()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	// 创建一个新的事务
	q := gen.New(tx)
	err = fn(q)
	if err != nil {
		// 如果回滚与事务同时失败, 则返回2错误
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("回滚错误: %v, 事务错误: %v", rbErr, err)
		}
		// 返回事务错误
		return err
	}

	// 返回错误给调用者
	return tx.Commit(ctx)
}

type TransferParams struct {
	FromAccountID int64 `json:"from_account_id,omitempty"`
	ToAccountID   int64 `json:"to_account_id,omitempty"`
	Amount        int64 `json:"amount,omitempty"`
}

// TransferResult 包含转账交易的结果
type TransferResult struct {
	// 创建的转账数据
	Transfer      *gen.Transfer `json:"transfer,omitempty"`
	// 转账后的余额
	FromAccountID *gen.Account `json:"from_account_id,omitempty"`
	// 转给的用户
	ToAccountID   *gen.Account `json:"to_account_id,omitempty"`
	// 记录支出后的用户余额
	FromEntry *gen.Entry `json:"from_entry,omitempty"`
	// 记录收入的用户余额
	ToEntry   *gen.Entry `json:"to_entry,omitempty"`
}

func (s *Store) Transaction(ctx context.Context, arg TransferParams) (TransferResult,error) {
	var result TransferResult
	// 在一个新的事务里面操作, 使用q来操作事务, 而不是db
	err := TransferTx(ctx, func(q *gen.Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, &gen.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, &gen.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, &gen.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result ,err
}

package query

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"user/internal/biz/repository"
)

// Store 扩展db, 因为sqlc生成的查询不支持事务
type Store struct {
	*repository.Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: repository.New(db),
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
func TransferTx(ctx context.Context, fn func(*repository.Queries) error) error {
	conn := Connect()
	defer conn.Close()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	// 创建一个新的事务
	q := repository.New(tx)
	err = fn(q)
	if err != nil {
		// 如果回滚与事务同时失败, 则返回2错误
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("回滚错误: %v, 事务错误: %v\n", rbErr, err)
		}
		// 返回事务错误
		return err
	}

	// 返回错误给调用者
	return tx.Commit(ctx)
}

type TransferParams struct {
	FromAccount int64 `json:"from_account_id,omitempty"`
	ToAccount   int64 `json:"to_account_id,omitempty"`
	Amount        int64 `json:"amount,omitempty"`
}

// TransferResult 包含转账交易的结果
type TransferResult struct {
	// 创建的转账数据
	Transfer *repository.Transfer `json:"transfer,omitempty"`
	// 转账后的余额
	FromAccount *repository.Account `json:"from_account_id,omitempty"`
	// 转给的用户
	ToAccount *repository.Account `json:"to_account_id,omitempty"`
	// 记录支出后的用户余额
	FromEntry *repository.Entry `json:"from_entry,omitempty"`
	// 记录收入的用户余额
	ToEntry *repository.Entry `json:"to_entry,omitempty"`
}

var txKey2 = struct{}{}

func (s *Store) Transaction(ctx context.Context, arg TransferParams) (TransferResult, error) {
	var result TransferResult
	// 在一个新的事务里面操作, 使用q来操作事务, 而不是db
	err := TransferTx(ctx, func(q *repository.Queries) error {
		var err error

		txName := ctx.Value(txKey2)

		fmt.Printf("%v: create transfer 创建转账记录, 由用户ID: %v 转到用户ID: %v 的金额是: %v\n",txName, arg.FromAccount, arg.ToAccount, arg.Amount)
		result.Transfer, err = q.CreateTransfer(ctx, &repository.CreateTransferParams{
			FromAccountID: arg.FromAccount,
			ToAccountID:   arg.ToAccount,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%v: create entry 1 创建账单记录, 记录用户: %v 的金额变化为: %v\n", txName,arg.FromAccount, -arg.Amount)
		result.FromEntry, err = q.CreateEntry(ctx, &repository.CreateEntryParams{
			AccountID: arg.FromAccount,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%v: create entry 2 创建账单记录, 记录用户: %v 的金额变化为: %v\n", txName, arg.ToAccount, arg.Amount)
		result.ToEntry, err = q.CreateEntry(ctx, &repository.CreateEntryParams{
			AccountID: arg.ToAccount,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%v: get account 1\n", txName)
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccount)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, &repository.UpdateAccountParams{
			ID:       arg.FromAccount,
			Balance:  account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%v: get account 2\n", txName)
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccount)
		if err != nil {
			return err
		}

		fmt.Printf("%v: update account 2 更新账户, 记录用户ID: %v 的余额变化为: %v\n",txName,  arg.ToAccount, arg.Amount)
		result.ToAccount, err = q.UpdateAccount(ctx, &repository.UpdateAccountParams{
			ID:       arg.ToAccount,
			Balance:  account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

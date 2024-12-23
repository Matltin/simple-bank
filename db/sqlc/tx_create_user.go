package db

import "context"

type CreateUserTxParams struct {
	CreatUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreatUser(ctx, arg.CreatUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}

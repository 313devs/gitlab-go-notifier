package commit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/313devs/gitlab-go-notifier/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}


func (r *RedisRepo) Insert(ctx context.Context, commit model.Commit) error {
	data, err := json.Marshal(commit)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, commit.Sha, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "commits", commit.Sha).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

var ErrNotExists = errors.New("no commits available")

type GetAllPage struct {
	size   uint64
	offset uint64
}
type FindResult struct {
	Commits []model.Commit
	Cursor  uint64
}

func (r *RedisRepo) GetAll(ctx context.Context, page GetAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "commits", page.offset, "*", int64(page.size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get commits: %w", err)
	}
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get commits: %w", err)
	}

	commits := make([]model.Commit, 0, len(xs))

	for i, x := range xs {
		x := x.(string)
		var commit model.Commit

		err := json.Unmarshal([]byte(x), &commit)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to unmarshal commit: %w", err)
		}
		commits[i] = commit
	}
	return FindResult{
		Commits: commits,
		Cursor:  cursor,
	}, nil
}

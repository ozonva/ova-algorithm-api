package repo

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
)

const (
	tableName         = "algorithms"
	idColumn          = "id"
	subjectColumn     = "subject"
	descriptionColumn = "description"
)

type Repo interface {

	// AddAlgorithms adds new entities to the repo. Ids are assigned by
	// the store
	AddAlgorithms(ctx context.Context, algorithm []algorithm.Algorithm) error

	// ListAlgorithms return list of entities provided limit and offset
	ListAlgorithms(ctx context.Context, limit, offset uint64) ([]algorithm.Algorithm, error)

	// DescribeAlgorithm returns entity details for provided algorithmID
	DescribeAlgorithm(ctx context.Context, algorithmID uint64) (*algorithm.Algorithm, error)

	// RemoveAlgorithm returns found id entity has been removed and error
	RemoveAlgorithm(ctx context.Context, algorithmID uint64) (bool, error)

	// UpdateAlgorithm updates fields of algorithm. Algorithm is selected
	// provided id. If no algorithm exists nothing is updates and false is
	// returned as the first return value
	UpdateAlgorithm(ctx context.Context, algorithm algorithm.Algorithm) (bool, error)
}

// NewRepo creates new Repo with provided database connection
func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}

type repo struct {
	db *sql.DB
}

func (r *repo) AddAlgorithms(ctx context.Context, algorithms []algorithm.Algorithm) error {
	if len(algorithms) == 0 {
		return nil
	}

	sql, _, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).Columns(subjectColumn, descriptionColumn).
		Suffix("RETURNING id").
		Values("", "").ToSql()

	if err != nil {
		return fmt.Errorf("failed to build sql template: %w", err)
	}


	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failer to start transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		retErr := fmt.Errorf("failed to build sql template: %w", err)

		if err := tx.Rollback(); err != nil {
			retErr = fmt.Errorf("cannot rollback: %w", err)
		}

		return retErr
	}
	defer stmt.Close()

	for i := 0; i < len(algorithms); i++ {
		id, err := addAlgorithmQuery(ctx, stmt, algorithms[i])
		if err != nil {
			retErr := fmt.Errorf("cannot add algorithm: %w", err)
			if err := tx.Rollback(); err != nil {
				retErr = fmt.Errorf("cannot rollback: %w", err)
			}
			return retErr
		}

		algorithms[i].UserID = id
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func addAlgorithmQuery(ctx context.Context, stmt *sql.Stmt, a algorithm.Algorithm) (uint64, error) {
	idsSQL, err := stmt.QueryContext(ctx, a.Subject, a.Description)
	if err != nil {
		return 0, fmt.Errorf("cannot execute prepared statement: %w", err)
	}
	defer idsSQL.Close()

	if !idsSQL.Next() {
		return 0, fmt.Errorf("no id returned: %w", idsSQL.Err())
	}

	var id uint64
	if err := idsSQL.Scan(&id); err != nil {
		return 0, fmt.Errorf("cannot parse sql row: %w", err)
	}

	return id, nil
}

func (r *repo) ListAlgorithms(ctx context.Context, limit, offset uint64) ([]algorithm.Algorithm, error) {
	users, err := sq.Select("*").
		From(tableName).
		OrderBy(idColumn).
		Limit(limit).
		Offset(offset).
		RunWith(r.db).QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("cannot run list query: %w", err)
	}
	defer users.Close()

	algorithms := make([]algorithm.Algorithm, 0, limit)

	for users.Next() {
		var algo algorithm.Algorithm
		if err := users.Scan(&algo.UserID, &algo.Subject, &algo.Description); err != nil {
			return algorithms, fmt.Errorf("cannot parse algortihm: %w", err)
		}
		algorithms = append(algorithms, algo)
	}
	if err := users.Err(); err != nil {
		return algorithms, fmt.Errorf("error list query %w", err)
	}

	log.Debug().Int("algorithmsLen", len(algorithms)).Msg("ListAlgorithms")

	return algorithms, nil
}

func (r *repo) DescribeAlgorithm(ctx context.Context, algorithmID uint64) (*algorithm.Algorithm, error) {
	users, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From(tableName).
		Where(sq.Eq{idColumn: algorithmID}).
		RunWith(r.db).
		QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("cannot run describe query: %w", err)
	}
	defer users.Close()

	if !users.Next() {
		return nil, nil
	}

	var algo = &algorithm.Algorithm{}

	if err := users.Scan(&algo.UserID, &algo.Subject, &algo.Description); err != nil {
		return nil, fmt.Errorf("cannot parse algortihm: %w", err)
	}

	if err := users.Err(); err != nil {
		return nil, fmt.Errorf("error list query %w", err)
	}

	return algo, nil
}

func (r *repo) RemoveAlgorithm(ctx context.Context, algorithmID uint64) (bool, error) {
	result, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("").
		From(tableName).
		Where(sq.Eq{idColumn: algorithmID}).
		RunWith(r.db).
		ExecContext(ctx)

	if err != nil {
		return false, fmt.Errorf("cannot run delete query: %w", err)
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("cannot get rows affected: %w", err)
	}

	return deletedRows > 0, nil
}

func (r *repo) UpdateAlgorithm(ctx context.Context, algorithm algorithm.Algorithm) (bool, error) {
	result, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tableName).
		Set(subjectColumn, algorithm.Subject).
		Set(descriptionColumn, algorithm.Description).
		Where(sq.Eq{idColumn: algorithm.UserID}).
		RunWith(r.db).
		ExecContext(ctx)

	if err != nil {
		return false, fmt.Errorf("cannot run delete query: %w", err)
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("cannot get rows affected: %w", err)
	}

	return updatedRows > 0, nil
}

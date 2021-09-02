package repo

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/rs/zerolog/log"
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
	AddAlgorithms(algorithm []algorithm.Algorithm) error

	// ListAlgorithms return list of entities provided limit and offset
	ListAlgorithms(limit, offset uint64) ([]algorithm.Algorithm, error)

	// DescribeAlgorithm returns entity details for provided algorithmID
	DescribeAlgorithm(algorithmID uint64) (*algorithm.Algorithm, error)

	// RemoveAlgorithm returns found id entity has been removed and error
	RemoveAlgorithm(algorithmID uint64) (bool, error)
}

func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}

type repo struct {
	db *sql.DB
}

func (r *repo) AddAlgorithms(algorithms []algorithm.Algorithm) error {
	if len(algorithms) == 0 {
		return nil
	}

	sql, _, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).Columns(subjectColumn, descriptionColumn).
		Values("", "").ToSql()

	if err != nil {
		return fmt.Errorf("failed to build sql template: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failer to start transaction: %w", err)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to build sql template: %w", err)
	}
	defer stmt.Close()

	for i := 0; i < len(algorithms); i++ {
		if _, err := stmt.Exec(algorithms[i].Subject, algorithms[i].Description); err != nil {
			tx.Rollback()
			return fmt.Errorf("cannot fill prepared statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (r *repo) ListAlgorithms(limit, offset uint64) ([]algorithm.Algorithm, error) {
	users, err := sq.Select("*").
		From(tableName).
		OrderBy(idColumn).
		Limit(limit).
		Offset(offset).
		RunWith(r.db).Query()

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

func (r *repo) DescribeAlgorithm(algorithmID uint64) (*algorithm.Algorithm, error) {
	users, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From(tableName).
		Where(sq.Eq{idColumn: algorithmID}).
		RunWith(r.db).
		Query()

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

func (r *repo) RemoveAlgorithm(algorithmID uint64) (bool, error) {
	result, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("").
		From(tableName).
		Where(sq.Eq{idColumn: algorithmID}).
		RunWith(r.db).
		Exec()

	if err != nil {
		return false, fmt.Errorf("cannot run delete query: %w", err)
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("cannot get rows affected: %w", err)
	}

	return deletedRows > 0, nil
}

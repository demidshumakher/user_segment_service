package postgresql

import (
	"context"
	"database/sql"
	"math"
	"segment_service/domain"
)

// PostgresSegmentRepository репозиторий для работы с сегментами в PostgreSQL.
type PostgresSegmentRepository struct {
	db *sql.DB
}

// NewSegmentRepository создаёт новый экземпляр PostgresSegmentRepository.
func NewSegmentRepository(db *sql.DB) *PostgresSegmentRepository {
	return &PostgresSegmentRepository{db: db}
}

// GetAll возвращает список всех сегментов.
func (r *PostgresSegmentRepository) GetAll() []domain.Segment {
	ctx := context.Background()

	rows, err := r.db.QueryContext(ctx, `SELECT name FROM segments ORDER BY name`)
	if err != nil {
		return []domain.Segment{}
	}
	defer rows.Close()

	result := make([]domain.Segment, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		result = append(result, domain.Segment(name))
	}
	return result
}

// Delete удаляет сегмент.
func (r *PostgresSegmentRepository) Delete(segment string) error {
	ctx := context.Background()
	_, err := r.db.ExecContext(ctx, `DELETE FROM segments WHERE name = $1`, segment)
	return err
}

// Create создаёт сегмент.
func (r *PostgresSegmentRepository) Create(segment string) error {
	ctx := context.Background()
	_, err := r.db.ExecContext(ctx, `INSERT INTO segments (name) VALUES ($1)`, segment)
	return err
}

// Distribute добавляет указанный сегмент случайно выбранному проценту пользователей, не имеющих его.
func (r *PostgresSegmentRepository) Distribute(segment string, percentage float64) error {
	ctx := context.Background()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Считаем количество пользователей, у которых нет этого сегмента
	var total int
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM users u
		LEFT JOIN user_segments us ON us.user_id = u.id AND us.segment = $1
		WHERE us.user_id IS NULL
	`, segment).Scan(&total)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	toAssign := int(math.Round(float64(total) * percentage / 100.0))
	if toAssign <= 0 {
		return tx.Commit()
	}

	// Вставляем сегмент для случайно выбранных пользователей, у которых его нет
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_segments (user_id, segment)
		SELECT u.id, $1
		FROM users u
		LEFT JOIN user_segments us ON us.user_id = u.id AND us.segment = $1
		WHERE us.user_id IS NULL
		ORDER BY random()
		LIMIT $2
	`, segment, toAssign)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

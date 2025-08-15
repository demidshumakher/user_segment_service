package postgresql

import (
	"context"
	"database/sql"
	"segment_service/domain"
)

// PostgresUserRepository репозиторий для работы с пользователями и их сегментами в PostgreSQL.
type PostgresUserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт новый экземпляр PostgresUserRepository.
func NewUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// GetAll возвращает карту пользователей и их сегментов.
// Возвращаются все пользователи, включая тех, у которых нет связанных сегментов.
func (r *PostgresUserRepository) GetAll() (map[domain.User][]domain.Segment, error) {
	ctx := context.Background()

	// Сначала получаем всех пользователей
	usersRows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT id FROM users
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer usersRows.Close()

	// Инициализируем карту для всех пользователей
	result := make(map[domain.User][]domain.Segment)
	for usersRows.Next() {
		var uid int
		if err := usersRows.Scan(&uid); err != nil {
			return nil, err
		}
		userKey := domain.User(uid)
		// Инициализируем пустой список сегментов для пользователя
		result[userKey] = []domain.Segment{}
	}
	if err := usersRows.Err(); err != nil {
		return nil, err
	}

	// Теперь получаем все сегменты пользователей
	segmentsRows, err := r.db.QueryContext(ctx, `
		SELECT us.user_id, s.name
		FROM user_segments us
		JOIN segments s ON s.name = us.segment
		ORDER BY us.user_id, s.name
	`)
	if err != nil {
		return nil, err
	}
	defer segmentsRows.Close()

	// Заполняем сегменты для соответствующих пользователей
	for segmentsRows.Next() {
		var uid int
		var sname string
		if err := segmentsRows.Scan(&uid, &sname); err != nil {
			return nil, err
		}
		userKey := domain.User(uid)
		result[userKey] = append(result[userKey], domain.Segment(sname))
	}
	if err := segmentsRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSegmentsById возвращает список сегментов пользователя.
func (r *PostgresUserRepository) GetSegmentsById(id int) ([]domain.Segment, error) {
	ctx := context.Background()

	rows, err := r.db.QueryContext(ctx, `
		SELECT s.name
		FROM user_segments us
		JOIN segments s ON s.name = us.segment
		WHERE us.user_id = $1
		ORDER BY s.name
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	segments := make([]domain.Segment, 0)
	for rows.Next() {
		var sname string
		if err := rows.Scan(&sname); err != nil {
			return nil, err
		}
		segments = append(segments, domain.Segment(sname))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}

// ClearUserSegments удаляет все сегменты у пользователя.
func (r *PostgresUserRepository) ClearUserSegments(id int) error {
	ctx := context.Background()
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_segments WHERE user_id = $1`, id)
	return err
}

// AddSegment добавляет сегмент пользователю.
func (r *PostgresUserRepository) AddSegment(id int, segment string) error {
	ctx := context.Background()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_segments (user_id, segment)
		VALUES ($1, $2)
		ON CONFLICT (user_id, segment) DO NOTHING
	`, id, segment)
	return err
}

// DeleteSegment удаляет сегмент у пользователя.
func (r *PostgresUserRepository) DeleteSegment(id int, segment string) error {
	ctx := context.Background()
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_segments WHERE user_id = $1 AND segment = $2`, id, segment)
	return err
}

package users

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stanislavCasciuc/atom-fit-go/internal/types"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	const op = "user.store.GetUserByEmail"

	rows, err := s.db.Queryx("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var u *types.User
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	const op = "user.store.GetUserByID"

	rows, err := s.db.Queryx("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var u *types.User
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return u, nil
}

func (s *Store) CreateUser(u types.RegisterUserPayload) (int, error) {
	const op = "user.store.CreateUser"

	stmt, err := s.db.Preparex(
		"INSERT INTO users(email, username, password, is_male, age, height, weight, goal, weight_goal) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int

	err = stmt.QueryRow(
		u.Email, u.Username, u.Password, u.IsMale, u.Age, u.Height, u.Weight, u.Goal, u.WeightGoal,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func scanRowIntoUser(rows *sqlx.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.IsActive,
		&user.IsSuperuser,
		&user.IsMale,
		&user.Age,
		&user.Height,
		&user.Weight,
		&user.Goal,
		&user.WeightGoal,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

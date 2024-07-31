package users

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stanislavCasciuc/atom-fit-go/internal/services/users/models"
)

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(useData models.RegisterUserPayload) (int, error)
}

var UserAlreadyExist = errors.New("user already exists")

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	const op = "user.store.GetUserByEmail"

	rows, err := s.db.Queryx("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var u *models.User
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*models.User, error) {
	const op = "user.store.GetUserByID"

	rows, err := s.db.Queryx("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var u *models.User
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return u, nil
}

func (s *Store) CreateUser(u models.RegisterUserPayload) (int, error) {
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
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, UserAlreadyExist
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func scanRowIntoUser(rows *sqlx.Rows) (*models.User, error) {
	user := new(models.User)

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

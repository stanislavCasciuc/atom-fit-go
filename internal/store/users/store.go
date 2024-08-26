package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stanislavCasciuc/atom-fit-go/internal/models"
)

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(userData models.RegisterUserPayload, passwordHash []byte) (int, string, error)
	UpdateUser(id int, userData models.User) error
}

var (
	UserAlreadyExist = errors.New("users already exists")
	UserNotFound     = errors.New("users not found")
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	const op = "users.store.GetUserByEmail"

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
	const op = "users.store.GetUserByID"

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
	if u.ID == 0 {
		return nil, fmt.Errorf("%s: %w", op, UserNotFound)
	}
	return u, nil
}

func (s *Store) CreateUser(u models.RegisterUserPayload, passHash []byte) (int, string, error) {
	const op = "users.store.CreateUser"

	stmt, err := s.db.Preparex(
		"INSERT INTO users(email, username, password, is_male, age, height, weight, goal, weight_goal, activation_code) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
	)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	var id int
	activationCode := uuid.New()

	err = stmt.QueryRow(
		u.Email, u.Username, passHash, u.IsMale, u.Age, u.Height, u.Weight, u.Goal, u.WeightGoal, activationCode,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, "", UserAlreadyExist
		}
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}
	return id, activationCode.String(), nil
}

func (s *Store) UpdateUser(id int, userData models.User) error {
	const op = "users.store.UpdateUser"

	// Prepare the SQL statement
	stmt, err := s.db.Preparex(
		"UPDATE users SET email = $1, username = $2, is_male = $3, age = $4, height = $5, weight = $6, goal = $7, weight_goal = $8, is_active = $9 WHERE id = $10",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Execute the SQL statement
	_, err = stmt.Exec(
		userData.Email, userData.Username, userData.IsMale, userData.Age, userData.Height, userData.Weight,
		userData.Goal, userData.WeightGoal, userData.IsActive, id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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
		&user.ActivationCode,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

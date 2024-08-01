package models

import "time"

type User struct {
	ID          int       `db:"id"`
	Email       string    `db:"email"`
	Username    string    `db:"username"`
	Password    []byte    `db:"password"`
	CreatedAt   time.Time `db:"created_at"`
	IsActive    bool      `db:"is_active"`
	IsSuperuser bool      `db:"is_superuser"`
	IsMale      bool      `db:"is_male"`
	Age         int       `db:"age"`
	Height      int       `db:"height"`
	Weight      int       `db:"weight"`
	Goal        string    `db:"goal"`
	WeightGoal  int       `db:"weight_goal"`
}

type RegisterUserPayload struct {
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required,min=3,max=30"`
	IsMale     bool   `json:"isMale" validate:"required"`
	Age        int    `json:"age" validate:"required"`
	Height     int    `json:"height" validate:"required"`
	Weight     int    `json:"weight" validate:"required"`
	Goal       string `json:"goal" validate:"required,oneof=lose maintain gain"`
	WeightGoal int    `json:"weightGoal" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=30"`
}

package models

type Cat struct {
	ID                int     `db:"id" json:"ID"`
	Name              string  `db:"name" json:"Name" validate:"required"`
	YearsOfExperience int     `db:"years_of_experience" json:"YearsOfExperience" validate:"required"`
	Breed             string  `db:"breed" json:"Breed" validate:"required"`
	Salary            float64 `db:"salary" json:"Salary" validate:"required"`
	CreatedAt         string  `db:"created_at" json:"CreatedAt"`
	UpdatedAt         string  `db:"updated_at" json:"UpdatedAt,omitempty"`
}

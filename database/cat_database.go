package database

import (
	"database/sql"
	"spyCat/database/models"
	"time"
)

type CatDatabaseInterface interface {
	SelectAll() (*[]models.Cat, error)
	SelectByID(id int) (*models.Cat, error)
	Insert(cat models.Cat) (int, error)
	Update(catID int, salary float64) error
	Delete(id int) error
}

type CatDatabase struct {
	*Database
}

func NewCatDatabase(Conn *Database) *CatDatabase {
	return &CatDatabase{Conn}
}

func (cd *CatDatabase) SelectAll() (*[]models.Cat, error) {
	var createdAt, updatedAt time.Time
	var cats []models.Cat

	query := `SELECT id, name, years_of_experience, breed, salary, created_at, updated_at FROM spy_cats`
	rows, err := cd.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat models.Cat

		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		cat.CreatedAt = createdAt.Format("15:04:05 02:01:06")
		cat.UpdatedAt = updatedAt.Format("15:04:05 02:01:06")
		cats = append(cats, cat)
	}

	return &cats, nil
}

func (cd *CatDatabase) SelectByID(id int) (*models.Cat, error) {
	var cat models.Cat
	var createdAt, updatedAt time.Time

	query := `SELECT id, name, years_of_experience, breed, salary, created_at, updated_at 
              FROM spy_cats WHERE id = $1`
	err := cd.Connection.QueryRow(query, id).Scan(
		&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	cat.CreatedAt = createdAt.Format("15:04:05 02:01:06")
	cat.UpdatedAt = updatedAt.Format("15:04:05 02:01:06")

	return &cat, nil
}

func (cd *CatDatabase) Insert(cat models.Cat) (int, error) {
	var id int
	query := `INSERT INTO spy_cats (name, years_of_experience, breed, salary) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := cd.Connection.QueryRow(query, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cd *CatDatabase) Update(catID int, salary float64) error {
	query := `UPDATE spy_cats SET salary = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	result, err := cd.Connection.Exec(query, salary, catID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (cd *CatDatabase) Delete(catID int) error {
	query := `DELETE FROM spy_cats WHERE id = $1`
	result, err := cd.Connection.Exec(query, catID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

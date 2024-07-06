package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"spyCat/database"
	"spyCat/database/models"
	"strings"
	"time"
)

type CatService struct {
	DbCat    database.CatDatabaseInterface
	validate *validator.Validate
}

func NewCatService(DbCat database.CatDatabaseInterface, validate *validator.Validate) *CatService {
	return &CatService{DbCat: DbCat, validate: validate}
}

type CatServiceInterface interface {
	GetAllCats() (*[]models.Cat, error)
	GetCat(catID int) (*models.Cat, error, int)
	CreateCat(cat models.Cat) (int, error, int)
	EditCatSalary(ID int, salary float64) (error, int)
	DeleteCat(catID int) (error, int)
	CatValidation(cat models.Cat) error
}

type CatBreed struct {
	Name string `json:"name"`
}

var (
	catAPIURL     = "https://api.thecatapi.com/v1/breeds"
	httpClient    = &http.Client{Timeout: 10 * time.Second}
	cachedBreeds  []CatBreed
	cacheExpireAt time.Time
)

func (cs *CatService) CreateCat(cat models.Cat) (int, error, int) {
	if !isValidBreed(cat.Breed) {
		return 0, errors.New("invalid breed: the provided breed does not match any known cat breeds"), http.StatusBadRequest
	}

	newCat := models.Cat{
		Name:              cat.Name,
		YearsOfExperience: cat.YearsOfExperience,
		Breed:             cat.Breed,
		Salary:            cat.Salary,
	}

	if err := cs.CatValidation(newCat); err != nil {
		return 0, err, http.StatusBadRequest
	}

	insertedId, err := cs.DbCat.Insert(newCat)
	if err != nil {
		return insertedId, err, http.StatusInternalServerError
	}

	return insertedId, err, http.StatusCreated
}

func (cs *CatService) GetCat(catID int) (*models.Cat, error, int) {
	cat, err := cs.DbCat.SelectByID(catID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("there is no cat with that ID"), http.StatusBadRequest
	} else if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return cat, nil, http.StatusOK
}

func (cs *CatService) EditCatSalary(ID int, salary float64) (error, int) {
	err := cs.DbCat.Update(ID, salary)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return err, http.StatusOK
}

func (cs *CatService) GetAllCats() (*[]models.Cat, error) {
	cats, err := cs.DbCat.SelectAll()
	if err != nil {
		return nil, err
	}

	return cats, err
}

func (cs *CatService) DeleteCat(catID int) (error, int) {
	err := cs.DbCat.Delete(catID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("there is no user with that ID"), http.StatusBadRequest
	} else if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (cs *CatService) CatValidation(cat models.Cat) error {
	if validationErr := cs.validate.Struct(&cat); validationErr != nil {
		return validationErr
	}
	return nil
}

func fetchBreeds() ([]CatBreed, error) {
	if time.Now().Before(cacheExpireAt) && len(cachedBreeds) > 0 {
		return cachedBreeds, nil
	}

	req, err := http.NewRequest("GET", catAPIURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch cat breeds")
	}

	var breeds []CatBreed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}

	cachedBreeds = breeds
	cacheExpireAt = time.Now().Add(24 * time.Hour)

	return breeds, nil
}

func isValidBreed(breed string) bool {
	breeds, err := fetchBreeds()
	if err != nil {
		// Log the error and return false, or handle it as appropriate for your application
		return false
	}

	normalizedUserBreed := strings.ToLower(strings.TrimSpace(breed))
	for _, apiBreed := range breeds {
		if strings.ToLower(apiBreed.Name) == normalizedUserBreed {
			return true
		}
	}

	return false
}

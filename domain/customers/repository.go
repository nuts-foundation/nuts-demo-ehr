package customers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(id string) (*domain.Customer, error)
	All() ([]domain.Customer, error)
}

type flatFileRepo struct {
	filepath string
	mutex    sync.Mutex
	// records is a cache
	records map[string]domain.Customer
}

func NewFlatFileRepository(filepath string) *flatFileRepo {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	return &flatFileRepo{
		filepath: filepath,
		mutex:    sync.Mutex{},
		records:  make(map[string]domain.Customer, 0),
	}
}

func (db *flatFileRepo) FindByID(id string) (*domain.Customer, error) {
	if len(db.records) == 0 {
		if err := db.ReadAll(); err != nil {
			return nil, err
		}
	}

	for _, r := range db.records {
		if r.Id == id {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	return nil, errors.New("not found")
}

func (db *flatFileRepo) ReadAll() error {
	//log.Debug("Reading full customer list from file")
	bytes, err := os.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("unable to read db from file: %w", err)
	}

	if len(bytes) == 0 {
		return nil
	}

	if err = json.Unmarshal(bytes, &db.records); err != nil {
		return fmt.Errorf("unable to unmarshall db from file: %w", err)
	}
	return nil
}

func (db *flatFileRepo) All() ([]domain.Customer, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if err := db.ReadAll(); err != nil {
		return nil, err
	}

	v := make([]domain.Customer, len(db.records))

	idx := 0
	for _, value := range db.records {
		v[idx] = value
		idx = idx + 1
	}
	sort.SliceStable(v, func(i, j int) bool {
		return v[i].Id < v[j].Id
	})
	return v, nil
}

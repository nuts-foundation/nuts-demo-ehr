package customers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"sort"
	"sync"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Repository interface {
	FindByID(id int) (*types.Customer, error)
	FindByDID(did string) (*types.Customer, error)
	All() ([]types.Customer, error)
}

type jsonFileRepo struct {
	filepath string
	mutex    sync.Mutex
	// records is a cache
	records map[string]types.Customer
}

func NewJsonFileRepository(filepath string) *jsonFileRepo {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Warnf("Could not open cusomers file (path=%s), it still needs to be created: %s", filepath, err)
		// But allow to proceed, since it is shared with nuts-registry-admin-demo, which creates it.
		// In Docker environments, it might not be there yet if demo-ehr starts first.
	}

	repo := jsonFileRepo{
		filepath: filepath,
		mutex:    sync.Mutex{},
		records:  make(map[string]types.Customer, 0),
	}

	return &repo
}

func (db *jsonFileRepo) FindByID(id int) (*types.Customer, error) {
	if len(db.records) == 0 {
		if err := db.readAll(); err != nil {
			return nil, err
		}
	}

	for _, r := range db.records {
		if r.Id == id {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	// Not found
	return nil, nil
}

func (db *jsonFileRepo) FindByDID(did string) (*types.Customer, error) {
	if len(db.records) == 0 {
		if err := db.readAll(); err != nil {
			return nil, err
		}
	}

	for _, r := range db.records {
		if r.Did != nil && *r.Did == did {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	// Not found
	return nil, nil
}

func (db *jsonFileRepo) readAll() error {
	bytes, err := os.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("unable to read db from file: %w", err)
	}

	if len(bytes) == 0 {
		return nil
	}

	v := map[string]types.Customer{}
	if err = json.Unmarshal(bytes, &v); err != nil {
		return fmt.Errorf("unable to unmarshall db from file: %w", err)
	}

	db.records = v

	return nil
}

func (db *jsonFileRepo) All() ([]types.Customer, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if err := db.readAll(); err != nil {
		return nil, err
	}

	v := make([]types.Customer, len(db.records))

	idx := 0
	for _, value := range db.records {
		v[idx] = value
		idx = idx + 1
	}
	sort.SliceStable(v, func(i, j int) bool {
		return v[i].Name < v[j].Name
	})
	return v, nil
}

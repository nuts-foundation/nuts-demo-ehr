package customers

import (
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
)

type Repository interface {
	FindByID(id int) (*types.Customer, error)
	FindByDID(did string) (*types.Customer, error)
	All() ([]types.Customer, error)
}

type jsonFileRepo struct {
	filepath string
}

func NewJsonFileRepository(filepath string) Repository {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		logrus.Warnf("Could not open cusomers file (path=%s), it still needs to be created: %s", filepath, err)
		// But allow to proceed, since it is shared with nuts-registry-admin-demo, which creates it.
		// In Docker environments, it might not be there yet if demo-ehr starts first.
	} else {
		defer f.Close()
	}

	return &jsonFileRepo{
		filepath: filepath,
	}
}

func (db *jsonFileRepo) FindByID(id int) (*types.Customer, error) {
	records, err := db.read()
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		if r.Id == id {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	// Not found
	return nil, nil
}

func (db *jsonFileRepo) FindByDID(did string) (*types.Customer, error) {
	records, err := db.read()
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		if r.Did != nil && *r.Did == did {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	// Not found
	return nil, nil
}

func (db *jsonFileRepo) All() ([]types.Customer, error) {
	records, err := db.read()
	if err != nil {
		return nil, err
	}

	v := make([]types.Customer, len(records))

	idx := 0
	for _, value := range records {
		v[idx] = value
		idx = idx + 1
	}
	sort.SliceStable(v, func(i, j int) bool {
		return v[i].Name < v[j].Name
	})
	return v, nil
}

func (db *jsonFileRepo) read() (map[string]types.Customer, error) {
	bytes, err := os.ReadFile(db.filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read db from file: %w", err)
	}

	if len(bytes) == 0 {
		return nil, nil
	}

	v := map[string]types.Customer{}
	if err = json.Unmarshal(bytes, &v); err != nil {
		return nil, fmt.Errorf("unable to unmarshall db from file: %w", err)
	}

	return v, nil
}

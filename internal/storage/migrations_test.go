package storage

import (
	"github.com/golang-migrate/migrate/v4"
	"os"
	"regexp"
	"slices"
	"strconv"
	"testing"
)

func getMigrationRevisions() ([]uint, error) {
	files, err := os.ReadDir("./migrations")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(\d+)_.*\.(up|down)\.sql$`)
	revisionSet := make(map[string]struct{})

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		match := re.FindStringSubmatch(file.Name())
		if len(match) > 1 {
			revisionSet[match[1]] = struct{}{}
		}
	}

	var revisions []uint
	for rev := range revisionSet {
		revNum, _ := strconv.Atoi(rev)
		revisions = append(revisions, uint(revNum))
	}

	slices.Sort(revisions)
	return revisions, nil
}

// TestStairway verifies the integrity of each migration revision by performing the following steps:
// 1. Applies the migration up to the given revision.
// 2. Rolls back a single step (Steps(-1)).
// 3. Reapplies the same revision.
// 4. Fully rolls back all applied migrations (Down).
func TestStairway(t *testing.T) {
	revisions, err := getMigrationRevisions()
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrator()
	if err != nil {
		t.Fatal(err)
	}

	defer func(m *migrate.Migrate) {
		err, err2 := m.Close()
		if err != nil {
			t.Fatal(err)
		}
		if err2 != nil {
			t.Fatal(err2)
		}
	}(m)

	for _, rev := range revisions {
		err = m.Migrate(rev)
		if err != nil {
			t.Fatal(err)
		}

		err = m.Steps(-1)
		if err != nil {
			t.Fatal(err)
		}

		err = m.Migrate(rev)
		if err != nil {
			t.Fatal(err)
		}

		err = m.Down()
		if err != nil {
			t.Fatal(err)
		}
	}
}

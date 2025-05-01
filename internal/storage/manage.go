package storage

import (
	"context"
	"fmt"
)

func createDB(name, template string) error {
	query := fmt.Sprintf(
		`CREATE DATABASE "%s" ENCODING 'UTF8' TEMPLATE "%s";`,
		name,
		template,
	)
	_, err := Pool().Exec(context.Background(), query)

	return err
}

func dropDB(name string) error {
	query := fmt.Sprintf(`DROP DATABASE "%s";`, name)
	_, err := Pool().Exec(context.Background(), query)

	return err
}

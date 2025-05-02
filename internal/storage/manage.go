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
	_, err := activePool().Exec(context.Background(), query)

	return err
}

func dropDB(name string) error { // coverage-ignore
	// (it's actually tested, but cover don't catch it...)
	query := fmt.Sprintf(`DROP DATABASE "%s";`, name)
	_, err := activePool().Exec(context.Background(), query)

	return err
}

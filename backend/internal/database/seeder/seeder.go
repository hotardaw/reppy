package seeder

import (
	"context"
	"fmt"
	"go-fitsync/backend/internal/database/seeder/seeds"
	"go-fitsync/backend/internal/database/sqlc"
)

func cleanTestData(queries *sqlc.Queries) error {
	// delete all user profiles first bc of foreign key constraints
	err := queries.DeleteAllUserProfiles(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean user profiles: %v", err)
	}

	err = queries.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean users: %v", err)
	}
	return nil

	// add DeleteAllExercises, DeleteAllWorkouts, etc.
}

func shouldSeed(queries *sqlc.Queries) (bool, error) {
	hasSeeded, err := queries.CheckInitialSeed(context.Background())
	if err != nil {
		return false, err
	}
	return !hasSeeded, nil
}

func markSeedComplete(queries *sqlc.Queries) error {
	return queries.MarkInitialSeed(context.Background())
}

func SeedTestData(queries *sqlc.Queries) error {
	shouldDoSeed, err := shouldSeed(queries)
	if err != nil {
		return fmt.Errorf("failed to check seed status: %v", err)
	}

	if !shouldDoSeed {
		return nil // meaning seeding alreadyt happened
	}

	if err := cleanTestData(queries); err != nil {
		return fmt.Errorf("failed to clean test data: %v", err)
	}

	// below starts seed user logic
	if err = seeds.SeedUsers(queries); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	if err = seeds.SeedMuscles(queries); err != nil {
		return fmt.Errorf("failed to seed muscles: %v", err)
	}

	if err = seeds.SeedExercises(queries); err != nil {
		return fmt.Errorf("failed to seed exercises: %v", err)
	}

	// if err = seedExercisesMuscles(queries); err != nil {
	// 	return fmt.Errorf("failed to seed exercises: %v", err)
	// }

	if err := markSeedComplete(queries); err != nil {
		return fmt.Errorf("failed to mark seeding complete: %v", err)
	}

	return nil
}

package seeder

import (
	"context"
	"fmt"
	"go-reppy/backend/internal/database/seeder/seeds"
	"go-reppy/backend/internal/database/sqlc"
)

func cleanTestData(queries *sqlc.Queries) error {
	// start by deleting tables with the most FDs (workout sets), then work inwards toward those depended on most often (users)
	err := queries.DeleteAllWorkoutSetsUnconditional(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean workout-sets: %v", err)
	}

	err = queries.DeleteAllWorkouts(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean workouts: %v", err)
	}

	err = queries.DeleteAllUserProfiles(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean user-profiles: %v", err)
	}

	err = queries.DeleteAllExerciseMuscles(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean exercise-muscles: %v", err)
	}

	err = queries.DeleteAllExercises(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean exercises: %v", err)
	}

	err = queries.DeleteAllMuscles(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean muscles: %v", err)
	}

	err = queries.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to clean users: %v", err)
	}
	return nil
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

	/*
	** seeding begins
	 */

	if err = seeds.SeedUsers(queries); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	if err = seeds.SeedMuscles(queries); err != nil {
		return fmt.Errorf("failed to seed muscles: %v", err)
	}

	if err = seeds.SeedExercises(queries); err != nil {
		return fmt.Errorf("failed to seed exercises: %v", err)
	}

	if err = seeds.SeedExerciseMuscles(queries); err != nil {
		return fmt.Errorf("failed to seed exercise-muscle associations: %v", err)
	}

	if err = seeds.SeedWorkouts(queries); err != nil {
		return fmt.Errorf("failed to seed workouts: %v", err)
	}

	if err = seeds.SeedWorkoutSets(queries); err != nil {
		return fmt.Errorf("failed to seed workouts-sets: %v", err)
	}

	if err := markSeedComplete(queries); err != nil {
		return fmt.Errorf("failed to mark seeding complete: %v", err)
	}

	return nil
}

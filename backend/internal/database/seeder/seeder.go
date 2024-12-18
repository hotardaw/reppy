package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TestUser struct {
	Email     string
	Password  string
	Username  string
	FirstName string
	LastName  string
	Gender    string
	DOB       time.Time
	Height    int32
	Weight    int32
}

func GetTestUsers() []TestUser {
	return []TestUser{
		{
			Email:     "aaron@test.com",
			Password:  "password123",
			Username:  "aaronhotard",
			FirstName: "Aaron",
			LastName:  "Hotard",
			Gender:    "Male",
			DOB:       time.Date(1999, 11, 19, 0, 0, 0, 0, time.UTC),
			Height:    74,
			Weight:    185,
		},
		{
			Email:     "aiyana@test.com",
			Password:  "password456",
			Username:  "aiyanathomas",
			FirstName: "Aiyana",
			LastName:  "Thomas",
			Gender:    "Female",
			DOB:       time.Date(2000, 7, 5, 0, 0, 0, 0, time.UTC),
			Height:    63,
			Weight:    130,
		},
	}
}

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

	for _, user := range GetTestUsers() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		createdUser, err := queries.CreateUser(context.Background(), sqlc.CreateUserParams{
			Email:        user.Email,
			PasswordHash: string(hashedPassword),
			Username:     user.Username,
		})
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.Email, err)
		}

		_, err = queries.CreateUserProfile(context.Background(), sqlc.CreateUserProfileParams{
			UserID: sql.NullInt32{
				Int32: createdUser.UserID,
				Valid: true,
			},
			FirstName: sql.NullString{
				String: user.FirstName,
				Valid:  true,
			},
			LastName: sql.NullString{
				String: user.LastName,
				Valid:  true,
			},
			DateOfBirth: sql.NullTime{
				Time:  user.DOB,
				Valid: true,
			},
			Gender: sql.NullString{
				String: user.Gender,
				Valid:  true,
			},
			HeightInches: sql.NullInt32{
				Int32: user.Height,
				Valid: true,
			},
			WeightPounds: sql.NullInt32{
				Int32: user.Weight,
				Valid: true,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create profile for user %s: %v", user.Email, err)
		}
	}

	if err := markSeedComplete(queries); err != nil {
		return fmt.Errorf("failed to mark seeding complete: %v", err)
	}

	return nil
}

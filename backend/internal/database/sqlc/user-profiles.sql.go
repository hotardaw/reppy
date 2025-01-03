// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user-profiles.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createUserProfile = `-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, first_name, last_name, date_of_birth, gender, height_inches, weight_pounds)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING profile_id, user_id, first_name, last_name, date_of_birth, gender, height_inches, weight_pounds, profile_picture_url, created_at, updated_at
`

type CreateUserProfileParams struct {
	UserID       sql.NullInt32
	FirstName    sql.NullString
	LastName     sql.NullString
	DateOfBirth  sql.NullTime
	Gender       sql.NullString
	HeightInches sql.NullInt32
	WeightPounds sql.NullInt32
}

func (q *Queries) CreateUserProfile(ctx context.Context, arg CreateUserProfileParams) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, createUserProfile,
		arg.UserID,
		arg.FirstName,
		arg.LastName,
		arg.DateOfBirth,
		arg.Gender,
		arg.HeightInches,
		arg.WeightPounds,
	)
	var i UserProfile
	err := row.Scan(
		&i.ProfileID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.Gender,
		&i.HeightInches,
		&i.WeightPounds,
		&i.ProfilePictureUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllUserProfiles = `-- name: DeleteAllUserProfiles :exec
DELETE FROM user_profiles
`

func (q *Queries) DeleteAllUserProfiles(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUserProfiles)
	return err
}

const deleteUserProfile = `-- name: DeleteUserProfile :one
DELETE FROM user_profiles
WHERE user_id = $1
RETURNING profile_id, user_id, first_name, last_name, date_of_birth, gender, height_inches, weight_pounds, profile_picture_url, created_at, updated_at
`

func (q *Queries) DeleteUserProfile(ctx context.Context, userID sql.NullInt32) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, deleteUserProfile, userID)
	var i UserProfile
	err := row.Scan(
		&i.ProfileID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.Gender,
		&i.HeightInches,
		&i.WeightPounds,
		&i.ProfilePictureUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllActiveUserProfiles = `-- name: GetAllActiveUserProfiles :many
SELECT
  ups.profile_id,
  ups.user_id,
  ups.first_name,
  ups.last_name,
  ups.date_of_birth,
  ups.gender,
  ups.height_inches,
  ups.weight_pounds,
  ups.profile_picture_url,
  ups.created_at,
  ups.updated_at,
  u.active
FROM user_profiles ups
JOIN users u ON ups.user_id = u.user_id
WHERE u.active = true
`

type GetAllActiveUserProfilesRow struct {
	ProfileID         int32
	UserID            sql.NullInt32
	FirstName         sql.NullString
	LastName          sql.NullString
	DateOfBirth       sql.NullTime
	Gender            sql.NullString
	HeightInches      sql.NullInt32
	WeightPounds      sql.NullInt32
	ProfilePictureUrl sql.NullString
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
	Active            sql.NullBool
}

// Weird SQLc error with using user_profiles.* here - generated empty select statement. Explicitly naming columns to select instead.
func (q *Queries) GetAllActiveUserProfiles(ctx context.Context) ([]GetAllActiveUserProfilesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllActiveUserProfiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllActiveUserProfilesRow
	for rows.Next() {
		var i GetAllActiveUserProfilesRow
		if err := rows.Scan(
			&i.ProfileID,
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.DateOfBirth,
			&i.Gender,
			&i.HeightInches,
			&i.WeightPounds,
			&i.ProfilePictureUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Active,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllInactiveUserProfiles = `-- name: GetAllInactiveUserProfiles :many
SELECT 
  ups.profile_id,
  ups.user_id,
  ups.first_name,
  ups.last_name,
  ups.date_of_birth,
  ups.gender,
  ups.height_inches,
  ups.weight_pounds,
  ups.profile_picture_url,
  ups.created_at,
  ups.updated_at,
  u.active
FROM user_profiles ups
JOIN users u ON ups.user_id = u.user_id
WHERE u.active = false
`

type GetAllInactiveUserProfilesRow struct {
	ProfileID         int32
	UserID            sql.NullInt32
	FirstName         sql.NullString
	LastName          sql.NullString
	DateOfBirth       sql.NullTime
	Gender            sql.NullString
	HeightInches      sql.NullInt32
	WeightPounds      sql.NullInt32
	ProfilePictureUrl sql.NullString
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
	Active            sql.NullBool
}

// Same thing here - explicitly naming columns to return.
func (q *Queries) GetAllInactiveUserProfiles(ctx context.Context) ([]GetAllInactiveUserProfilesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllInactiveUserProfiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllInactiveUserProfilesRow
	for rows.Next() {
		var i GetAllInactiveUserProfilesRow
		if err := rows.Scan(
			&i.ProfileID,
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.DateOfBirth,
			&i.Gender,
			&i.HeightInches,
			&i.WeightPounds,
			&i.ProfilePictureUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Active,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUserProfiles = `-- name: GetAllUserProfiles :many
SELECT user_profiles.profile_id, user_profiles.user_id, user_profiles.first_name, user_profiles.last_name, user_profiles.date_of_birth, user_profiles.gender, user_profiles.height_inches, user_profiles.weight_pounds, user_profiles.profile_picture_url, user_profiles.created_at, user_profiles.updated_at
FROM user_profiles
`

func (q *Queries) GetAllUserProfiles(ctx context.Context) ([]UserProfile, error) {
	rows, err := q.db.QueryContext(ctx, getAllUserProfiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserProfile
	for rows.Next() {
		var i UserProfile
		if err := rows.Scan(
			&i.ProfileID,
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.DateOfBirth,
			&i.Gender,
			&i.HeightInches,
			&i.WeightPounds,
			&i.ProfilePictureUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserProfile = `-- name: GetUserProfile :one
SELECT user_profiles.profile_id, user_profiles.user_id, user_profiles.first_name, user_profiles.last_name, user_profiles.date_of_birth, user_profiles.gender, user_profiles.height_inches, user_profiles.weight_pounds, user_profiles.profile_picture_url, user_profiles.created_at, user_profiles.updated_at
FROM user_profiles
WHERE user_profiles.user_id = $1
`

func (q *Queries) GetUserProfile(ctx context.Context, userID sql.NullInt32) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, getUserProfile, userID)
	var i UserProfile
	err := row.Scan(
		&i.ProfileID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.Gender,
		&i.HeightInches,
		&i.WeightPounds,
		&i.ProfilePictureUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserProfile = `-- name: UpdateUserProfile :one
UPDATE user_profiles
SET 
  first_name = $2, 
  last_name = $3, 
  date_of_birth = $4,
  gender = $5,
  height_inches = $6,
  weight_pounds = $7,
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING profile_id, user_id, first_name, last_name, date_of_birth, gender, height_inches, weight_pounds, profile_picture_url, created_at, updated_at
`

type UpdateUserProfileParams struct {
	UserID       sql.NullInt32
	FirstName    sql.NullString
	LastName     sql.NullString
	DateOfBirth  sql.NullTime
	Gender       sql.NullString
	HeightInches sql.NullInt32
	WeightPounds sql.NullInt32
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, updateUserProfile,
		arg.UserID,
		arg.FirstName,
		arg.LastName,
		arg.DateOfBirth,
		arg.Gender,
		arg.HeightInches,
		arg.WeightPounds,
	)
	var i UserProfile
	err := row.Scan(
		&i.ProfileID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.Gender,
		&i.HeightInches,
		&i.WeightPounds,
		&i.ProfilePictureUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

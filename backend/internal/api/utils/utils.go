package utils

import (
	"database/sql"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
	"time"
)

// TODO: upgrade these to accept a variadic input for valid checking in instances like "01-seed-users.go"

// Used in APIs' SQLc params section for compact conversions.
func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// Used in APIs' SQLc params section for compact conversions.
func ToNullInt32(num interface{}) sql.NullInt32 {
	switch n := num.(type) {
	case int:
		return sql.NullInt32{
			Int32: int32(n),
			Valid: true,
		}
	case int64:
		return sql.NullInt32{
			Int32: int32(n),
			Valid: true,
		}
	case int32:
		return sql.NullInt32{
			Int32: n,
			Valid: true,
		}
	default:
		return sql.NullInt32{Valid: false} // invalid NullInt32 is safer than panicking
	}
}

// Used in APIs' SQLc params section for compact conversions.
func ToNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

// Used in APIs with optional fields that utilize pointers
func IntPtr(i int32) *int32         { return &i }
func StrPtr(s string) *string       { return &s }
func Float32Ptr(f float32) *float32 { return &f }

// For converting pointers to SQL null types
func NullIntFromIntPtr(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func NullStringFromStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NullStringFromFloat32Ptr(f *float32) sql.NullString {
	if f == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: fmt.Sprintf("%.1f", *f), Valid: true}
}

func NullResistanceTypeFromPtr(s *string) sqlc.NullResistanceTypeEnum {
	if s == nil {
		return sqlc.NullResistanceTypeEnum{}
	}
	return sqlc.NullResistanceTypeEnum{ResistanceTypeEnum: sqlc.ResistanceTypeEnum(*s), Valid: true}
}

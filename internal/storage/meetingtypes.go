package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/one-thousand-twoth/calendar/internal/meetingtype"
	"github.com/one-thousand-twoth/calendar/internal/storage/ent"
	entmeetingtype "github.com/one-thousand-twoth/calendar/internal/storage/ent/meetingtype"
	"modernc.org/sqlite"
)

var _ meetingtype.Repository = (*Store)(nil)

func (s *Store) CreateMeetingType(ctx context.Context, input meetingtype.MeetingType, key string) (meetingtype.MeetingType, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return meetingtype.MeetingType{}, storageError(err)
	}
	defer func() { _ = tx.Rollback() }()
	// The IMMEDIATE transaction serializes this check and insertion across writers.
	exists, err := tx.MeetingType.Query().Where(entmeetingtype.NameKeyEQ(key)).Exist(ctx)
	if err != nil {
		return meetingtype.MeetingType{}, storageError(err)
	}
	if exists {
		return meetingtype.MeetingType{}, meetingtype.ErrDuplicate
	}
	row, err := tx.MeetingType.Create().SetName(input.Name).SetNameKey(key).
		SetDurationMinutes(input.DurationMinutes).SetNillableDescription(input.Description).Save(ctx)
	if err != nil {
		return meetingtype.MeetingType{}, storageError(err)
	}
	if err := tx.Commit(); err != nil {
		return meetingtype.MeetingType{}, storageError(err)
	}
	return fromRow(row), nil
}

func (s *Store) ListMeetingTypes(ctx context.Context) ([]meetingtype.MeetingType, error) {
	rows, err := s.client.MeetingType.Query().Order(ent.Asc(entmeetingtype.FieldID)).All(ctx)
	if err != nil {
		return nil, storageError(err)
	}
	result := make([]meetingtype.MeetingType, 0, len(rows))
	for _, row := range rows {
		result = append(result, fromRow(row))
	}
	return result, nil
}

func storageError(err error) error {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		// The low byte is the primary code, including for extended SQLite errors.
		code := sqliteErr.Code() & 0xff
		if code == 5 || code == 6 { // SQLITE_BUSY or SQLITE_LOCKED
			return fmt.Errorf("%w: %w", meetingtype.ErrUnavailable, err)
		}
	}
	return err
}

func fromRow(row *ent.MeetingType) meetingtype.MeetingType {
	return meetingtype.MeetingType{
		ID: row.ID, Name: row.Name, DurationMinutes: row.DurationMinutes, Description: row.Description,
	}
}

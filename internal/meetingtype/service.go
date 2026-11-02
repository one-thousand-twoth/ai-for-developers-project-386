// Package meetingtype owns the rules for the calendar owner's meeting types.
package meetingtype

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
)

var (
	ErrInvalid     = errors.New("invalid meeting type")
	ErrDuplicate   = errors.New("meeting type name already exists")
	ErrUnavailable = errors.New("storage temporarily unavailable")
)

type MeetingType struct {
	ID              string
	Name            string
	DurationMinutes int
	Description     *string
}

type Repository interface {
	CreateMeetingType(context.Context, MeetingType, string) (MeetingType, error)
	ListMeetingTypes(context.Context) ([]MeetingType, error)
}

type Service struct{ repository Repository }

func New(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, input MeetingType) (MeetingType, error) {
	if strings.TrimSpace(input.Name) == "" || utf8.RuneCountInString(input.Name) > 200 {
		return MeetingType{}, fmt.Errorf("%w: name must contain a non-whitespace character and be at most 200 Unicode characters", ErrInvalid)
	}
	if input.DurationMinutes < 1 || input.DurationMinutes > 480 {
		return MeetingType{}, fmt.Errorf("%w: durationMinutes must be between 1 and 480", ErrInvalid)
	}
	if input.Description != nil && utf8.RuneCountInString(*input.Description) > 2000 {
		return MeetingType{}, fmt.Errorf("%w: description must be at most 2000 Unicode characters", ErrInvalid)
	}
	input.Name = strings.TrimSpace(input.Name)
	key := cases.Fold().String(input.Name)
	return s.repository.CreateMeetingType(ctx, input, key)
}

func (s *Service) List(ctx context.Context) ([]MeetingType, error) {
	return s.repository.ListMeetingTypes(ctx)
}

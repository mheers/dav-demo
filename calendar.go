package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type calendarBackend struct {
	prefix string

	calendars []caldav.Calendar
	objectMap map[string][]caldav.CalendarObject
}

func (b *calendarBackend) CreateCalendar(ctx context.Context, calendar *caldav.Calendar) error {
	return nil
}

func (b *calendarBackend) Calendar(ctx context.Context) (*caldav.Calendar, error) {
	return &b.calendars[0], nil
}

func (b *calendarBackend) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	return b.calendars, nil
}

func (b *calendarBackend) GetCalendar(ctx context.Context, path string) (*caldav.Calendar, error) {
	for _, cal := range b.calendars {
		if cal.Path == path {
			return &cal, nil
		}
	}
	return nil, fmt.Errorf("Calendar for path: %s not found", path)
}

func (b *calendarBackend) CalendarHomeSetPath(ctx context.Context) (string, error) {
	username, err := currentUsername(ctx)
	return fmt.Sprintf("/%s/%s/calendars/", b.prefix, username), err
}

// must begin and end with a slash
func (b *calendarBackend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	username, err := currentUsername(ctx)
	return "/" + url.PathEscape(username) + "/", err
}

func (b *calendarBackend) DeleteCalendarObject(ctx context.Context, path string) error {
	return nil
}

func (b *calendarBackend) GetCalendarObject(ctx context.Context, path string, req *caldav.CalendarCompRequest) (*caldav.CalendarObject, error) {
	for _, objs := range b.objectMap {
		for _, obj := range objs {
			if obj.Path == path {
				return &obj, nil
			}
		}
	}
	return nil, fmt.Errorf("Couldn't find calendar object at: %s", path)
}

func (b *calendarBackend) PutCalendarObject(ctx context.Context, path string, calendar *ical.Calendar, opts *caldav.PutCalendarObjectOptions) (*caldav.CalendarObject, error) {
	return nil, nil
}

func (b *calendarBackend) ListCalendarObjects(ctx context.Context, path string, req *caldav.CalendarCompRequest) ([]caldav.CalendarObject, error) {
	return b.objectMap[path], nil
}

func (b *calendarBackend) QueryCalendarObjects(ctx context.Context, path string, query *caldav.CalendarQuery) ([]caldav.CalendarObject, error) {
	return nil, nil
}

func NewCalDavHandler(prefix, username string) http.Handler {
	sessionsCal := caldav.Calendar{
		Path:                  fmt.Sprintf("/%s/%s/calendars/sessions", prefix, username),
		SupportedComponentSet: []string{ical.CompEvent},
	}

	calendars := []caldav.Calendar{
		sessionsCal,
	}
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")

	eventSummary := fmt.Sprintf("Jetzt: %s", username)
	event := ical.NewEvent()
	event.Name = ical.CompEvent
	event.Props.SetText(ical.PropUID, "46bbf47a-1861-41a3-ae06-8d8268c6d41e")
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
	event.Props.SetText(ical.PropSummary, eventSummary)
	event.Props.SetDateTime(ical.PropDateTimeStart, time.Now())
	event.Props.SetDateTime(ical.PropDateTimeEnd, time.Now().Add(1*time.Hour))

	cal.Children = []*ical.Component{
		event.Component,
	}

	object := caldav.CalendarObject{
		Path: fmt.Sprintf("/%s/%s/calendars/sessions/test.ics", prefix, username),
		Data: cal,
	}

	return &caldav.Handler{
		Prefix: fmt.Sprintf("/%s/", prefix),
		Backend: &calendarBackend{
			calendars: calendars,
			objectMap: map[string][]caldav.CalendarObject{
				sessionsCal.Path: {object},
			},
		},
	}
}

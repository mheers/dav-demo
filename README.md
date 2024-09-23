# dav-demo

This is a simple WebDAV, CardDAV and CalDAV server written in Go.

# Usage

## WebDAV

Open nautilus and type `dav://marcel@localhost:8086/files` in the address bar.

## CardDAV

In thunderbird address book add a new entry for http://localhost:8086/carddav/marcel/contacts/private

## CalDAV

In thunderbird calendar add a new entry for http://localhost:8086/caldav/marcel/calendars/sessions/test.ics

# TODO

- [x] basic auth
- [x] serve files
- [x] serve contacts
- [x] serve calendar
- [x] carddav: add missing current-user-privilege-set to response
- [ ] fix *sync* issue with vcards under thunderbird (ctag seems to be wrong)
- [ ] write file sync on base of github.com/Elnus/webdav-cli

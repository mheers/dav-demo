package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav/carddav"
)

type contactBackend struct {
	prefix   string
	contacts []User
}

// must begin and end with a slash
func (b *contactBackend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	username, err := currentUsername(ctx)
	return "/" + url.PathEscape(username) + "/", err
}

func (b *contactBackend) GetAddressObject(ctx context.Context, path string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	username, err := currentUsername(ctx)
	if err != nil {
		return nil, err
	}
	p := fmt.Sprintf("/%s/contacts/private", username)
	if b.prefix != "" {
		p = fmt.Sprintf("/%s%s", b.prefix, p)
	}
	return &carddav.AddressObject{
		Path:          p,
		Card:          vcardFromUser(b.contacts[0]),
		ModTime:       time.Now(),
		ContentLength: 100,
		ETag:          "125",
	}, nil
}
func (b *contactBackend) ListAddressObjects(ctx context.Context, path string, req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	entry, err := b.GetAddressObject(ctx, path, req)
	if err != nil {
		return nil, err
	}
	return []carddav.AddressObject{
		*entry,
	}, nil
}
func (b *contactBackend) QueryAddressObjects(ctx context.Context, path string, query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	return nil, nil
}
func (b *contactBackend) PutAddressObject(ctx context.Context, path string, card vcard.Card, opts *carddav.PutAddressObjectOptions) (*carddav.AddressObject, error) {
	return nil, nil
}
func (b *contactBackend) DeleteAddressObject(ctx context.Context, path string) error {
	return nil
}
func (b *contactBackend) CreateAddressBook(ctx context.Context, addressBook *carddav.AddressBook) error {
	return nil
}
func (b *contactBackend) DeleteAddressBook(ctx context.Context, path string) error {
	return nil
}
func (b *contactBackend) GetAddressBook(ctx context.Context, path string) (*carddav.AddressBook, error) {
	username, err := currentUsername(ctx)
	if err != nil {
		return nil, err
	}
	p := fmt.Sprintf("/%s/contacts/private", username)
	if b.prefix != "" {
		p = fmt.Sprintf("/%s%s", b.prefix, p)
	}
	return &carddav.AddressBook{
		Path:            p,
		Description:     "My address book",
		Name:            "My address book",
		MaxResourceSize: 1024,
		SupportedAddressData: []carddav.AddressDataType{
			{ContentType: vcard.MIMEType, Version: "3.0"},
			{ContentType: vcard.MIMEType, Version: "4.0"},
		},
		CTag: "001",
	}, nil
}
func (b *contactBackend) ListAddressBooks(ctx context.Context) ([]carddav.AddressBook, error) {
	username, err := currentUsername(ctx)
	if err != nil {
		return nil, err
	}
	p := fmt.Sprintf("/%s/contacts/private", username)
	if b.prefix != "" {
		p = fmt.Sprintf("/%s%s", b.prefix, p)
	}
	addressBook, err := b.GetAddressBook(ctx, p)
	if err != nil {
		return nil, err
	}
	return []carddav.AddressBook{
		*addressBook,
	}, nil
}

func (b *contactBackend) AddressBookHomeSetPath(ctx context.Context) (string, error) {
	principal, err := b.CurrentUserPrincipal(ctx)
	return principal + "contacts/", err
}

func NewCardDAVHandler(prefix string, contacts []User) http.Handler {
	p := ""
	if prefix != "" {
		p = fmt.Sprintf("/%s/", prefix)
	}
	return &carddav.Handler{
		Prefix: p,
		Backend: &contactBackend{
			prefix:   prefix,
			contacts: contacts,
		},
	}
}

func currentUsername(ctx context.Context) (string, error) {
	if v, ok := ctx.Value(CtxKey{}).(CtxValue); ok {
		return v.Username, nil
	}
	return "", errors.New("not authenticated")
}

func utf8Field(v string) *vcard.Field {
	return &vcard.Field{
		Value: v,
		Params: vcard.Params{
			"CHARSET": []string{"UTF-8"},
		},
	}
}

type User struct {
	Firstname string
	Lastname  string
	Unit      string
	UpdatedAt time.Time
	Extid     string
	Emails    map[string]string
	Phones    map[string]string
}

func vcardFromUser(u User) vcard.Card {
	c := vcard.Card{}

	c.Set(vcard.FieldFormattedName, utf8Field(u.Firstname+" "+u.Lastname))
	c.SetName(&vcard.Name{
		Field:      utf8Field(""),
		FamilyName: u.Lastname,
		GivenName:  u.Firstname,
	})
	c.SetRevision(u.UpdatedAt)
	c.SetValue(vcard.FieldUID, u.Extid)

	c.Set(vcard.FieldOrganization, utf8Field(u.Unit))

	// addFields sorts the key to ensure a stable order
	addFields := func(fieldName string, values map[string]string) {
		for _, k := range slices.Sorted(maps.Keys(values)) {
			v := values[k]
			c.Add(fieldName, &vcard.Field{
				Value: v,
				Params: vcard.Params{
					vcard.ParamType: []string{k + ";CHARSET=UTF-8"}, // hacky but prevent maps ordering issues
					// "CHARSET":       []string{"UTF-8"},
				},
			})
		}
	}

	addFields(vcard.FieldEmail, u.Emails)
	addFields(vcard.FieldTelephone, u.Phones)

	vcard.ToV4(c)
	return c
}

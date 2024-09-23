module github.com/mheers/dav-demo

go 1.23.1

require (
	github.com/emersion/go-ical v0.0.0-20240127095438-fc1c9d8fb2b6
	github.com/emersion/go-vcard v0.0.0-20230815062825-8fda7d206ec9
	github.com/emersion/go-webdav v0.5.1-0.20240419143909-21f251fa1de2
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.29.0
)

require (
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/teambition/rrule-go v1.8.2 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

replace github.com/emersion/go-webdav => github.com/mheers/go-webdav v0.0.0-20240923131844-0fc5736b465f

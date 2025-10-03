module github.com/Nemutagk/govalidator/v2

go 1.25.0

require (
	github.com/Nemutagk/goerrors v1.2.2
	github.com/gofrs/uuid v4.4.0+incompatible
)

require github.com/Nemutagk/godb/v2 v2.0.0 // indirect

// replace github.com/Nemutagk/godb/v2 => /opt/modules/db

// replace github.com/Nemutagk/goerrors => /opt/modules/goerrors

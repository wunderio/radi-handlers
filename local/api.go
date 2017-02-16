package local

import (
	"context"
	"os/user"

	handler_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

// Settings needed to make a local API
type LocalAPISettings struct {
	handler_bytesource.BytesourceFileSettings
	Context context.Context
	User    user.User
}

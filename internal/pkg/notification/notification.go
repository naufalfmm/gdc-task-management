package notification

import (
	"context"
	"encoding/json"
)

type NotificationData interface {
	json.Marshaler
}

type Notification interface {
	Send(ctx context.Context, message string, data NotificationData) error
}

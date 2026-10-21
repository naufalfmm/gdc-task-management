package dummy

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/naufalfmm/gdc-task-management/internal/pkg/logger"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/notification"
)

type dummyNotification struct {
	logger logger.Logger
}

func NewDummyNotification(logger logger.Logger) notification.Notification {
	return &dummyNotification{
		logger: logger,
	}
}

type dummyLoggerData struct {
	data notification.NotificationData
}

func (d dummyLoggerData) Args() []any {
	if d.data == nil {
		return nil
	}

	v := reflect.ValueOf(d.data)

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		result := make([]any, 0, v.NumField()*2)

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			if !field.IsExported() {
				continue // Abaikan private field
			}

			key := getFieldName(field)
			if key == "-" {
				continue
			}

			result = append(result, key, fieldValue.Interface())
		}

		return result
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}

		// Urutkan key map agar hasilnya konsisten (tidak acak saat di-run)
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
		})

		result := make([]any, 0, len(keys)*2)
		for _, k := range keys {
			val := v.MapIndex(k)
			result = append(result, k.Interface(), val.Interface())
		}

		return result
	default:
		typeName := v.Type().Name()
		if typeName == "" {
			typeName = v.Kind().String()
		}

		return []any{typeName, v.Interface()}
	}
}

func getFieldName(field reflect.StructField) string {
	if tag := field.Tag.Get("json"); tag != "" {
		name := strings.Split(tag, ",")[0]
		if name != "" {
			return name
		}
	}
	if tag := field.Tag.Get("db"); tag != "" {
		name := strings.Split(tag, ",")[0]
		if name != "" {
			return name
		}
	}
	return field.Name
}

func (d dummyNotification) Send(ctx context.Context, message string, data notification.NotificationData) error {
	d.logger.Info(ctx, message, dummyLoggerData{data: data})
	return nil
}

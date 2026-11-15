// Package jsoncompat preserves the Kratos v2 application/json wire format
// while the HTTP transport migrates to Kratos v3.
package jsoncompat

import (
	"encoding/json"
	"reflect"

	"github.com/go-kratos/kratos/v3/encoding"
	_ "github.com/go-kratos/kratos/v3/transport" // register core codecs first
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}

func (codec) Name() string { return "json" }

func (codec) Marshal(value any) ([]byte, error) {
	switch typed := value.(type) {
	case json.Marshaler:
		return typed.MarshalJSON()
	case proto.Message:
		return (protojson.MarshalOptions{EmitUnpopulated: true}).Marshal(typed)
	default:
		return json.Marshal(value)
	}
}

func (codec) Unmarshal(data []byte, value any) error {
	switch typed := value.(type) {
	case json.Unmarshaler:
		return typed.UnmarshalJSON(data)
	case proto.Message:
		return (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, typed)
	default:
		reflected := reflect.ValueOf(value)
		for reflected.Kind() == reflect.Pointer {
			if reflected.IsNil() {
				reflected.Set(reflect.New(reflected.Type().Elem()))
			}
			reflected = reflected.Elem()
		}
		if message, ok := reflect.Indirect(reflected).Interface().(proto.Message); ok {
			return (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, message)
		}
		return json.Unmarshal(data, value)
	}
}

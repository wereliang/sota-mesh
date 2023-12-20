package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/mitchellh/mapstructure"
)

type TypeConfig interface {
	Type() string
	Unmarshal(obj interface{}) error
}

type TypeConfigImpl map[string]interface{}

func (tc TypeConfigImpl) Type() string {
	return tc["@type"].(string)
}

func (tc TypeConfigImpl) Unmarshal(obj interface{}) error {

	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   obj,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	return decoder.Decode(tc)
}

type AnyTypeConfig struct {
	A *any.Any
}

func (tc *AnyTypeConfig) Type() string {
	return tc.A.GetTypeUrl()
}

func (tc *AnyTypeConfig) Unmarshal(obj interface{}) error {
	return ptypes.UnmarshalAny(tc.A, obj.(proto.Message))
}

//////////////////////////////

type ConfigProto interface {
	Type() string
}

func UnmarshalConfigProto(a TypeConfig, cfg interface{}) error {
	switch v := cfg.(type) {
	case proto.Message:
		// 特殊处理
		if _, ok := a.(TypeConfigImpl); ok {
			data, err := json.Marshal(a)
			if err != nil {
				return fmt.Errorf("json marshal error. %s", err)
			}
			us := jsonpb.Unmarshaler{AllowUnknownFields: true}
			return us.Unmarshal(strings.NewReader(string(data)), v)
		} else {
			// pb any对pb直接unmarshal
			return a.Unmarshal(cfg)
		}
	case ConfigProto:
		return a.Unmarshal(v)
	default:
		return fmt.Errorf("invalid config proto type")
	}
}

func GetConfigProtoType(cfg interface{}) (string, error) {
	switch v := cfg.(type) {
	case proto.Message:
		any, err := ptypes.MarshalAny(v)
		if err != nil {
			return "", err
		}
		return any.GetTypeUrl(), nil
	case ConfigProto:
		return v.Type(), nil
	default:
		return "", fmt.Errorf("invalid config type")
	}
}

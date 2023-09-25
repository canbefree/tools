package marshaler

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
)

type Marshaler struct {
}

func NewMarshaler() *Marshaler {
	return &Marshaler{}
}

type protoBody struct {
	Body []byte
	Err  error
}

func (m *Marshaler) MarshalWrapper(respBody ...interface{}) ([]byte, error) {
	if len(respBody) != 2 {
		return nil, fmt.Errorf("MarshalWrapper err len %d", len(respBody))
	}
	var respbody proto.Message
	var resperr error
	if respBody[0] != nil {
		respbody = respBody[0].(proto.Message)
	}
	if respBody[1] != nil {
		resperr = respBody[1].(error)
	}

	v, err := proto.Marshal(respbody)
	if err != nil {
		return nil, err
	}

	protoBody := &protoBody{
		Body: v,
		Err:  resperr,
	}

	return json.Marshal(protoBody)
}

func (m *Marshaler) UnMarshalWrapper(strings []byte, resp any) ([]interface{}, error) {
	protoBody := &protoBody{}
	err := json.Unmarshal(strings, protoBody)
	if err != nil {
		return nil, err
	}

	if protoBody.Body == nil {
		return []interface{}{protoBody.Body, protoBody.Err}, nil
	}

	v := resp.([]interface{})
	bodyF, ok := v[0].(proto.Message)
	if !ok {
		return nil, fmt.Errorf("UnMarshalWrapper err type %T", v)
	}

	var errF error
	if v[1] != nil {
		errF, ok = v[1].(error)
		if !ok {
			return nil, fmt.Errorf("UnMarshalWrapper err type %T", v)
		}
	}

	if err := proto.Unmarshal(protoBody.Body, bodyF); err != nil {
		return nil, err
	}

	return []interface{}{bodyF, errF}, nil
}

package utils

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"

	"github.com/flyteorg/flytestdlib/utils/prototest"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

type SimpleType struct {
	StringValue string `json:"string_value,omitempty"`
}

func TestMarshalPbToString(t *testing.T) {
	type args struct {
		msg proto.Message
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{msg: &prototest.TestProto{}}, "{}", false},
		{"has value", args{msg: &prototest.TestProto{StringValue: "hello"}}, `{"stringValue":"hello"}`, false},
		{"nil input", args{msg: nil}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalPbToString(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MarshalToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalObjToStruct(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *structpb.Struct
		wantErr bool
	}{
		{"has proto value", args{input: &prototest.TestProto{StringValue: "hello"}}, &structpb.Struct{Fields: map[string]*structpb.Value{
			"stringValue": {Kind: &structpb.Value_StringValue{StringValue: "hello"}},
		}}, false},
		{"has struct value", args{input: SimpleType{StringValue: "hello"}}, &structpb.Struct{Fields: map[string]*structpb.Value{
			"string_value": {Kind: &structpb.Value_StringValue{StringValue: "hello"}},
		}}, false},
		{"has string value", args{input: "hello"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalObjToStruct(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalObjToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("MarshalObjToStruct() = %v, want %v, diff: %v", got, tt.want, diff)
			}
		})
	}
}

func TestUnmarshalStructToPb(t *testing.T) {
	type args struct {
		structObj *structpb.Struct
		msg       proto.Message
	}
	tests := []struct {
		name     string
		args     args
		expected proto.Message
		wantErr  bool
	}{
		{"empty", args{structObj: &structpb.Struct{Fields: map[string]*structpb.Value{}}, msg: &prototest.TestProto{}}, &prototest.TestProto{}, false},
		{"has value", args{structObj: &structpb.Struct{Fields: map[string]*structpb.Value{
			"stringValue": {Kind: &structpb.Value_StringValue{StringValue: "hello"}},
		}}, msg: &prototest.TestProto{}}, &prototest.TestProto{StringValue: "hello"}, false},
		{"nil input", args{structObj: nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnmarshalStructToPb(tt.args.structObj, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalStructToPb() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.expected == nil {
				assert.Nil(t, tt.args.msg)
			} else {
				assert.Equal(t, (tt.expected.(*prototest.TestProto)).StringValue, (tt.args.msg.(*prototest.TestProto)).StringValue)
			}
		})
	}
}

func TestMarshalPbToStruct(t *testing.T) {
	type args struct {
		in proto.Message
	}
	tests := []struct {
		name     string
		args     args
		expected *structpb.Struct
		wantErr  bool
	}{
		{"empty", args{in: &prototest.TestProto{}}, &structpb.Struct{Fields: map[string]*structpb.Value{}}, false},
		{"has value",
			args{
				in: &prototest.TestProto{StringValue: "hello"},
			},
			&structpb.Struct{
				Fields: map[string]*structpb.Value{
					"stringValue": {Kind: &structpb.Value_StringValue{StringValue: "hello"}},
				},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MarshalPbToStruct(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("MarshalPbToStruct() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(tt.expected.Fields) == 0 {
				assert.Empty(t, got.Fields)
			} else {
				assert.Equal(t, tt.expected.Fields["stringValue"].Kind, got.Fields["stringValue"].Kind)
			}
		})
	}
}

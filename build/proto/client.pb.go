// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.4
// source: proto/client.proto

package client

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageContainer_Status int32

const (
	MessageContainer_UNPROCESSED MessageContainer_Status = 0
	MessageContainer_PROCESSING  MessageContainer_Status = 1
	MessageContainer_PROCESSED   MessageContainer_Status = 2
)

// Enum value maps for MessageContainer_Status.
var (
	MessageContainer_Status_name = map[int32]string{
		0: "UNPROCESSED",
		1: "PROCESSING",
		2: "PROCESSED",
	}
	MessageContainer_Status_value = map[string]int32{
		"UNPROCESSED": 0,
		"PROCESSING":  1,
		"PROCESSED":   2,
	}
)

func (x MessageContainer_Status) Enum() *MessageContainer_Status {
	p := new(MessageContainer_Status)
	*p = x
	return p
}

func (x MessageContainer_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageContainer_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_client_proto_enumTypes[0].Descriptor()
}

func (MessageContainer_Status) Type() protoreflect.EnumType {
	return &file_proto_client_proto_enumTypes[0]
}

func (x MessageContainer_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageContainer_Status.Descriptor instead.
func (MessageContainer_Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_client_proto_rawDescGZIP(), []int{1, 0}
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_client_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_proto_client_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_proto_client_proto_rawDescGZIP(), []int{0}
}

func (x *Status) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type MessageContainer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dummy string `protobuf:"bytes,2,opt,name=Dummy,proto3" json:"Dummy,omitempty"`
}

func (x *MessageContainer) Reset() {
	*x = MessageContainer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_client_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageContainer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageContainer) ProtoMessage() {}

func (x *MessageContainer) ProtoReflect() protoreflect.Message {
	mi := &file_proto_client_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageContainer.ProtoReflect.Descriptor instead.
func (*MessageContainer) Descriptor() ([]byte, []int) {
	return file_proto_client_proto_rawDescGZIP(), []int{1}
}

func (x *MessageContainer) GetDummy() string {
	if x != nil {
		return x.Dummy
	}
	return ""
}

var File_proto_client_proto protoreflect.FileDescriptor

var file_proto_client_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x62, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x44, 0x75, 0x6d, 0x6d, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x44, 0x75, 0x6d,
	0x6d, 0x79, 0x22, 0x38, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0f, 0x0a, 0x0b,
	0x55, 0x4e, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a,
	0x0a, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0d, 0x0a,
	0x09, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x45, 0x44, 0x10, 0x02, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_client_proto_rawDescOnce sync.Once
	file_proto_client_proto_rawDescData = file_proto_client_proto_rawDesc
)

func file_proto_client_proto_rawDescGZIP() []byte {
	file_proto_client_proto_rawDescOnce.Do(func() {
		file_proto_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_client_proto_rawDescData)
	})
	return file_proto_client_proto_rawDescData
}

var file_proto_client_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_client_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_client_proto_goTypes = []interface{}{
	(MessageContainer_Status)(0), // 0: MessageContainer.Status
	(*Status)(nil),               // 1: Status
	(*MessageContainer)(nil),     // 2: MessageContainer
}
var file_proto_client_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_client_proto_init() }
func file_proto_client_proto_init() {
	if File_proto_client_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_client_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_client_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageContainer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_client_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_client_proto_goTypes,
		DependencyIndexes: file_proto_client_proto_depIdxs,
		EnumInfos:         file_proto_client_proto_enumTypes,
		MessageInfos:      file_proto_client_proto_msgTypes,
	}.Build()
	File_proto_client_proto = out.File
	file_proto_client_proto_rawDesc = nil
	file_proto_client_proto_goTypes = nil
	file_proto_client_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: events/v1/ufo.proto

package events_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Наблюдение НЛО зарегистрировано
type UFORecorded struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uuid          string                 `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`                               // Уникальный идентификатор наблюдения
	ObservedAt    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=observed_at,json=observedAt,proto3" json:"observed_at,omitempty"` // Дата и время наблюдения
	Location      string                 `protobuf:"bytes,3,opt,name=location,proto3" json:"location,omitempty"`                       // Место наблюдения
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`                 // Описание наблюдаемого объекта
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UFORecorded) Reset() {
	*x = UFORecorded{}
	mi := &file_events_v1_ufo_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UFORecorded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UFORecorded) ProtoMessage() {}

func (x *UFORecorded) ProtoReflect() protoreflect.Message {
	mi := &file_events_v1_ufo_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UFORecorded.ProtoReflect.Descriptor instead.
func (*UFORecorded) Descriptor() ([]byte, []int) {
	return file_events_v1_ufo_proto_rawDescGZIP(), []int{0}
}

func (x *UFORecorded) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *UFORecorded) GetObservedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ObservedAt
	}
	return nil
}

func (x *UFORecorded) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *UFORecorded) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var File_events_v1_ufo_proto protoreflect.FileDescriptor

const file_events_v1_ufo_proto_rawDesc = "" +
	"\n" +
	"\x13events/v1/ufo.proto\x12\tevents.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\x9c\x01\n" +
	"\vUFORecorded\x12\x12\n" +
	"\x04uuid\x18\x01 \x01(\tR\x04uuid\x12;\n" +
	"\vobserved_at\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"observedAt\x12\x1a\n" +
	"\blocation\x18\x03 \x01(\tR\blocation\x12 \n" +
	"\vdescription\x18\x04 \x01(\tR\vdescriptionBkZigithub.com/olezhek28/microservices-course-examples/week_5/clean_arch/shared/pkg/proto/events/v1;events_v1b\x06proto3"

var (
	file_events_v1_ufo_proto_rawDescOnce sync.Once
	file_events_v1_ufo_proto_rawDescData []byte
)

func file_events_v1_ufo_proto_rawDescGZIP() []byte {
	file_events_v1_ufo_proto_rawDescOnce.Do(func() {
		file_events_v1_ufo_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_events_v1_ufo_proto_rawDesc), len(file_events_v1_ufo_proto_rawDesc)))
	})
	return file_events_v1_ufo_proto_rawDescData
}

var file_events_v1_ufo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_events_v1_ufo_proto_goTypes = []any{
	(*UFORecorded)(nil),           // 0: events.v1.UFORecorded
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_events_v1_ufo_proto_depIdxs = []int32{
	1, // 0: events.v1.UFORecorded.observed_at:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_events_v1_ufo_proto_init() }
func file_events_v1_ufo_proto_init() {
	if File_events_v1_ufo_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_events_v1_ufo_proto_rawDesc), len(file_events_v1_ufo_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_v1_ufo_proto_goTypes,
		DependencyIndexes: file_events_v1_ufo_proto_depIdxs,
		MessageInfos:      file_events_v1_ufo_proto_msgTypes,
	}.Build()
	File_events_v1_ufo_proto = out.File
	file_events_v1_ufo_proto_goTypes = nil
	file_events_v1_ufo_proto_depIdxs = nil
}

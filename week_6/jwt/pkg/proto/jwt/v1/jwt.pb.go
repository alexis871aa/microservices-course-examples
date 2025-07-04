// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: jwt/v1/jwt.proto

package jwt_v1

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Запрос на логин
type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_jwt_v1_jwt_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_v1_jwt_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_jwt_v1_jwt_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Ответ на логин
type LoginResponse struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	AccessToken           string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken          string                 `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	AccessTokenExpiresAt  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=access_token_expires_at,json=accessTokenExpiresAt,proto3" json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=refresh_token_expires_at,json=refreshTokenExpiresAt,proto3" json:"refresh_token_expires_at,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_jwt_v1_jwt_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_v1_jwt_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_jwt_v1_jwt_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *LoginResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *LoginResponse) GetAccessTokenExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AccessTokenExpiresAt
	}
	return nil
}

func (x *LoginResponse) GetRefreshTokenExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.RefreshTokenExpiresAt
	}
	return nil
}

// Запрос на получение access токена
type GetAccessTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefreshToken  string                 `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAccessTokenRequest) Reset() {
	*x = GetAccessTokenRequest{}
	mi := &file_jwt_v1_jwt_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessTokenRequest) ProtoMessage() {}

func (x *GetAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_v1_jwt_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*GetAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_jwt_v1_jwt_proto_rawDescGZIP(), []int{2}
}

func (x *GetAccessTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// Ответ с новым access токеном
type GetAccessTokenResponse struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	AccessToken          string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	AccessTokenExpiresAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=access_token_expires_at,json=accessTokenExpiresAt,proto3" json:"access_token_expires_at,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *GetAccessTokenResponse) Reset() {
	*x = GetAccessTokenResponse{}
	mi := &file_jwt_v1_jwt_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAccessTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessTokenResponse) ProtoMessage() {}

func (x *GetAccessTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_v1_jwt_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessTokenResponse.ProtoReflect.Descriptor instead.
func (*GetAccessTokenResponse) Descriptor() ([]byte, []int) {
	return file_jwt_v1_jwt_proto_rawDescGZIP(), []int{3}
}

func (x *GetAccessTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *GetAccessTokenResponse) GetAccessTokenExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AccessTokenExpiresAt
	}
	return nil
}

// Запрос на получение refresh токена
type GetRefreshTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefreshToken  string                 `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRefreshTokenRequest) Reset() {
	*x = GetRefreshTokenRequest{}
	mi := &file_jwt_v1_jwt_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRefreshTokenRequest) ProtoMessage() {}

func (x *GetRefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_v1_jwt_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*GetRefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_jwt_v1_jwt_proto_rawDescGZIP(), []int{4}
}

func (x *GetRefreshTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// Ответ с новым refresh токеном
type GetRefreshTokenResponse struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	RefreshToken          string                 `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=refresh_token_expires_at,json=refreshTokenExpiresAt,proto3" json:"refresh_token_expires_at,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *GetRefreshTokenResponse) Reset() {
	*x = GetRefreshTokenResponse{}
	mi := &file_jwt_v1_jwt_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRefreshTokenResponse) ProtoMessage() {}

func (x *GetRefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_v1_jwt_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*GetRefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_jwt_v1_jwt_proto_rawDescGZIP(), []int{5}
}

func (x *GetRefreshTokenResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *GetRefreshTokenResponse) GetRefreshTokenExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.RefreshTokenExpiresAt
	}
	return nil
}

var File_jwt_v1_jwt_proto protoreflect.FileDescriptor

const file_jwt_v1_jwt_proto_rawDesc = "" +
	"\n" +
	"\x10jwt/v1/jwt.proto\x12\x06jwt.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"F\n" +
	"\fLoginRequest\x12\x1a\n" +
	"\busername\x18\x01 \x01(\tR\busername\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"\xff\x01\n" +
	"\rLoginResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\x12#\n" +
	"\rrefresh_token\x18\x02 \x01(\tR\frefreshToken\x12Q\n" +
	"\x17access_token_expires_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x14accessTokenExpiresAt\x12S\n" +
	"\x18refresh_token_expires_at\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\x15refreshTokenExpiresAt\"<\n" +
	"\x15GetAccessTokenRequest\x12#\n" +
	"\rrefresh_token\x18\x01 \x01(\tR\frefreshToken\"\x8e\x01\n" +
	"\x16GetAccessTokenResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\x12Q\n" +
	"\x17access_token_expires_at\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x14accessTokenExpiresAt\"=\n" +
	"\x16GetRefreshTokenRequest\x12#\n" +
	"\rrefresh_token\x18\x01 \x01(\tR\frefreshToken\"\x93\x01\n" +
	"\x17GetRefreshTokenResponse\x12#\n" +
	"\rrefresh_token\x18\x01 \x01(\tR\frefreshToken\x12S\n" +
	"\x18refresh_token_expires_at\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\x15refreshTokenExpiresAt2\xe7\x01\n" +
	"\n" +
	"JWTService\x124\n" +
	"\x05Login\x12\x14.jwt.v1.LoginRequest\x1a\x15.jwt.v1.LoginResponse\x12O\n" +
	"\x0eGetAccessToken\x12\x1d.jwt.v1.GetAccessTokenRequest\x1a\x1e.jwt.v1.GetAccessTokenResponse\x12R\n" +
	"\x0fGetRefreshToken\x12\x1e.jwt.v1.GetRefreshTokenRequest\x1a\x1f.jwt.v1.GetRefreshTokenResponseBQZOgithub.com/olezhek28/microservices-course-examples/week_6/jwt/api/jwt/v1;jwt_v1b\x06proto3"

var (
	file_jwt_v1_jwt_proto_rawDescOnce sync.Once
	file_jwt_v1_jwt_proto_rawDescData []byte
)

func file_jwt_v1_jwt_proto_rawDescGZIP() []byte {
	file_jwt_v1_jwt_proto_rawDescOnce.Do(func() {
		file_jwt_v1_jwt_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_jwt_v1_jwt_proto_rawDesc), len(file_jwt_v1_jwt_proto_rawDesc)))
	})
	return file_jwt_v1_jwt_proto_rawDescData
}

var (
	file_jwt_v1_jwt_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
	file_jwt_v1_jwt_proto_goTypes  = []any{
		(*LoginRequest)(nil),            // 0: jwt.v1.LoginRequest
		(*LoginResponse)(nil),           // 1: jwt.v1.LoginResponse
		(*GetAccessTokenRequest)(nil),   // 2: jwt.v1.GetAccessTokenRequest
		(*GetAccessTokenResponse)(nil),  // 3: jwt.v1.GetAccessTokenResponse
		(*GetRefreshTokenRequest)(nil),  // 4: jwt.v1.GetRefreshTokenRequest
		(*GetRefreshTokenResponse)(nil), // 5: jwt.v1.GetRefreshTokenResponse
		(*timestamppb.Timestamp)(nil),   // 6: google.protobuf.Timestamp
	}
)
var file_jwt_v1_jwt_proto_depIdxs = []int32{
	6, // 0: jwt.v1.LoginResponse.access_token_expires_at:type_name -> google.protobuf.Timestamp
	6, // 1: jwt.v1.LoginResponse.refresh_token_expires_at:type_name -> google.protobuf.Timestamp
	6, // 2: jwt.v1.GetAccessTokenResponse.access_token_expires_at:type_name -> google.protobuf.Timestamp
	6, // 3: jwt.v1.GetRefreshTokenResponse.refresh_token_expires_at:type_name -> google.protobuf.Timestamp
	0, // 4: jwt.v1.JWTService.Login:input_type -> jwt.v1.LoginRequest
	2, // 5: jwt.v1.JWTService.GetAccessToken:input_type -> jwt.v1.GetAccessTokenRequest
	4, // 6: jwt.v1.JWTService.GetRefreshToken:input_type -> jwt.v1.GetRefreshTokenRequest
	1, // 7: jwt.v1.JWTService.Login:output_type -> jwt.v1.LoginResponse
	3, // 8: jwt.v1.JWTService.GetAccessToken:output_type -> jwt.v1.GetAccessTokenResponse
	5, // 9: jwt.v1.JWTService.GetRefreshToken:output_type -> jwt.v1.GetRefreshTokenResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_jwt_v1_jwt_proto_init() }
func file_jwt_v1_jwt_proto_init() {
	if File_jwt_v1_jwt_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_jwt_v1_jwt_proto_rawDesc), len(file_jwt_v1_jwt_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jwt_v1_jwt_proto_goTypes,
		DependencyIndexes: file_jwt_v1_jwt_proto_depIdxs,
		MessageInfos:      file_jwt_v1_jwt_proto_msgTypes,
	}.Build()
	File_jwt_v1_jwt_proto = out.File
	file_jwt_v1_jwt_proto_goTypes = nil
	file_jwt_v1_jwt_proto_depIdxs = nil
}

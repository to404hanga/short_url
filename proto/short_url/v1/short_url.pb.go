// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: short_url.proto

package short_url_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type GenerateShortUrlRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginUrl     string                 `protobuf:"bytes,1,opt,name=origin_url,json=originUrl,proto3" json:"origin_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenerateShortUrlRequest) Reset() {
	*x = GenerateShortUrlRequest{}
	mi := &file_short_url_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateShortUrlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateShortUrlRequest) ProtoMessage() {}

func (x *GenerateShortUrlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_short_url_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateShortUrlRequest.ProtoReflect.Descriptor instead.
func (*GenerateShortUrlRequest) Descriptor() ([]byte, []int) {
	return file_short_url_proto_rawDescGZIP(), []int{0}
}

func (x *GenerateShortUrlRequest) GetOriginUrl() string {
	if x != nil {
		return x.OriginUrl
	}
	return ""
}

type GenerateShortUrlResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenerateShortUrlResponse) Reset() {
	*x = GenerateShortUrlResponse{}
	mi := &file_short_url_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateShortUrlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateShortUrlResponse) ProtoMessage() {}

func (x *GenerateShortUrlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_short_url_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateShortUrlResponse.ProtoReflect.Descriptor instead.
func (*GenerateShortUrlResponse) Descriptor() ([]byte, []int) {
	return file_short_url_proto_rawDescGZIP(), []int{1}
}

func (x *GenerateShortUrlResponse) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type GetOriginUrlRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOriginUrlRequest) Reset() {
	*x = GetOriginUrlRequest{}
	mi := &file_short_url_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOriginUrlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOriginUrlRequest) ProtoMessage() {}

func (x *GetOriginUrlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_short_url_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOriginUrlRequest.ProtoReflect.Descriptor instead.
func (*GetOriginUrlRequest) Descriptor() ([]byte, []int) {
	return file_short_url_proto_rawDescGZIP(), []int{2}
}

func (x *GetOriginUrlRequest) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type GetOriginUrlResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginUrl     string                 `protobuf:"bytes,1,opt,name=origin_url,json=originUrl,proto3" json:"origin_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOriginUrlResponse) Reset() {
	*x = GetOriginUrlResponse{}
	mi := &file_short_url_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOriginUrlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOriginUrlResponse) ProtoMessage() {}

func (x *GetOriginUrlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_short_url_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOriginUrlResponse.ProtoReflect.Descriptor instead.
func (*GetOriginUrlResponse) Descriptor() ([]byte, []int) {
	return file_short_url_proto_rawDescGZIP(), []int{3}
}

func (x *GetOriginUrlResponse) GetOriginUrl() string {
	if x != nil {
		return x.OriginUrl
	}
	return ""
}

var File_short_url_proto protoreflect.FileDescriptor

const file_short_url_proto_rawDesc = "" +
	"\n" +
	"\x0fshort_url.proto\x12\fshort_url.v1\"8\n" +
	"\x17GenerateShortUrlRequest\x12\x1d\n" +
	"\n" +
	"origin_url\x18\x01 \x01(\tR\toriginUrl\"7\n" +
	"\x18GenerateShortUrlResponse\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\"2\n" +
	"\x13GetOriginUrlRequest\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\"5\n" +
	"\x14GetOriginUrlResponse\x12\x1d\n" +
	"\n" +
	"origin_url\x18\x01 \x01(\tR\toriginUrl2\xcb\x01\n" +
	"\x0fShortUrlService\x12a\n" +
	"\x10GenerateShortUrl\x12%.short_url.v1.GenerateShortUrlRequest\x1a&.short_url.v1.GenerateShortUrlResponse\x12U\n" +
	"\fGetOriginUrl\x12!.short_url.v1.GetOriginUrlRequest\x1a\".short_url.v1.GetOriginUrlResponseB\x1bZ\x19short_url/v1;short_url_v1b\x06proto3"

var (
	file_short_url_proto_rawDescOnce sync.Once
	file_short_url_proto_rawDescData []byte
)

func file_short_url_proto_rawDescGZIP() []byte {
	file_short_url_proto_rawDescOnce.Do(func() {
		file_short_url_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_short_url_proto_rawDesc), len(file_short_url_proto_rawDesc)))
	})
	return file_short_url_proto_rawDescData
}

var file_short_url_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_short_url_proto_goTypes = []any{
	(*GenerateShortUrlRequest)(nil),  // 0: short_url.v1.GenerateShortUrlRequest
	(*GenerateShortUrlResponse)(nil), // 1: short_url.v1.GenerateShortUrlResponse
	(*GetOriginUrlRequest)(nil),      // 2: short_url.v1.GetOriginUrlRequest
	(*GetOriginUrlResponse)(nil),     // 3: short_url.v1.GetOriginUrlResponse
}
var file_short_url_proto_depIdxs = []int32{
	0, // 0: short_url.v1.ShortUrlService.GenerateShortUrl:input_type -> short_url.v1.GenerateShortUrlRequest
	2, // 1: short_url.v1.ShortUrlService.GetOriginUrl:input_type -> short_url.v1.GetOriginUrlRequest
	1, // 2: short_url.v1.ShortUrlService.GenerateShortUrl:output_type -> short_url.v1.GenerateShortUrlResponse
	3, // 3: short_url.v1.ShortUrlService.GetOriginUrl:output_type -> short_url.v1.GetOriginUrlResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_short_url_proto_init() }
func file_short_url_proto_init() {
	if File_short_url_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_short_url_proto_rawDesc), len(file_short_url_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_short_url_proto_goTypes,
		DependencyIndexes: file_short_url_proto_depIdxs,
		MessageInfos:      file_short_url_proto_msgTypes,
	}.Build()
	File_short_url_proto = out.File
	file_short_url_proto_goTypes = nil
	file_short_url_proto_depIdxs = nil
}

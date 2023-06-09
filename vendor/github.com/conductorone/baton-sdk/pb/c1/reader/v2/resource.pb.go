// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0-devel
// 	protoc        (unknown)
// source: c1/reader/v2/resource.proto

package v2

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResourceTypesReaderServiceGetResourceTypeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceTypeId string `protobuf:"bytes,1,opt,name=resource_type_id,json=resourceTypeId,proto3" json:"resource_type_id,omitempty"`
}

func (x *ResourceTypesReaderServiceGetResourceTypeRequest) Reset() {
	*x = ResourceTypesReaderServiceGetResourceTypeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_c1_reader_v2_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceTypesReaderServiceGetResourceTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceTypesReaderServiceGetResourceTypeRequest) ProtoMessage() {}

func (x *ResourceTypesReaderServiceGetResourceTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_c1_reader_v2_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceTypesReaderServiceGetResourceTypeRequest.ProtoReflect.Descriptor instead.
func (*ResourceTypesReaderServiceGetResourceTypeRequest) Descriptor() ([]byte, []int) {
	return file_c1_reader_v2_resource_proto_rawDescGZIP(), []int{0}
}

func (x *ResourceTypesReaderServiceGetResourceTypeRequest) GetResourceTypeId() string {
	if x != nil {
		return x.ResourceTypeId
	}
	return ""
}

type ResourceTypesReaderServiceGetResourceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceId  *v2.ResourceId `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Annotations []*anypb.Any   `protobuf:"bytes,2,rep,name=annotations,proto3" json:"annotations,omitempty"`
}

func (x *ResourceTypesReaderServiceGetResourceRequest) Reset() {
	*x = ResourceTypesReaderServiceGetResourceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_c1_reader_v2_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceTypesReaderServiceGetResourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceTypesReaderServiceGetResourceRequest) ProtoMessage() {}

func (x *ResourceTypesReaderServiceGetResourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_c1_reader_v2_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceTypesReaderServiceGetResourceRequest.ProtoReflect.Descriptor instead.
func (*ResourceTypesReaderServiceGetResourceRequest) Descriptor() ([]byte, []int) {
	return file_c1_reader_v2_resource_proto_rawDescGZIP(), []int{1}
}

func (x *ResourceTypesReaderServiceGetResourceRequest) GetResourceId() *v2.ResourceId {
	if x != nil {
		return x.ResourceId
	}
	return nil
}

func (x *ResourceTypesReaderServiceGetResourceRequest) GetAnnotations() []*anypb.Any {
	if x != nil {
		return x.Annotations
	}
	return nil
}

var File_c1_reader_v2_resource_proto protoreflect.FileDescriptor

var file_c1_reader_v2_resource_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x63, 0x31, 0x2f, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63,
	0x31, 0x2e, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x1a, 0x1e, 0x63, 0x31, 0x2f,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x32, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x5c, 0x0a, 0x30, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73,
	0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x22, 0xae, 0x01,
	0x0a, 0x2c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x46,
	0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x32, 0x8e,
	0x01, 0x0a, 0x1a, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73,
	0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x70, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x3e, 0x2e, 0x63, 0x31, 0x2e, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x63, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x76, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x32,
	0x7e, 0x0a, 0x16, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x64, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x3a, 0x2e, 0x63, 0x31, 0x2e, 0x72, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x31, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42,
	0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f,
	0x6e, 0x64, 0x75, 0x63, 0x74, 0x6f, 0x72, 0x6f, 0x6e, 0x65, 0x2f, 0x62, 0x61, 0x74, 0x6f, 0x6e,
	0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x70, 0x62, 0x2f, 0x63, 0x31, 0x2f, 0x72, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x2f, 0x76, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_c1_reader_v2_resource_proto_rawDescOnce sync.Once
	file_c1_reader_v2_resource_proto_rawDescData = file_c1_reader_v2_resource_proto_rawDesc
)

func file_c1_reader_v2_resource_proto_rawDescGZIP() []byte {
	file_c1_reader_v2_resource_proto_rawDescOnce.Do(func() {
		file_c1_reader_v2_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_c1_reader_v2_resource_proto_rawDescData)
	})
	return file_c1_reader_v2_resource_proto_rawDescData
}

var file_c1_reader_v2_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_c1_reader_v2_resource_proto_goTypes = []interface{}{
	(*ResourceTypesReaderServiceGetResourceTypeRequest)(nil), // 0: c1.reader.v2.ResourceTypesReaderServiceGetResourceTypeRequest
	(*ResourceTypesReaderServiceGetResourceRequest)(nil),     // 1: c1.reader.v2.ResourceTypesReaderServiceGetResourceRequest
	(*v2.ResourceId)(nil),   // 2: c1.connector.v2.ResourceId
	(*anypb.Any)(nil),       // 3: google.protobuf.Any
	(*v2.ResourceType)(nil), // 4: c1.connector.v2.ResourceType
	(*v2.Resource)(nil),     // 5: c1.connector.v2.Resource
}
var file_c1_reader_v2_resource_proto_depIdxs = []int32{
	2, // 0: c1.reader.v2.ResourceTypesReaderServiceGetResourceRequest.resource_id:type_name -> c1.connector.v2.ResourceId
	3, // 1: c1.reader.v2.ResourceTypesReaderServiceGetResourceRequest.annotations:type_name -> google.protobuf.Any
	0, // 2: c1.reader.v2.ResourceTypesReaderService.GetResourceType:input_type -> c1.reader.v2.ResourceTypesReaderServiceGetResourceTypeRequest
	1, // 3: c1.reader.v2.ResourcesReaderService.GetResource:input_type -> c1.reader.v2.ResourceTypesReaderServiceGetResourceRequest
	4, // 4: c1.reader.v2.ResourceTypesReaderService.GetResourceType:output_type -> c1.connector.v2.ResourceType
	5, // 5: c1.reader.v2.ResourcesReaderService.GetResource:output_type -> c1.connector.v2.Resource
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_c1_reader_v2_resource_proto_init() }
func file_c1_reader_v2_resource_proto_init() {
	if File_c1_reader_v2_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_c1_reader_v2_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceTypesReaderServiceGetResourceTypeRequest); i {
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
		file_c1_reader_v2_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceTypesReaderServiceGetResourceRequest); i {
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
			RawDescriptor: file_c1_reader_v2_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_c1_reader_v2_resource_proto_goTypes,
		DependencyIndexes: file_c1_reader_v2_resource_proto_depIdxs,
		MessageInfos:      file_c1_reader_v2_resource_proto_msgTypes,
	}.Build()
	File_c1_reader_v2_resource_proto = out.File
	file_c1_reader_v2_resource_proto_rawDesc = nil
	file_c1_reader_v2_resource_proto_goTypes = nil
	file_c1_reader_v2_resource_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v6.31.1
// source: v1/sensor/sensor.proto

package sensorpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents a single sensor value.
type SensorValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name or identifier of the sensor.
	SensorName string `protobuf:"bytes,1,opt,name=sensor_name,json=sensorName,proto3" json:"sensor_name,omitempty"`
	// Represents a single sensor measurement.
	Measurement *SensorValue_Measurement `protobuf:"bytes,2,opt,name=measurement,proto3" json:"measurement,omitempty"`
}

func (x *SensorValue) Reset() {
	*x = SensorValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_sensor_sensor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SensorValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SensorValue) ProtoMessage() {}

func (x *SensorValue) ProtoReflect() protoreflect.Message {
	mi := &file_v1_sensor_sensor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SensorValue.ProtoReflect.Descriptor instead.
func (*SensorValue) Descriptor() ([]byte, []int) {
	return file_v1_sensor_sensor_proto_rawDescGZIP(), []int{0}
}

func (x *SensorValue) GetSensorName() string {
	if x != nil {
		return x.SensorName
	}
	return ""
}

func (x *SensorValue) GetMeasurement() *SensorValue_Measurement {
	if x != nil {
		return x.Measurement
	}
	return nil
}

// Request contains list of sensor values.
type SensorValuesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*SensorValue `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SensorValuesRequest) Reset() {
	*x = SensorValuesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_sensor_sensor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SensorValuesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SensorValuesRequest) ProtoMessage() {}

func (x *SensorValuesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_sensor_sensor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SensorValuesRequest.ProtoReflect.Descriptor instead.
func (*SensorValuesRequest) Descriptor() ([]byte, []int) {
	return file_v1_sensor_sensor_proto_rawDescGZIP(), []int{1}
}

func (x *SensorValuesRequest) GetItems() []*SensorValue {
	if x != nil {
		return x.Items
	}
	return nil
}

type SensorValuesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SensorValuesResponse) Reset() {
	*x = SensorValuesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_sensor_sensor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SensorValuesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SensorValuesResponse) ProtoMessage() {}

func (x *SensorValuesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_sensor_sensor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SensorValuesResponse.ProtoReflect.Descriptor instead.
func (*SensorValuesResponse) Descriptor() ([]byte, []int) {
	return file_v1_sensor_sensor_proto_rawDescGZIP(), []int{2}
}

type SensorValue_Measurement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Measurement value
	// Notice: int64 wrapper is used to distinguish the absence of field and its default zero value.
	SensorValue *wrapperspb.Int64Value `protobuf:"bytes,2,opt,name=sensor_value,json=sensorValue,proto3" json:"sensor_value,omitempty"`
	// Measurement timestamp
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *SensorValue_Measurement) Reset() {
	*x = SensorValue_Measurement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_sensor_sensor_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SensorValue_Measurement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SensorValue_Measurement) ProtoMessage() {}

func (x *SensorValue_Measurement) ProtoReflect() protoreflect.Message {
	mi := &file_v1_sensor_sensor_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SensorValue_Measurement.ProtoReflect.Descriptor instead.
func (*SensorValue_Measurement) Descriptor() ([]byte, []int) {
	return file_v1_sensor_sensor_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SensorValue_Measurement) GetSensorValue() *wrapperspb.Int64Value {
	if x != nil {
		return x.SensorValue
	}
	return nil
}

func (x *SensorValue_Measurement) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_v1_sensor_sensor_proto protoreflect.FileDescriptor

var file_v1_sensor_sensor_proto_rawDesc = []byte{
	0x0a, 0x16, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x2f, 0x73, 0x65, 0x6e, 0x73,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72,
	0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xff, 0x01, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x6e, 0x73, 0x6f,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x44, 0x0a, 0x0b, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x73, 0x65, 0x6e,
	0x73, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x2e, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b,
	0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x88, 0x01, 0x0a, 0x0b,
	0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x3e, 0x0a, 0x0c, 0x73,
	0x65, 0x6e, 0x73, 0x6f, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b,
	0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x43, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73,
	0x65, 0x6e, 0x73, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x16, 0x0a, 0x14, 0x53,
	0x65, 0x6e, 0x73, 0x6f, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x70, 0x61, 0x76, 0x6c, 0x6f, 0x76, 0x39, 0x33, 0x2f, 0x74, 0x65, 0x6c, 0x65,
	0x6d, 0x65, 0x74, 0x72, 0x79, 0x2d, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x3b, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_sensor_sensor_proto_rawDescOnce sync.Once
	file_v1_sensor_sensor_proto_rawDescData = file_v1_sensor_sensor_proto_rawDesc
)

func file_v1_sensor_sensor_proto_rawDescGZIP() []byte {
	file_v1_sensor_sensor_proto_rawDescOnce.Do(func() {
		file_v1_sensor_sensor_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_sensor_sensor_proto_rawDescData)
	})
	return file_v1_sensor_sensor_proto_rawDescData
}

var file_v1_sensor_sensor_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_v1_sensor_sensor_proto_goTypes = []interface{}{
	(*SensorValue)(nil),             // 0: sensor.v1.SensorValue
	(*SensorValuesRequest)(nil),     // 1: sensor.v1.SensorValuesRequest
	(*SensorValuesResponse)(nil),    // 2: sensor.v1.SensorValuesResponse
	(*SensorValue_Measurement)(nil), // 3: sensor.v1.SensorValue.Measurement
	(*wrapperspb.Int64Value)(nil),   // 4: google.protobuf.Int64Value
	(*timestamppb.Timestamp)(nil),   // 5: google.protobuf.Timestamp
}
var file_v1_sensor_sensor_proto_depIdxs = []int32{
	3, // 0: sensor.v1.SensorValue.measurement:type_name -> sensor.v1.SensorValue.Measurement
	0, // 1: sensor.v1.SensorValuesRequest.items:type_name -> sensor.v1.SensorValue
	4, // 2: sensor.v1.SensorValue.Measurement.sensor_value:type_name -> google.protobuf.Int64Value
	5, // 3: sensor.v1.SensorValue.Measurement.created_at:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_v1_sensor_sensor_proto_init() }
func file_v1_sensor_sensor_proto_init() {
	if File_v1_sensor_sensor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_sensor_sensor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SensorValue); i {
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
		file_v1_sensor_sensor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SensorValuesRequest); i {
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
		file_v1_sensor_sensor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SensorValuesResponse); i {
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
		file_v1_sensor_sensor_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SensorValue_Measurement); i {
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
			RawDescriptor: file_v1_sensor_sensor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_sensor_sensor_proto_goTypes,
		DependencyIndexes: file_v1_sensor_sensor_proto_depIdxs,
		MessageInfos:      file_v1_sensor_sensor_proto_msgTypes,
	}.Build()
	File_v1_sensor_sensor_proto = out.File
	file_v1_sensor_sensor_proto_rawDesc = nil
	file_v1_sensor_sensor_proto_goTypes = nil
	file_v1_sensor_sensor_proto_depIdxs = nil
}

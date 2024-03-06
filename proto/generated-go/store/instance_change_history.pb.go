// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: store/instance_change_history.proto

package store

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

type InstanceChangeHistoryPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PushEvent        *PushEvent        `protobuf:"bytes,1,opt,name=push_event,json=pushEvent,proto3" json:"push_event,omitempty"`
	ChangedResources *ChangedResources `protobuf:"bytes,2,opt,name=changed_resources,json=changedResources,proto3" json:"changed_resources,omitempty"`
}

func (x *InstanceChangeHistoryPayload) Reset() {
	*x = InstanceChangeHistoryPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_instance_change_history_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceChangeHistoryPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceChangeHistoryPayload) ProtoMessage() {}

func (x *InstanceChangeHistoryPayload) ProtoReflect() protoreflect.Message {
	mi := &file_store_instance_change_history_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceChangeHistoryPayload.ProtoReflect.Descriptor instead.
func (*InstanceChangeHistoryPayload) Descriptor() ([]byte, []int) {
	return file_store_instance_change_history_proto_rawDescGZIP(), []int{0}
}

func (x *InstanceChangeHistoryPayload) GetPushEvent() *PushEvent {
	if x != nil {
		return x.PushEvent
	}
	return nil
}

func (x *InstanceChangeHistoryPayload) GetChangedResources() *ChangedResources {
	if x != nil {
		return x.ChangedResources
	}
	return nil
}

type ChangedResources struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Databases []*ChangedResourceDatabase `protobuf:"bytes,1,rep,name=databases,proto3" json:"databases,omitempty"`
}

func (x *ChangedResources) Reset() {
	*x = ChangedResources{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_instance_change_history_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangedResources) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResources) ProtoMessage() {}

func (x *ChangedResources) ProtoReflect() protoreflect.Message {
	mi := &file_store_instance_change_history_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResources.ProtoReflect.Descriptor instead.
func (*ChangedResources) Descriptor() ([]byte, []int) {
	return file_store_instance_change_history_proto_rawDescGZIP(), []int{1}
}

func (x *ChangedResources) GetDatabases() []*ChangedResourceDatabase {
	if x != nil {
		return x.Databases
	}
	return nil
}

type ChangedResourceDatabase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string                   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Schemas []*ChangedResourceSchema `protobuf:"bytes,2,rep,name=schemas,proto3" json:"schemas,omitempty"`
}

func (x *ChangedResourceDatabase) Reset() {
	*x = ChangedResourceDatabase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_instance_change_history_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangedResourceDatabase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceDatabase) ProtoMessage() {}

func (x *ChangedResourceDatabase) ProtoReflect() protoreflect.Message {
	mi := &file_store_instance_change_history_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceDatabase.ProtoReflect.Descriptor instead.
func (*ChangedResourceDatabase) Descriptor() ([]byte, []int) {
	return file_store_instance_change_history_proto_rawDescGZIP(), []int{2}
}

func (x *ChangedResourceDatabase) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceDatabase) GetSchemas() []*ChangedResourceSchema {
	if x != nil {
		return x.Schemas
	}
	return nil
}

type ChangedResourceSchema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Tables []*ChangedResourceTable `protobuf:"bytes,2,rep,name=tables,proto3" json:"tables,omitempty"`
}

func (x *ChangedResourceSchema) Reset() {
	*x = ChangedResourceSchema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_instance_change_history_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangedResourceSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceSchema) ProtoMessage() {}

func (x *ChangedResourceSchema) ProtoReflect() protoreflect.Message {
	mi := &file_store_instance_change_history_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceSchema.ProtoReflect.Descriptor instead.
func (*ChangedResourceSchema) Descriptor() ([]byte, []int) {
	return file_store_instance_change_history_proto_rawDescGZIP(), []int{3}
}

func (x *ChangedResourceSchema) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceSchema) GetTables() []*ChangedResourceTable {
	if x != nil {
		return x.Tables
	}
	return nil
}

type ChangedResourceTable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// estimated row count of the table
	TableRows int64 `protobuf:"varint,2,opt,name=table_rows,json=tableRows,proto3" json:"table_rows,omitempty"`
}

func (x *ChangedResourceTable) Reset() {
	*x = ChangedResourceTable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_instance_change_history_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangedResourceTable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceTable) ProtoMessage() {}

func (x *ChangedResourceTable) ProtoReflect() protoreflect.Message {
	mi := &file_store_instance_change_history_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceTable.ProtoReflect.Descriptor instead.
func (*ChangedResourceTable) Descriptor() ([]byte, []int) {
	return file_store_instance_change_history_proto_rawDescGZIP(), []int{4}
}

func (x *ChangedResourceTable) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceTable) GetTableRows() int64 {
	if x != nil {
		return x.TableRows
	}
	return 0
}

var File_store_instance_change_history_proto protoreflect.FileDescriptor

var file_store_instance_change_history_proto_rawDesc = []byte{
	0x0a, 0x23, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x0f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x63, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x01, 0x0a, 0x1c, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x38, 0x0a, 0x0a, 0x70, 0x75, 0x73, 0x68, 0x5f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x62, 0x79,
	0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x70, 0x75, 0x73, 0x68, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x4d, 0x0a, 0x11, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x62,
	0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x10,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x22, 0x59, 0x0a, 0x10, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x12, 0x45, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x73, 0x22, 0x6e, 0x0a, 0x17, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x73, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x62, 0x79,
	0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x22, 0x69, 0x0a, 0x15, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3c, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62,
	0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x06,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x22, 0x49, 0x0a, 0x14, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x77, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x6f, 0x77,
	0x73, 0x42, 0x14, 0x5a, 0x12, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x67,
	0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_instance_change_history_proto_rawDescOnce sync.Once
	file_store_instance_change_history_proto_rawDescData = file_store_instance_change_history_proto_rawDesc
)

func file_store_instance_change_history_proto_rawDescGZIP() []byte {
	file_store_instance_change_history_proto_rawDescOnce.Do(func() {
		file_store_instance_change_history_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_instance_change_history_proto_rawDescData)
	})
	return file_store_instance_change_history_proto_rawDescData
}

var file_store_instance_change_history_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_store_instance_change_history_proto_goTypes = []interface{}{
	(*InstanceChangeHistoryPayload)(nil), // 0: bytebase.store.InstanceChangeHistoryPayload
	(*ChangedResources)(nil),             // 1: bytebase.store.ChangedResources
	(*ChangedResourceDatabase)(nil),      // 2: bytebase.store.ChangedResourceDatabase
	(*ChangedResourceSchema)(nil),        // 3: bytebase.store.ChangedResourceSchema
	(*ChangedResourceTable)(nil),         // 4: bytebase.store.ChangedResourceTable
	(*PushEvent)(nil),                    // 5: bytebase.store.PushEvent
}
var file_store_instance_change_history_proto_depIdxs = []int32{
	5, // 0: bytebase.store.InstanceChangeHistoryPayload.push_event:type_name -> bytebase.store.PushEvent
	1, // 1: bytebase.store.InstanceChangeHistoryPayload.changed_resources:type_name -> bytebase.store.ChangedResources
	2, // 2: bytebase.store.ChangedResources.databases:type_name -> bytebase.store.ChangedResourceDatabase
	3, // 3: bytebase.store.ChangedResourceDatabase.schemas:type_name -> bytebase.store.ChangedResourceSchema
	4, // 4: bytebase.store.ChangedResourceSchema.tables:type_name -> bytebase.store.ChangedResourceTable
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_store_instance_change_history_proto_init() }
func file_store_instance_change_history_proto_init() {
	if File_store_instance_change_history_proto != nil {
		return
	}
	file_store_vcs_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_store_instance_change_history_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceChangeHistoryPayload); i {
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
		file_store_instance_change_history_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangedResources); i {
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
		file_store_instance_change_history_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangedResourceDatabase); i {
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
		file_store_instance_change_history_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangedResourceSchema); i {
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
		file_store_instance_change_history_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangedResourceTable); i {
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
			RawDescriptor: file_store_instance_change_history_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_instance_change_history_proto_goTypes,
		DependencyIndexes: file_store_instance_change_history_proto_depIdxs,
		MessageInfos:      file_store_instance_change_history_proto_msgTypes,
	}.Build()
	File_store_instance_change_history_proto = out.File
	file_store_instance_change_history_proto_rawDesc = nil
	file_store_instance_change_history_proto_goTypes = nil
	file_store_instance_change_history_proto_depIdxs = nil
}

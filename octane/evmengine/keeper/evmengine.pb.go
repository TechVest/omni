// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: octane/evmengine/keeper/evmengine.proto

package keeper

import (
	_ "cosmossdk.io/api/cosmos/orm/v1"
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

// ExecutionHead defines the execution chain head.
// It is a singleton table; it only has a single row with ID==1.
type ExecutionHead struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                            // Auto-incremented ID (always and only 1).
	CreatedHeight uint64                 `protobuf:"varint,2,opt,name=created_height,json=createdHeight,proto3" json:"created_height,omitempty"` // Consensus chain height this execution block was created in.
	BlockHeight   uint64                 `protobuf:"varint,3,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`       // Execution block height.
	BlockHash     []byte                 `protobuf:"bytes,4,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`              // Execution block hash.
	BlockTime     uint64                 `protobuf:"varint,5,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`             // Execution block time.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecutionHead) Reset() {
	*x = ExecutionHead{}
	mi := &file_octane_evmengine_keeper_evmengine_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecutionHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutionHead) ProtoMessage() {}

func (x *ExecutionHead) ProtoReflect() protoreflect.Message {
	mi := &file_octane_evmengine_keeper_evmengine_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutionHead.ProtoReflect.Descriptor instead.
func (*ExecutionHead) Descriptor() ([]byte, []int) {
	return file_octane_evmengine_keeper_evmengine_proto_rawDescGZIP(), []int{0}
}

func (x *ExecutionHead) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ExecutionHead) GetCreatedHeight() uint64 {
	if x != nil {
		return x.CreatedHeight
	}
	return 0
}

func (x *ExecutionHead) GetBlockHeight() uint64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *ExecutionHead) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *ExecutionHead) GetBlockTime() uint64 {
	if x != nil {
		return x.BlockTime
	}
	return 0
}

// Withdrawal defines a wirthdrawal request.
type Withdrawal struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                            // Auto-incremented ID.
	Address       []byte                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`                                   // Target address of the withdrawal.
	CreatedHeight uint64                 `protobuf:"varint,3,opt,name=created_height,json=createdHeight,proto3" json:"created_height,omitempty"` // Consensus chain height this withdrawal was created in.
	AmountGwei    uint64                 `protobuf:"varint,4,opt,name=amount_gwei,json=amountGwei,proto3" json:"amount_gwei,omitempty"`          // Value of withdrawal in Gwei.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Withdrawal) Reset() {
	*x = Withdrawal{}
	mi := &file_octane_evmengine_keeper_evmengine_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Withdrawal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Withdrawal) ProtoMessage() {}

func (x *Withdrawal) ProtoReflect() protoreflect.Message {
	mi := &file_octane_evmengine_keeper_evmengine_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Withdrawal.ProtoReflect.Descriptor instead.
func (*Withdrawal) Descriptor() ([]byte, []int) {
	return file_octane_evmengine_keeper_evmengine_proto_rawDescGZIP(), []int{1}
}

func (x *Withdrawal) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Withdrawal) GetAddress() []byte {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Withdrawal) GetCreatedHeight() uint64 {
	if x != nil {
		return x.CreatedHeight
	}
	return 0
}

func (x *Withdrawal) GetAmountGwei() uint64 {
	if x != nil {
		return x.AmountGwei
	}
	return 0
}

var File_octane_evmengine_keeper_evmengine_proto protoreflect.FileDescriptor

var file_octane_evmengine_keeper_evmengine_proto_rawDesc = string([]byte{
	0x0a, 0x27, 0x6f, 0x63, 0x74, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x6d, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x65, 0x76, 0x6d, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x6f, 0x63, 0x74, 0x61, 0x6e,
	0x65, 0x2e, 0x65, 0x76, 0x6d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x1a, 0x17, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x6f, 0x72, 0x6d, 0x2f, 0x76,
	0x31, 0x2f, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb9, 0x01, 0x0a, 0x0d,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x54, 0x69, 0x6d, 0x65, 0x3a, 0x10, 0xf2, 0x9e, 0xd3, 0x8e, 0x03, 0x0a, 0x0a, 0x06, 0x0a,
	0x02, 0x69, 0x64, 0x10, 0x01, 0x18, 0x01, 0x22, 0x90, 0x01, 0x0a, 0x0a, 0x57, 0x69, 0x74, 0x68,
	0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x25, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x68, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x67, 0x77, 0x65, 0x69, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x47, 0x77, 0x65, 0x69, 0x3a, 0x10, 0xf2, 0x9e, 0xd3, 0x8e, 0x03, 0x0a,
	0x0a, 0x06, 0x0a, 0x02, 0x69, 0x64, 0x10, 0x01, 0x18, 0x02, 0x42, 0xe1, 0x01, 0x0a, 0x1b, 0x63,
	0x6f, 0x6d, 0x2e, 0x6f, 0x63, 0x74, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x6d, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x42, 0x0e, 0x45, 0x76, 0x6d, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2d, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x6f, 0x63, 0x74, 0x61, 0x6e,
	0x65, 0x2f, 0x65, 0x76, 0x6d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0xa2, 0x02, 0x03, 0x4f, 0x45, 0x4b, 0xaa, 0x02, 0x17, 0x4f, 0x63, 0x74, 0x61, 0x6e,
	0x65, 0x2e, 0x45, 0x76, 0x6d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0xca, 0x02, 0x17, 0x4f, 0x63, 0x74, 0x61, 0x6e, 0x65, 0x5c, 0x45, 0x76, 0x6d, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5c, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0xe2, 0x02, 0x23, 0x4f,
	0x63, 0x74, 0x61, 0x6e, 0x65, 0x5c, 0x45, 0x76, 0x6d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5c,
	0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x19, 0x4f, 0x63, 0x74, 0x61, 0x6e, 0x65, 0x3a, 0x3a, 0x45, 0x76, 0x6d,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x3a, 0x3a, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_octane_evmengine_keeper_evmengine_proto_rawDescOnce sync.Once
	file_octane_evmengine_keeper_evmengine_proto_rawDescData []byte
)

func file_octane_evmengine_keeper_evmengine_proto_rawDescGZIP() []byte {
	file_octane_evmengine_keeper_evmengine_proto_rawDescOnce.Do(func() {
		file_octane_evmengine_keeper_evmengine_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_octane_evmengine_keeper_evmengine_proto_rawDesc), len(file_octane_evmengine_keeper_evmengine_proto_rawDesc)))
	})
	return file_octane_evmengine_keeper_evmengine_proto_rawDescData
}

var file_octane_evmengine_keeper_evmengine_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_octane_evmengine_keeper_evmengine_proto_goTypes = []any{
	(*ExecutionHead)(nil), // 0: octane.evmengine.keeper.ExecutionHead
	(*Withdrawal)(nil),    // 1: octane.evmengine.keeper.Withdrawal
}
var file_octane_evmengine_keeper_evmengine_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_octane_evmengine_keeper_evmengine_proto_init() }
func file_octane_evmengine_keeper_evmengine_proto_init() {
	if File_octane_evmengine_keeper_evmengine_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_octane_evmengine_keeper_evmengine_proto_rawDesc), len(file_octane_evmengine_keeper_evmengine_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_octane_evmengine_keeper_evmengine_proto_goTypes,
		DependencyIndexes: file_octane_evmengine_keeper_evmengine_proto_depIdxs,
		MessageInfos:      file_octane_evmengine_keeper_evmengine_proto_msgTypes,
	}.Build()
	File_octane_evmengine_keeper_evmengine_proto = out.File
	file_octane_evmengine_keeper_evmengine_proto_goTypes = nil
	file_octane_evmengine_keeper_evmengine_proto_depIdxs = nil
}

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: task.proto

package proto

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TaskType int32

const (
	TaskType_SHELL TaskType = 0
	TaskType_JSON  TaskType = 2
)

var TaskType_name = map[int32]string{
	0: "SHELL",
	2: "JSON",
}

var TaskType_value = map[string]int32{
	"SHELL": 0,
	"JSON":  2,
}

func (x TaskType) String() string {
	return proto.EnumName(TaskType_name, int32(x))
}

func (TaskType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}

type TaskStatus int32

const (
	TaskStatus_TASK_SUCCESS TaskStatus = 0
	TaskStatus_TASK_FAILED  TaskStatus = 1
	TaskStatus_TASK_RUNNING TaskStatus = 2
)

var TaskStatus_name = map[int32]string{
	0: "TASK_SUCCESS",
	1: "TASK_FAILED",
	2: "TASK_RUNNING",
}

var TaskStatus_value = map[string]int32{
	"TASK_SUCCESS": 0,
	"TASK_FAILED":  1,
	"TASK_RUNNING": 2,
}

func (x TaskStatus) String() string {
	return proto.EnumName(TaskStatus_name, int32(x))
}

func (TaskStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}

type RunTaskRequest struct {
	Type                 TaskType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.TaskType" json:"type,omitempty"`
	Data                 string   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Id                   string   `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunTaskRequest) Reset()         { *m = RunTaskRequest{} }
func (m *RunTaskRequest) String() string { return proto.CompactTextString(m) }
func (*RunTaskRequest) ProtoMessage()    {}
func (*RunTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}
func (m *RunTaskRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RunTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RunTaskRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RunTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunTaskRequest.Merge(m, src)
}
func (m *RunTaskRequest) XXX_Size() int {
	return m.Size()
}
func (m *RunTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RunTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RunTaskRequest proto.InternalMessageInfo

func (m *RunTaskRequest) GetType() TaskType {
	if m != nil {
		return m.Type
	}
	return TaskType_SHELL
}

func (m *RunTaskRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *RunTaskRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type RunTaskResponse struct {
	Status               TaskStatus `protobuf:"varint,1,opt,name=status,proto3,enum=proto.TaskStatus" json:"status,omitempty"`
	Data                 string     `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Id                   string     `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RunTaskResponse) Reset()         { *m = RunTaskResponse{} }
func (m *RunTaskResponse) String() string { return proto.CompactTextString(m) }
func (*RunTaskResponse) ProtoMessage()    {}
func (*RunTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}
func (m *RunTaskResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RunTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RunTaskResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RunTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunTaskResponse.Merge(m, src)
}
func (m *RunTaskResponse) XXX_Size() int {
	return m.Size()
}
func (m *RunTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RunTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RunTaskResponse proto.InternalMessageInfo

func (m *RunTaskResponse) GetStatus() TaskStatus {
	if m != nil {
		return m.Status
	}
	return TaskStatus_TASK_SUCCESS
}

func (m *RunTaskResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *RunTaskResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type StreamRunTaskResponse struct {
	Status               TaskStatus `protobuf:"varint,1,opt,name=status,proto3,enum=proto.TaskStatus" json:"status,omitempty"`
	Data                 string     `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Id                   string     `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *StreamRunTaskResponse) Reset()         { *m = StreamRunTaskResponse{} }
func (m *StreamRunTaskResponse) String() string { return proto.CompactTextString(m) }
func (*StreamRunTaskResponse) ProtoMessage()    {}
func (*StreamRunTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{2}
}
func (m *StreamRunTaskResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StreamRunTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StreamRunTaskResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StreamRunTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamRunTaskResponse.Merge(m, src)
}
func (m *StreamRunTaskResponse) XXX_Size() int {
	return m.Size()
}
func (m *StreamRunTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamRunTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamRunTaskResponse proto.InternalMessageInfo

func (m *StreamRunTaskResponse) GetStatus() TaskStatus {
	if m != nil {
		return m.Status
	}
	return TaskStatus_TASK_SUCCESS
}

func (m *StreamRunTaskResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *StreamRunTaskResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterEnum("proto.TaskType", TaskType_name, TaskType_value)
	proto.RegisterEnum("proto.TaskStatus", TaskStatus_name, TaskStatus_value)
	proto.RegisterType((*RunTaskRequest)(nil), "proto.RunTaskRequest")
	proto.RegisterType((*RunTaskResponse)(nil), "proto.RunTaskResponse")
	proto.RegisterType((*StreamRunTaskResponse)(nil), "proto.StreamRunTaskResponse")
}

func init() { proto.RegisterFile("task.proto", fileDescriptor_ce5d8dd45b4a91ff) }

var fileDescriptor_ce5d8dd45b4a91ff = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2c, 0xce,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x91, 0x5c, 0x7c, 0x41, 0xa5,
	0x79, 0x21, 0x89, 0xc5, 0xd9, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0xca, 0x5c, 0x2c,
	0x25, 0x95, 0x05, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x7c, 0x46, 0xfc, 0x10, 0xe5, 0x7a, 0x20,
	0x15, 0x21, 0x95, 0x05, 0xa9, 0x41, 0x60, 0x49, 0x21, 0x21, 0x2e, 0x96, 0x94, 0xc4, 0x92, 0x44,
	0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x30, 0x5b, 0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82,
	0x19, 0x2c, 0xc2, 0x94, 0x99, 0xa2, 0x94, 0xc0, 0xc5, 0x0f, 0x37, 0xba, 0xb8, 0x20, 0x3f, 0xaf,
	0x38, 0x55, 0x48, 0x93, 0x8b, 0xad, 0xb8, 0x24, 0xb1, 0xa4, 0xb4, 0x18, 0x6a, 0xba, 0x20, 0x92,
	0xe9, 0xc1, 0x60, 0x89, 0x20, 0xa8, 0x02, 0xa2, 0x6c, 0x48, 0xe3, 0x12, 0x0d, 0x2e, 0x29, 0x4a,
	0x4d, 0xcc, 0xa5, 0xad, 0x3d, 0x5a, 0xf2, 0x5c, 0x1c, 0x30, 0xff, 0x0b, 0x71, 0x72, 0xb1, 0x06,
	0x7b, 0xb8, 0xfa, 0xf8, 0x08, 0x30, 0x08, 0x71, 0x70, 0xb1, 0x78, 0x05, 0xfb, 0xfb, 0x09, 0x30,
	0x69, 0x39, 0x72, 0x71, 0x21, 0x8c, 0x16, 0x12, 0xe0, 0xe2, 0x09, 0x71, 0x0c, 0xf6, 0x8e, 0x0f,
	0x0e, 0x75, 0x76, 0x76, 0x0d, 0x0e, 0x16, 0x60, 0x10, 0xe2, 0xe7, 0xe2, 0x06, 0x8b, 0xb8, 0x39,
	0x7a, 0xfa, 0xb8, 0xba, 0x08, 0x30, 0xc2, 0x95, 0x04, 0x85, 0xfa, 0xf9, 0x79, 0xfa, 0xb9, 0x0b,
	0x30, 0x19, 0xf5, 0x33, 0x72, 0x71, 0x83, 0xcd, 0x48, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x15, 0xb2,
	0xe0, 0x62, 0x87, 0xfa, 0x4a, 0x48, 0x14, 0xea, 0x7a, 0xd4, 0x88, 0x92, 0x12, 0x43, 0x17, 0x86,
	0x7a, 0xde, 0x8d, 0x8b, 0x37, 0xa8, 0x34, 0x0f, 0x12, 0x30, 0xf8, 0xf4, 0xcb, 0x40, 0x85, 0xb1,
	0x06, 0xa1, 0x01, 0xa3, 0x93, 0xc0, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78,
	0x24, 0xc7, 0x38, 0xe3, 0xb1, 0x1c, 0x43, 0x12, 0x1b, 0x58, 0x83, 0x31, 0x20, 0x00, 0x00, 0xff,
	0xff, 0x9e, 0x55, 0x05, 0xc6, 0x48, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskServiceClient interface {
	RunTask(ctx context.Context, in *RunTaskRequest, opts ...grpc.CallOption) (*RunTaskResponse, error)
	RunStreamTask(ctx context.Context, in *RunTaskRequest, opts ...grpc.CallOption) (TaskService_RunStreamTaskClient, error)
}

type taskServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskServiceClient(cc *grpc.ClientConn) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) RunTask(ctx context.Context, in *RunTaskRequest, opts ...grpc.CallOption) (*RunTaskResponse, error) {
	out := new(RunTaskResponse)
	err := c.cc.Invoke(ctx, "/proto.TaskService/RunTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) RunStreamTask(ctx context.Context, in *RunTaskRequest, opts ...grpc.CallOption) (TaskService_RunStreamTaskClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TaskService_serviceDesc.Streams[0], "/proto.TaskService/RunStreamTask", opts...)
	if err != nil {
		return nil, err
	}
	x := &taskServiceRunStreamTaskClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TaskService_RunStreamTaskClient interface {
	Recv() (*StreamRunTaskResponse, error)
	grpc.ClientStream
}

type taskServiceRunStreamTaskClient struct {
	grpc.ClientStream
}

func (x *taskServiceRunStreamTaskClient) Recv() (*StreamRunTaskResponse, error) {
	m := new(StreamRunTaskResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskServiceServer is the server API for TaskService service.
type TaskServiceServer interface {
	RunTask(context.Context, *RunTaskRequest) (*RunTaskResponse, error)
	RunStreamTask(*RunTaskRequest, TaskService_RunStreamTaskServer) error
}

// UnimplementedTaskServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (*UnimplementedTaskServiceServer) RunTask(ctx context.Context, req *RunTaskRequest) (*RunTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunTask not implemented")
}
func (*UnimplementedTaskServiceServer) RunStreamTask(req *RunTaskRequest, srv TaskService_RunStreamTaskServer) error {
	return status.Errorf(codes.Unimplemented, "method RunStreamTask not implemented")
}

func RegisterTaskServiceServer(s *grpc.Server, srv TaskServiceServer) {
	s.RegisterService(&_TaskService_serviceDesc, srv)
}

func _TaskService_RunTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).RunTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskService/RunTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).RunTask(ctx, req.(*RunTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_RunStreamTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RunTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServiceServer).RunStreamTask(m, &taskServiceRunStreamTaskServer{stream})
}

type TaskService_RunStreamTaskServer interface {
	Send(*StreamRunTaskResponse) error
	grpc.ServerStream
}

type taskServiceRunStreamTaskServer struct {
	grpc.ServerStream
}

func (x *taskServiceRunStreamTaskServer) Send(m *StreamRunTaskResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _TaskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunTask",
			Handler:    _TaskService_RunTask_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RunStreamTask",
			Handler:       _TaskService_RunStreamTask_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "task.proto",
}

func (m *RunTaskRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RunTaskRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RunTaskRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTask(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTask(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x12
	}
	if m.Type != 0 {
		i = encodeVarintTask(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RunTaskResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RunTaskResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RunTaskResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTask(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTask(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x12
	}
	if m.Status != 0 {
		i = encodeVarintTask(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *StreamRunTaskResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StreamRunTaskResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StreamRunTaskResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTask(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTask(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x12
	}
	if m.Status != 0 {
		i = encodeVarintTask(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTask(dAtA []byte, offset int, v uint64) int {
	offset -= sovTask(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RunTaskRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovTask(uint64(m.Type))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTask(uint64(l))
	}
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTask(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *RunTaskResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovTask(uint64(m.Status))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTask(uint64(l))
	}
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTask(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *StreamRunTaskResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovTask(uint64(m.Status))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTask(uint64(l))
	}
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTask(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTask(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTask(x uint64) (n int) {
	return sovTask(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RunTaskRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTask
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RunTaskRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RunTaskRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= TaskType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTask
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTask
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTask
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTask
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTask(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTask
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RunTaskResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTask
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RunTaskResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RunTaskResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= TaskStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTask
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTask
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTask
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTask
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTask(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTask
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *StreamRunTaskResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTask
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StreamRunTaskResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamRunTaskResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= TaskStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTask
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTask
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTask
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTask
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTask
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTask(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTask
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTask(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTask
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTask
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTask
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTask
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTask
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTask
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTask        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTask          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTask = fmt.Errorf("proto: unexpected end of group")
)

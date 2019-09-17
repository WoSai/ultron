// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ultron.proto

package ultron

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message_Type int32

const (
	Message_Disconnect          Message_Type = 0
	Message_RefreshRunnerConfig Message_Type = 1
	Message_StartAttack         Message_Type = 2
	Message_StopAttack          Message_Type = 3
)

var Message_Type_name = map[int32]string{
	0: "Disconnect",
	1: "RefreshRunnerConfig",
	2: "StartAttack",
	3: "StopAttack",
}

var Message_Type_value = map[string]int32{
	"Disconnect":          0,
	"RefreshRunnerConfig": 1,
	"StartAttack":         2,
	"StopAttack":          3,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}

func (Message_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f7eac1054816c15a, []int{4, 0}
}

type Result struct {
	Name                 string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Duration             int64          `protobuf:"varint,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Error                *AttackerError `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7eac1054816c15a, []int{0}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Result) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Result) GetError() *AttackerError {
	if m != nil {
		return m.Error
	}
	return nil
}

type AttackerError struct {
	CausedBy             string   `protobuf:"bytes,1,opt,name=causedBy,proto3" json:"causedBy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AttackerError) Reset()         { *m = AttackerError{} }
func (m *AttackerError) String() string { return proto.CompactTextString(m) }
func (*AttackerError) ProtoMessage()    {}
func (*AttackerError) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7eac1054816c15a, []int{1}
}

func (m *AttackerError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttackerError.Unmarshal(m, b)
}
func (m *AttackerError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttackerError.Marshal(b, m, deterministic)
}
func (m *AttackerError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttackerError.Merge(m, src)
}
func (m *AttackerError) XXX_Size() int {
	return xxx_messageInfo_AttackerError.Size(m)
}
func (m *AttackerError) XXX_DiscardUnknown() {
	xxx_messageInfo_AttackerError.DiscardUnknown(m)
}

var xxx_messageInfo_AttackerError proto.InternalMessageInfo

func (m *AttackerError) GetCausedBy() string {
	if m != nil {
		return m.CausedBy
	}
	return ""
}

type SlaveInfo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SlaveInfo) Reset()         { *m = SlaveInfo{} }
func (m *SlaveInfo) String() string { return proto.CompactTextString(m) }
func (*SlaveInfo) ProtoMessage()    {}
func (*SlaveInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7eac1054816c15a, []int{2}
}

func (m *SlaveInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SlaveInfo.Unmarshal(m, b)
}
func (m *SlaveInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SlaveInfo.Marshal(b, m, deterministic)
}
func (m *SlaveInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SlaveInfo.Merge(m, src)
}
func (m *SlaveInfo) XXX_Size() int {
	return xxx_messageInfo_SlaveInfo.Size(m)
}
func (m *SlaveInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SlaveInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SlaveInfo proto.InternalMessageInfo

func (m *SlaveInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Ack struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}
func (*Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7eac1054816c15a, []int{3}
}

func (m *Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ack.Unmarshal(m, b)
}
func (m *Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ack.Marshal(b, m, deterministic)
}
func (m *Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ack.Merge(m, src)
}
func (m *Ack) XXX_Size() int {
	return xxx_messageInfo_Ack.Size(m)
}
func (m *Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_Ack proto.InternalMessageInfo

type Message struct {
	Type                 Message_Type `protobuf:"varint,1,opt,name=type,proto3,enum=ultron.Message_Type" json:"type,omitempty"`
	Data                 []byte       `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7eac1054816c15a, []int{4}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetType() Message_Type {
	if m != nil {
		return m.Type
	}
	return Message_Disconnect
}

func (m *Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("ultron.Message_Type", Message_Type_name, Message_Type_value)
	proto.RegisterType((*Result)(nil), "ultron.Result")
	proto.RegisterType((*AttackerError)(nil), "ultron.AttackerError")
	proto.RegisterType((*SlaveInfo)(nil), "ultron.SlaveInfo")
	proto.RegisterType((*Ack)(nil), "ultron.Ack")
	proto.RegisterType((*Message)(nil), "ultron.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UltronClient is the client API for Ultron service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UltronClient interface {
	Send(ctx context.Context, opts ...grpc.CallOption) (Ultron_SendClient, error)
	Subscribe(ctx context.Context, in *SlaveInfo, opts ...grpc.CallOption) (Ultron_SubscribeClient, error)
}

type ultronClient struct {
	cc *grpc.ClientConn
}

func NewUltronClient(cc *grpc.ClientConn) UltronClient {
	return &ultronClient{cc}
}

func (c *ultronClient) Send(ctx context.Context, opts ...grpc.CallOption) (Ultron_SendClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Ultron_serviceDesc.Streams[0], "/ultron.Ultron/Send", opts...)
	if err != nil {
		return nil, err
	}
	x := &ultronSendClient{stream}
	return x, nil
}

type Ultron_SendClient interface {
	Send(*Result) error
	CloseAndRecv() (*Ack, error)
	grpc.ClientStream
}

type ultronSendClient struct {
	grpc.ClientStream
}

func (x *ultronSendClient) Send(m *Result) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ultronSendClient) CloseAndRecv() (*Ack, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Ack)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ultronClient) Subscribe(ctx context.Context, in *SlaveInfo, opts ...grpc.CallOption) (Ultron_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Ultron_serviceDesc.Streams[1], "/ultron.Ultron/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &ultronSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Ultron_SubscribeClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type ultronSubscribeClient struct {
	grpc.ClientStream
}

func (x *ultronSubscribeClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UltronServer is the server API for Ultron service.
type UltronServer interface {
	Send(Ultron_SendServer) error
	Subscribe(*SlaveInfo, Ultron_SubscribeServer) error
}

func RegisterUltronServer(s *grpc.Server, srv UltronServer) {
	s.RegisterService(&_Ultron_serviceDesc, srv)
}

func _Ultron_Send_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UltronServer).Send(&ultronSendServer{stream})
}

type Ultron_SendServer interface {
	SendAndClose(*Ack) error
	Recv() (*Result, error)
	grpc.ServerStream
}

type ultronSendServer struct {
	grpc.ServerStream
}

func (x *ultronSendServer) SendAndClose(m *Ack) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ultronSendServer) Recv() (*Result, error) {
	m := new(Result)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Ultron_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SlaveInfo)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UltronServer).Subscribe(m, &ultronSubscribeServer{stream})
}

type Ultron_SubscribeServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type ultronSubscribeServer struct {
	grpc.ServerStream
}

func (x *ultronSubscribeServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

var _Ultron_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ultron.Ultron",
	HandlerType: (*UltronServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Send",
			Handler:       _Ultron_Send_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _Ultron_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ultron.proto",
}

func init() { proto.RegisterFile("ultron.proto", fileDescriptor_f7eac1054816c15a) }

var fileDescriptor_f7eac1054816c15a = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0xcd, 0x4e, 0xf3, 0x30,
	0x10, 0xac, 0x93, 0x34, 0xdf, 0xd7, 0x6d, 0x49, 0xcb, 0x02, 0xa2, 0x2a, 0x97, 0xca, 0x17, 0x22,
	0x55, 0xaa, 0x50, 0xfb, 0x04, 0xe5, 0xe7, 0xc0, 0x01, 0x09, 0x39, 0xf0, 0x00, 0x6e, 0xb2, 0x29,
	0x51, 0x8b, 0x1d, 0x39, 0x0e, 0x52, 0x1f, 0x85, 0xb7, 0x45, 0xf9, 0x69, 0x24, 0xb8, 0xed, 0x78,
	0xc6, 0xeb, 0x99, 0x31, 0x8c, 0xca, 0x83, 0x35, 0x5a, 0x2d, 0x73, 0xa3, 0xad, 0x46, 0xbf, 0x41,
	0x9c, 0xc0, 0x17, 0x54, 0x94, 0x07, 0x8b, 0x08, 0x9e, 0x92, 0x9f, 0x34, 0x65, 0x73, 0x16, 0x0e,
	0x44, 0x3d, 0xe3, 0x0c, 0xfe, 0x27, 0xa5, 0x91, 0x36, 0xd3, 0x6a, 0xea, 0xcc, 0x59, 0xe8, 0x8a,
	0x0e, 0xe3, 0x02, 0xfa, 0x64, 0x8c, 0x36, 0x53, 0x77, 0xce, 0xc2, 0xe1, 0xea, 0x6a, 0xd9, 0xee,
	0xdf, 0x58, 0x2b, 0xe3, 0x3d, 0x99, 0xa7, 0x8a, 0x14, 0x8d, 0x86, 0x2f, 0xe0, 0xec, 0xd7, 0x79,
	0xb5, 0x39, 0x96, 0x65, 0x41, 0xc9, 0xfd, 0xb1, 0x7d, 0xb1, 0xc3, 0xfc, 0x06, 0x06, 0xd1, 0x41,
	0x7e, 0xd1, 0xb3, 0x4a, 0x35, 0x06, 0xe0, 0x64, 0x49, 0x2b, 0x71, 0xb2, 0x84, 0xf7, 0xc1, 0xdd,
	0xc4, 0x7b, 0xfe, 0xcd, 0xe0, 0xdf, 0x0b, 0x15, 0x85, 0xdc, 0x11, 0x86, 0xe0, 0xd9, 0x63, 0xde,
	0x38, 0x0f, 0x56, 0x97, 0x27, 0x23, 0x2d, 0xbd, 0x7c, 0x3b, 0xe6, 0x24, 0x6a, 0x45, 0x95, 0x31,
	0x91, 0x56, 0xd6, 0x59, 0x46, 0xa2, 0x9e, 0xf9, 0x2b, 0x78, 0x95, 0x02, 0x03, 0x80, 0xc7, 0xac,
	0x88, 0xb5, 0x52, 0x14, 0xdb, 0x49, 0x0f, 0xaf, 0xe1, 0x42, 0x50, 0x6a, 0xa8, 0xf8, 0x10, 0xa5,
	0x52, 0x64, 0x1e, 0xb4, 0x4a, 0xb3, 0xdd, 0x84, 0xe1, 0x18, 0x86, 0x91, 0x95, 0xc6, 0x36, 0x81,
	0x26, 0x4e, 0x75, 0x33, 0xb2, 0x3a, 0x6f, 0xb1, 0xbb, 0x4a, 0xc1, 0x7f, 0xaf, 0x2d, 0xe0, 0x2d,
	0x78, 0x11, 0xa9, 0x04, 0x83, 0x93, 0xa7, 0xa6, 0xeb, 0xd9, 0xb0, 0x2b, 0x2b, 0xde, 0xf3, 0x5e,
	0xc8, 0x70, 0x0d, 0x83, 0xa8, 0xdc, 0x16, 0xb1, 0xc9, 0xb6, 0x84, 0xe7, 0x27, 0xb6, 0x6b, 0x61,
	0x36, 0xfe, 0x13, 0x8a, 0xf7, 0xee, 0xd8, 0xd6, 0xaf, 0xbf, 0x72, 0xfd, 0x13, 0x00, 0x00, 0xff,
	0xff, 0x5d, 0x8c, 0x85, 0x74, 0xda, 0x01, 0x00, 0x00,
}
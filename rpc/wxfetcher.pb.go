// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wxfetcher.proto

package rpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
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

type FetchURLRequest struct {
	OriginalUrl          string   `protobuf:"bytes,1,opt,name=originalUrl,proto3" json:"originalUrl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FetchURLRequest) Reset()         { *m = FetchURLRequest{} }
func (m *FetchURLRequest) String() string { return proto.CompactTextString(m) }
func (*FetchURLRequest) ProtoMessage()    {}
func (*FetchURLRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cc394fee09db7d6, []int{0}
}

func (m *FetchURLRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FetchURLRequest.Unmarshal(m, b)
}
func (m *FetchURLRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FetchURLRequest.Marshal(b, m, deterministic)
}
func (m *FetchURLRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FetchURLRequest.Merge(m, src)
}
func (m *FetchURLRequest) XXX_Size() int {
	return xxx_messageInfo_FetchURLRequest.Size(m)
}
func (m *FetchURLRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FetchURLRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FetchURLRequest proto.InternalMessageInfo

func (m *FetchURLRequest) GetOriginalUrl() string {
	if m != nil {
		return m.OriginalUrl
	}
	return ""
}

type FetchURLResponse struct {
	ShortenedKey         string   `protobuf:"bytes,1,opt,name=shortenedKey,proto3" json:"shortenedKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FetchURLResponse) Reset()         { *m = FetchURLResponse{} }
func (m *FetchURLResponse) String() string { return proto.CompactTextString(m) }
func (*FetchURLResponse) ProtoMessage()    {}
func (*FetchURLResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cc394fee09db7d6, []int{1}
}

func (m *FetchURLResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FetchURLResponse.Unmarshal(m, b)
}
func (m *FetchURLResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FetchURLResponse.Marshal(b, m, deterministic)
}
func (m *FetchURLResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FetchURLResponse.Merge(m, src)
}
func (m *FetchURLResponse) XXX_Size() int {
	return xxx_messageInfo_FetchURLResponse.Size(m)
}
func (m *FetchURLResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FetchURLResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FetchURLResponse proto.InternalMessageInfo

func (m *FetchURLResponse) GetShortenedKey() string {
	if m != nil {
		return m.ShortenedKey
	}
	return ""
}

func init() {
	proto.RegisterType((*FetchURLRequest)(nil), "FetchURLRequest")
	proto.RegisterType((*FetchURLResponse)(nil), "FetchURLResponse")
}

func init() { proto.RegisterFile("wxfetcher.proto", fileDescriptor_6cc394fee09db7d6) }

var fileDescriptor_6cc394fee09db7d6 = []byte{
	// 154 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xaf, 0x48, 0x4b,
	0x2d, 0x49, 0xce, 0x48, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x32, 0xe6, 0xe2, 0x77,
	0x03, 0x09, 0x84, 0x06, 0xf9, 0x04, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x29, 0x70, 0x71,
	0xe7, 0x17, 0x65, 0xa6, 0x67, 0xe6, 0x25, 0xe6, 0x84, 0x16, 0xe5, 0x48, 0x30, 0x2a, 0x30, 0x6a,
	0x70, 0x06, 0x21, 0x0b, 0x29, 0x99, 0x71, 0x09, 0x20, 0x34, 0x15, 0x17, 0xe4, 0xe7, 0x15, 0xa7,
	0x0a, 0x29, 0x71, 0xf1, 0x14, 0x67, 0xe4, 0x17, 0x95, 0xa4, 0xe6, 0xa5, 0xa6, 0x78, 0xa7, 0x56,
	0x42, 0xb5, 0xa1, 0x88, 0x19, 0xd9, 0x71, 0x71, 0x86, 0x57, 0xb8, 0x41, 0xec, 0x17, 0x32, 0xe4,
	0xe2, 0x80, 0x19, 0x22, 0x24, 0xa0, 0x87, 0xe6, 0x08, 0x29, 0x41, 0x3d, 0x74, 0x1b, 0x94, 0x18,
	0x9c, 0x58, 0xa3, 0x98, 0x8b, 0x0a, 0x92, 0x93, 0xd8, 0xc0, 0x4e, 0x37, 0x06, 0x04, 0x00, 0x00,
	0xff, 0xff, 0xca, 0x93, 0x17, 0xd7, 0xcd, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WxFetcherClient is the client API for WxFetcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WxFetcherClient interface {
	FetchURL(ctx context.Context, in *FetchURLRequest, opts ...grpc.CallOption) (*FetchURLResponse, error)
}

type wxFetcherClient struct {
	cc *grpc.ClientConn
}

func NewWxFetcherClient(cc *grpc.ClientConn) WxFetcherClient {
	return &wxFetcherClient{cc}
}

func (c *wxFetcherClient) FetchURL(ctx context.Context, in *FetchURLRequest, opts ...grpc.CallOption) (*FetchURLResponse, error) {
	out := new(FetchURLResponse)
	err := c.cc.Invoke(ctx, "/WxFetcher/FetchURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WxFetcherServer is the server API for WxFetcher service.
type WxFetcherServer interface {
	FetchURL(context.Context, *FetchURLRequest) (*FetchURLResponse, error)
}

func RegisterWxFetcherServer(s *grpc.Server, srv WxFetcherServer) {
	s.RegisterService(&_WxFetcher_serviceDesc, srv)
}

func _WxFetcher_FetchURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WxFetcherServer).FetchURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WxFetcher/FetchURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WxFetcherServer).FetchURL(ctx, req.(*FetchURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WxFetcher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "WxFetcher",
	HandlerType: (*WxFetcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchURL",
			Handler:    _WxFetcher_FetchURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wxfetcher.proto",
}
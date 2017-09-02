// Code generated by protoc-gen-go. DO NOT EDIT.
// source: all.proto

/*
Package protobufs is a generated protocol buffer package.

It is generated from these files:
	all.proto

It has these top-level messages:
	ErrorResponse
	CreateSessionRequest
	CreateSessionResponse
	CreateCompanyRequest
	CreateCompanyResponse
	Studio
	CreateStudioRequest
	CreateStudioResponse
	DeleteStudioRequest
	DeleteStudioResponse
	Staff
	CreateStaffRequest
	CreateStaffResponse
	DeleteStaffRequest
	DeleteStaffResponse
	Card
	Customer
	CreateCustomerRequest
	CreateCustomerResponse
*/
package protobufs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ErrorType int32

const (
	ErrorType_INVALID_REQUEST_FORMAT ErrorType = 0
	ErrorType_NOT_FOUND              ErrorType = 1
	ErrorType_ALREADY_EXIST          ErrorType = 2
	ErrorType_INVALID_PASSWORD       ErrorType = 3
	ErrorType_INTERNAL_SERVER_ERROR  ErrorType = 4
	ErrorType_FORBIDDEN              ErrorType = 5
)

var ErrorType_name = map[int32]string{
	0: "INVALID_REQUEST_FORMAT",
	1: "NOT_FOUND",
	2: "ALREADY_EXIST",
	3: "INVALID_PASSWORD",
	4: "INTERNAL_SERVER_ERROR",
	5: "FORBIDDEN",
}
var ErrorType_value = map[string]int32{
	"INVALID_REQUEST_FORMAT": 0,
	"NOT_FOUND":              1,
	"ALREADY_EXIST":          2,
	"INVALID_PASSWORD":       3,
	"INTERNAL_SERVER_ERROR":  4,
	"FORBIDDEN":              5,
}

func (x ErrorType) String() string {
	return proto.EnumName(ErrorType_name, int32(x))
}
func (ErrorType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ErrorResponse struct {
	ErrorType ErrorType `protobuf:"varint,1,opt,name=errorType,enum=protobufs.ErrorType" json:"errorType,omitempty"`
	Message   string    `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *ErrorResponse) Reset()                    { *m = ErrorResponse{} }
func (m *ErrorResponse) String() string            { return proto.CompactTextString(m) }
func (*ErrorResponse) ProtoMessage()               {}
func (*ErrorResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ErrorResponse) GetErrorType() ErrorType {
	if m != nil {
		return m.ErrorType
	}
	return ErrorType_INVALID_REQUEST_FORMAT
}

func (m *ErrorResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type CreateSessionRequest struct {
	EmailAddress string `protobuf:"bytes,1,opt,name=emailAddress" json:"emailAddress,omitempty"`
	Password     string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *CreateSessionRequest) Reset()                    { *m = CreateSessionRequest{} }
func (m *CreateSessionRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateSessionRequest) ProtoMessage()               {}
func (*CreateSessionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateSessionRequest) GetEmailAddress() string {
	if m != nil {
		return m.EmailAddress
	}
	return ""
}

func (m *CreateSessionRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type CreateSessionResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *CreateSessionResponse) Reset()                    { *m = CreateSessionResponse{} }
func (m *CreateSessionResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateSessionResponse) ProtoMessage()               {}
func (*CreateSessionResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CreateSessionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type CreateCompanyRequest struct {
	Name         string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	EmailAddress string `protobuf:"bytes,2,opt,name=emailAddress" json:"emailAddress,omitempty"`
	Password     string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
}

func (m *CreateCompanyRequest) Reset()                    { *m = CreateCompanyRequest{} }
func (m *CreateCompanyRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateCompanyRequest) ProtoMessage()               {}
func (*CreateCompanyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CreateCompanyRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateCompanyRequest) GetEmailAddress() string {
	if m != nil {
		return m.EmailAddress
	}
	return ""
}

func (m *CreateCompanyRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type CreateCompanyResponse struct {
	Id           int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	EmailAddress string `protobuf:"bytes,3,opt,name=emailAddress" json:"emailAddress,omitempty"`
	CreatedAt    int64  `protobuf:"varint,4,opt,name=createdAt" json:"createdAt,omitempty"`
}

func (m *CreateCompanyResponse) Reset()                    { *m = CreateCompanyResponse{} }
func (m *CreateCompanyResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateCompanyResponse) ProtoMessage()               {}
func (*CreateCompanyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *CreateCompanyResponse) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CreateCompanyResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateCompanyResponse) GetEmailAddress() string {
	if m != nil {
		return m.EmailAddress
	}
	return ""
}

func (m *CreateCompanyResponse) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

type Studio struct {
	Id        int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	CreatedAt int64  `protobuf:"varint,4,opt,name=createdAt" json:"createdAt,omitempty"`
}

func (m *Studio) Reset()                    { *m = Studio{} }
func (m *Studio) String() string            { return proto.CompactTextString(m) }
func (*Studio) ProtoMessage()               {}
func (*Studio) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Studio) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Studio) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Studio) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

type CreateStudioRequest struct {
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *CreateStudioRequest) Reset()                    { *m = CreateStudioRequest{} }
func (m *CreateStudioRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateStudioRequest) ProtoMessage()               {}
func (*CreateStudioRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreateStudioRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateStudioResponse struct {
	Studio *Studio `protobuf:"bytes,1,opt,name=studio" json:"studio,omitempty"`
}

func (m *CreateStudioResponse) Reset()                    { *m = CreateStudioResponse{} }
func (m *CreateStudioResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateStudioResponse) ProtoMessage()               {}
func (*CreateStudioResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CreateStudioResponse) GetStudio() *Studio {
	if m != nil {
		return m.Studio
	}
	return nil
}

type DeleteStudioRequest struct {
	Id int32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteStudioRequest) Reset()                    { *m = DeleteStudioRequest{} }
func (m *DeleteStudioRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteStudioRequest) ProtoMessage()               {}
func (*DeleteStudioRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *DeleteStudioRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type DeleteStudioResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *DeleteStudioResponse) Reset()                    { *m = DeleteStudioResponse{} }
func (m *DeleteStudioResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteStudioResponse) ProtoMessage()               {}
func (*DeleteStudioResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *DeleteStudioResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type Staff struct {
	Id        int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	CreatedAt int64  `protobuf:"varint,4,opt,name=createdAt" json:"createdAt,omitempty"`
}

func (m *Staff) Reset()                    { *m = Staff{} }
func (m *Staff) String() string            { return proto.CompactTextString(m) }
func (*Staff) ProtoMessage()               {}
func (*Staff) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Staff) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Staff) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Staff) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

type CreateStaffRequest struct {
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *CreateStaffRequest) Reset()                    { *m = CreateStaffRequest{} }
func (m *CreateStaffRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateStaffRequest) ProtoMessage()               {}
func (*CreateStaffRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *CreateStaffRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateStaffResponse struct {
	Staff *Staff `protobuf:"bytes,1,opt,name=staff" json:"staff,omitempty"`
}

func (m *CreateStaffResponse) Reset()                    { *m = CreateStaffResponse{} }
func (m *CreateStaffResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateStaffResponse) ProtoMessage()               {}
func (*CreateStaffResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *CreateStaffResponse) GetStaff() *Staff {
	if m != nil {
		return m.Staff
	}
	return nil
}

type DeleteStaffRequest struct {
	Id int32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteStaffRequest) Reset()                    { *m = DeleteStaffRequest{} }
func (m *DeleteStaffRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteStaffRequest) ProtoMessage()               {}
func (*DeleteStaffRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *DeleteStaffRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type DeleteStaffResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *DeleteStaffResponse) Reset()                    { *m = DeleteStaffResponse{} }
func (m *DeleteStaffResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteStaffResponse) ProtoMessage()               {}
func (*DeleteStaffResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *DeleteStaffResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type Card struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Credit int32  `protobuf:"varint,2,opt,name=credit" json:"credit,omitempty"`
}

func (m *Card) Reset()                    { *m = Card{} }
func (m *Card) String() string            { return proto.CompactTextString(m) }
func (*Card) ProtoMessage()               {}
func (*Card) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *Card) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Card) GetCredit() int32 {
	if m != nil {
		return m.Credit
	}
	return 0
}

type Customer struct {
	Id          int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
	Card        *Card  `protobuf:"bytes,4,opt,name=card" json:"card,omitempty"`
}

func (m *Customer) Reset()                    { *m = Customer{} }
func (m *Customer) String() string            { return proto.CompactTextString(m) }
func (*Customer) ProtoMessage()               {}
func (*Customer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *Customer) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Customer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Customer) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Customer) GetCard() *Card {
	if m != nil {
		return m.Card
	}
	return nil
}

type CreateCustomerRequest struct {
	Name        string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Card        *Card  `protobuf:"bytes,3,opt,name=card" json:"card,omitempty"`
}

func (m *CreateCustomerRequest) Reset()                    { *m = CreateCustomerRequest{} }
func (m *CreateCustomerRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateCustomerRequest) ProtoMessage()               {}
func (*CreateCustomerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *CreateCustomerRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateCustomerRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateCustomerRequest) GetCard() *Card {
	if m != nil {
		return m.Card
	}
	return nil
}

type CreateCustomerResponse struct {
	Customer *Customer `protobuf:"bytes,1,opt,name=customer" json:"customer,omitempty"`
}

func (m *CreateCustomerResponse) Reset()                    { *m = CreateCustomerResponse{} }
func (m *CreateCustomerResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateCustomerResponse) ProtoMessage()               {}
func (*CreateCustomerResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *CreateCustomerResponse) GetCustomer() *Customer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func init() {
	proto.RegisterType((*ErrorResponse)(nil), "protobufs.ErrorResponse")
	proto.RegisterType((*CreateSessionRequest)(nil), "protobufs.CreateSessionRequest")
	proto.RegisterType((*CreateSessionResponse)(nil), "protobufs.CreateSessionResponse")
	proto.RegisterType((*CreateCompanyRequest)(nil), "protobufs.CreateCompanyRequest")
	proto.RegisterType((*CreateCompanyResponse)(nil), "protobufs.CreateCompanyResponse")
	proto.RegisterType((*Studio)(nil), "protobufs.Studio")
	proto.RegisterType((*CreateStudioRequest)(nil), "protobufs.CreateStudioRequest")
	proto.RegisterType((*CreateStudioResponse)(nil), "protobufs.CreateStudioResponse")
	proto.RegisterType((*DeleteStudioRequest)(nil), "protobufs.DeleteStudioRequest")
	proto.RegisterType((*DeleteStudioResponse)(nil), "protobufs.DeleteStudioResponse")
	proto.RegisterType((*Staff)(nil), "protobufs.Staff")
	proto.RegisterType((*CreateStaffRequest)(nil), "protobufs.CreateStaffRequest")
	proto.RegisterType((*CreateStaffResponse)(nil), "protobufs.CreateStaffResponse")
	proto.RegisterType((*DeleteStaffRequest)(nil), "protobufs.DeleteStaffRequest")
	proto.RegisterType((*DeleteStaffResponse)(nil), "protobufs.DeleteStaffResponse")
	proto.RegisterType((*Card)(nil), "protobufs.Card")
	proto.RegisterType((*Customer)(nil), "protobufs.Customer")
	proto.RegisterType((*CreateCustomerRequest)(nil), "protobufs.CreateCustomerRequest")
	proto.RegisterType((*CreateCustomerResponse)(nil), "protobufs.CreateCustomerResponse")
	proto.RegisterEnum("protobufs.ErrorType", ErrorType_name, ErrorType_value)
}

func init() { proto.RegisterFile("all.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 605 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcf, 0x73, 0x93, 0x40,
	0x18, 0x15, 0xf2, 0xa3, 0xe1, 0xab, 0xad, 0x74, 0x9b, 0x66, 0x62, 0xc7, 0x43, 0x66, 0xfd, 0x31,
	0xa9, 0x87, 0x54, 0xe3, 0xd9, 0x03, 0x06, 0x3a, 0x83, 0x13, 0x89, 0x2e, 0x69, 0xd4, 0x83, 0x93,
	0xa1, 0xb0, 0x71, 0x70, 0x92, 0x80, 0xbb, 0x30, 0x4e, 0x0f, 0xfe, 0x01, 0xfe, 0xd7, 0x0e, 0xeb,
	0x42, 0x08, 0xda, 0xb4, 0x07, 0x4f, 0xf0, 0x7d, 0xfb, 0x78, 0xef, 0x7d, 0x6f, 0x77, 0x01, 0xcd,
	0x5b, 0x2e, 0x07, 0x31, 0x8b, 0x92, 0x08, 0x69, 0xe2, 0x71, 0x95, 0x2e, 0x38, 0xfe, 0x02, 0x07,
	0x16, 0x63, 0x11, 0x23, 0x94, 0xc7, 0xd1, 0x9a, 0x53, 0x34, 0x04, 0x8d, 0x66, 0x8d, 0xe9, 0x75,
	0x4c, 0xbb, 0x4a, 0x4f, 0xe9, 0x1f, 0x0e, 0xdb, 0x83, 0x02, 0x3f, 0xb0, 0xf2, 0x35, 0xb2, 0x81,
	0xa1, 0x2e, 0xec, 0xad, 0x28, 0xe7, 0xde, 0x57, 0xda, 0x55, 0x7b, 0x4a, 0x5f, 0x23, 0x79, 0x89,
	0x67, 0xd0, 0x1e, 0x31, 0xea, 0x25, 0xd4, 0xa5, 0x9c, 0x87, 0xd1, 0x9a, 0xd0, 0xef, 0x29, 0xe5,
	0x09, 0xc2, 0x70, 0x9f, 0xae, 0xbc, 0x70, 0x69, 0x04, 0x01, 0xa3, 0x9c, 0x0b, 0x21, 0x8d, 0x6c,
	0xf5, 0xd0, 0x29, 0xb4, 0x62, 0x8f, 0xf3, 0x1f, 0x11, 0x0b, 0x24, 0x6d, 0x51, 0xe3, 0x97, 0x70,
	0x52, 0xe1, 0x95, 0xf6, 0xbb, 0xb0, 0xc7, 0x53, 0xdf, 0xcf, 0x39, 0x5b, 0x24, 0x2f, 0xf1, 0xb7,
	0xdc, 0xca, 0x28, 0x5a, 0xc5, 0xde, 0xfa, 0x3a, 0xb7, 0x82, 0xa0, 0xbe, 0xf6, 0x56, 0x54, 0x5a,
	0x10, 0xef, 0x7f, 0xd9, 0x53, 0x6f, 0xb1, 0x57, 0xab, 0xd8, 0xfb, 0x99, 0xdb, 0x2b, 0xb4, 0xa4,
	0xbd, 0x43, 0x50, 0xc3, 0x40, 0x48, 0x35, 0x88, 0x1a, 0x06, 0x85, 0xb8, 0xba, 0x43, 0xbc, 0xf6,
	0x0f, 0xf1, 0x47, 0xa0, 0xf9, 0x42, 0x20, 0x30, 0x92, 0x6e, 0xbd, 0xa7, 0xf4, 0x6b, 0x64, 0xd3,
	0xc0, 0x6f, 0xa1, 0xe9, 0x26, 0x69, 0x10, 0x46, 0x77, 0xd2, 0xdb, 0xcd, 0x75, 0x06, 0xc7, 0x32,
	0x69, 0xc1, 0x58, 0x4d, 0xad, 0x44, 0x84, 0x8d, 0x62, 0xb3, 0x25, 0x54, 0x0e, 0x7d, 0x06, 0x4d,
	0x2e, 0x3a, 0xc2, 0xc8, 0xfe, 0xf0, 0xa8, 0x74, 0x9e, 0x24, 0x54, 0x02, 0xf0, 0x53, 0x38, 0x36,
	0xe9, 0x92, 0x56, 0xd5, 0x2a, 0x63, 0xe0, 0x17, 0xd0, 0xde, 0x86, 0xdd, 0xba, 0xfb, 0x36, 0x34,
	0xdc, 0xc4, 0x5b, 0x2c, 0xfe, 0x43, 0x22, 0x7d, 0x40, 0xf9, 0x98, 0xde, 0x62, 0xb1, 0x2b, 0x90,
	0xd7, 0x9b, 0xec, 0x04, 0x52, 0xba, 0x7c, 0x06, 0x0d, 0x9e, 0x35, 0x64, 0x1c, 0xfa, 0x56, 0x1c,
	0x19, 0xf0, 0xcf, 0x32, 0x7e, 0x02, 0x28, 0x9f, 0xb2, 0x24, 0x54, 0xcd, 0xe2, 0x7c, 0x13, 0x59,
	0x59, 0xe4, 0xe6, 0x28, 0x06, 0x50, 0x1f, 0x79, 0x2c, 0x28, 0x11, 0x69, 0x22, 0x89, 0x0e, 0x34,
	0x7d, 0x46, 0x83, 0x30, 0x11, 0x33, 0x34, 0x88, 0xac, 0x70, 0x0a, 0xad, 0x51, 0xca, 0x93, 0x68,
	0x45, 0xd9, 0x9d, 0xd2, 0xeb, 0xc1, 0x7e, 0x40, 0xb9, 0xcf, 0xc2, 0x38, 0x09, 0xa3, 0xb5, 0x3c,
	0xbe, 0xe5, 0x16, 0x7a, 0x0c, 0x75, 0xdf, 0x63, 0x81, 0x88, 0x76, 0x7f, 0xf8, 0xa0, 0x34, 0x7f,
	0x66, 0x8c, 0x88, 0x45, 0xcc, 0x8a, 0x3b, 0x24, 0xc5, 0x77, 0x5d, 0xd8, 0x8a, 0xa6, 0x7a, 0xb3,
	0x66, 0x6d, 0x97, 0xa6, 0x0d, 0x9d, 0xaa, 0xa6, 0x8c, 0xf3, 0x1c, 0x5a, 0xbe, 0xec, 0xc9, 0x6d,
	0x3b, 0x2e, 0x53, 0xe4, 0xf0, 0x02, 0xf4, 0xfc, 0x97, 0x02, 0x5a, 0xf1, 0xb3, 0x44, 0xa7, 0xd0,
	0xb1, 0x9d, 0x99, 0x31, 0xb6, 0xcd, 0x39, 0xb1, 0x3e, 0x5c, 0x5a, 0xee, 0x74, 0x7e, 0x31, 0x21,
	0xef, 0x8c, 0xa9, 0x7e, 0x0f, 0x1d, 0x80, 0xe6, 0x4c, 0xb2, 0xfa, 0xd2, 0x31, 0x75, 0x05, 0x1d,
	0xc1, 0x81, 0x31, 0x26, 0x96, 0x61, 0x7e, 0x9e, 0x5b, 0x9f, 0x6c, 0x77, 0xaa, 0xab, 0xa8, 0x0d,
	0x7a, 0xfe, 0xf5, 0x7b, 0xc3, 0x75, 0x3f, 0x4e, 0x88, 0xa9, 0xd7, 0xd0, 0x43, 0x38, 0xb1, 0x9d,
	0xa9, 0x45, 0x1c, 0x63, 0x3c, 0x77, 0x2d, 0x32, 0xb3, 0xc8, 0xdc, 0x22, 0x64, 0x42, 0xf4, 0x7a,
	0x46, 0x79, 0x31, 0x21, 0x6f, 0x6c, 0xd3, 0xb4, 0x1c, 0xbd, 0x71, 0xd5, 0x14, 0x4e, 0x5f, 0xfd,
	0x0e, 0x00, 0x00, 0xff, 0xff, 0xbb, 0xa9, 0xa6, 0x81, 0x03, 0x06, 0x00, 0x00,
}

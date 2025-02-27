// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.0
// source: proto/chat/chat.proto

package chat

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

// Enum for the type of encryption used for the message
type EncryptionType int32

const (
	EncryptionType_PLAIN  EncryptionType = 0
	EncryptionType_SIGNAL EncryptionType = 1
	EncryptionType_PREKEY EncryptionType = 2
)

// Enum value maps for EncryptionType.
var (
	EncryptionType_name = map[int32]string{
		0: "PLAIN",
		1: "SIGNAL",
		2: "PREKEY",
	}
	EncryptionType_value = map[string]int32{
		"PLAIN":  0,
		"SIGNAL": 1,
		"PREKEY": 2,
	}
)

func (x EncryptionType) Enum() *EncryptionType {
	p := new(EncryptionType)
	*p = x
	return p
}

func (x EncryptionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EncryptionType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_chat_chat_proto_enumTypes[0].Descriptor()
}

func (EncryptionType) Type() protoreflect.EnumType {
	return &file_proto_chat_chat_proto_enumTypes[0]
}

func (x EncryptionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EncryptionType.Descriptor instead.
func (EncryptionType) EnumDescriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{0}
}

// MessageRequest is used by the client to send messages or files to another user
type MessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecipientId      uint32         `protobuf:"varint,1,opt,name=recipient_id,json=recipientId,proto3" json:"recipient_id,omitempty"`                                   // Unique identifier of the recipient (user or group ID)
	EncryptedMessage []byte         `protobuf:"bytes,2,opt,name=encrypted_message,json=encryptedMessage,proto3" json:"encrypted_message,omitempty"`                     // The message content encrypted with the recipient's public key
	MessageId        string         `protobuf:"bytes,3,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`                                          // Unique identifier for this message, generated by the client
	Timestamp        string         `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                                                           // Timestamp of when the message was sent, in ISO 8601 format
	FileContent      []byte         `protobuf:"bytes,5,opt,name=file_content,json=fileContent,proto3" json:"file_content,omitempty"`                                    // (Optional) Encrypted file data being sent (if any)
	FileName         string         `protobuf:"bytes,6,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`                                             // (Optional) Original name of the file being sent (if any)
	FileType         string         `protobuf:"bytes,7,opt,name=file_type,json=fileType,proto3" json:"file_type,omitempty"`                                             // (Optional) MIME type of the file (e.g., "image/png", "application/pdf")
	FileSize         uint64         `protobuf:"varint,8,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`                                            // (Optional) Size of the file in bytes
	EncryptionType   EncryptionType `protobuf:"varint,9,opt,name=encryption_type,json=encryptionType,proto3,enum=chat.EncryptionType" json:"encryption_type,omitempty"` // Type of encryption (Plain, Signal, or PreKey)
}

func (x *MessageRequest) Reset() {
	*x = MessageRequest{}
	mi := &file_proto_chat_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageRequest) ProtoMessage() {}

func (x *MessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageRequest.ProtoReflect.Descriptor instead.
func (*MessageRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{0}
}

func (x *MessageRequest) GetRecipientId() uint32 {
	if x != nil {
		return x.RecipientId
	}
	return 0
}

func (x *MessageRequest) GetEncryptedMessage() []byte {
	if x != nil {
		return x.EncryptedMessage
	}
	return nil
}

func (x *MessageRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *MessageRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *MessageRequest) GetFileContent() []byte {
	if x != nil {
		return x.FileContent
	}
	return nil
}

func (x *MessageRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *MessageRequest) GetFileType() string {
	if x != nil {
		return x.FileType
	}
	return ""
}

func (x *MessageRequest) GetFileSize() uint64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *MessageRequest) GetEncryptionType() EncryptionType {
	if x != nil {
		return x.EncryptionType
	}
	return EncryptionType_PLAIN
}

// MessageResponse is used by the server to deliver messages to the recipient.
type MessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId         uint32         `protobuf:"varint,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`                                            // Unique identifier of the sender
	SenderUsername   string         `protobuf:"bytes,2,opt,name=sender_username,json=senderUsername,proto3" json:"sender_username,omitempty"`                           // Username of the sender
	RecipientId      uint32         `protobuf:"varint,3,opt,name=recipient_id,json=recipientId,proto3" json:"recipient_id,omitempty"`                                   // Unique identifier of the recipient
	MessageId        string         `protobuf:"bytes,4,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`                                          // The message ID of the message being acknowledged or delivered
	EncryptedMessage []byte         `protobuf:"bytes,5,opt,name=encrypted_message,json=encryptedMessage,proto3" json:"encrypted_message,omitempty"`                     // The message content encrypted with the recipient's public key
	Status           string         `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`                                                                 // Status of the message (e.g., "delivered", "read", "received")
	Timestamp        string         `protobuf:"bytes,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                                                           // Timestamp of when the server processed or delivered the message
	EncryptionType   EncryptionType `protobuf:"varint,8,opt,name=encryption_type,json=encryptionType,proto3,enum=chat.EncryptionType" json:"encryption_type,omitempty"` // Type of encryption (Plain, Signal, or PreKey)
	FileName         string         `protobuf:"bytes,9,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`                                             // (Optional) Original name of the file being sent (if any)
	FileType         string         `protobuf:"bytes,10,opt,name=file_type,json=fileType,proto3" json:"file_type,omitempty"`                                            // (Optional) MIME type of the file (e.g., "image/png", "application/pdf")
	FileSize         uint64         `protobuf:"varint,11,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`                                           // (Optional) Size of the file in bytes
}

func (x *MessageResponse) Reset() {
	*x = MessageResponse{}
	mi := &file_proto_chat_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageResponse) ProtoMessage() {}

func (x *MessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageResponse.ProtoReflect.Descriptor instead.
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{1}
}

func (x *MessageResponse) GetSenderId() uint32 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *MessageResponse) GetSenderUsername() string {
	if x != nil {
		return x.SenderUsername
	}
	return ""
}

func (x *MessageResponse) GetRecipientId() uint32 {
	if x != nil {
		return x.RecipientId
	}
	return 0
}

func (x *MessageResponse) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *MessageResponse) GetEncryptedMessage() []byte {
	if x != nil {
		return x.EncryptedMessage
	}
	return nil
}

func (x *MessageResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *MessageResponse) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *MessageResponse) GetEncryptionType() EncryptionType {
	if x != nil {
		return x.EncryptionType
	}
	return EncryptionType_PLAIN
}

func (x *MessageResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *MessageResponse) GetFileType() string {
	if x != nil {
		return x.FileType
	}
	return ""
}

func (x *MessageResponse) GetFileSize() uint64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

var File_proto_chat_chat_proto protoreflect.FileDescriptor

var file_proto_chat_chat_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68, 0x61, 0x74, 0x22, 0xd6, 0x02,
	0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x11, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x21, 0x0a,
	0x0c, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x3d, 0x0a, 0x0f, 0x65, 0x6e, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x14, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x92, 0x03, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x49, 0x64, 0x12, 0x2b, 0x0a, 0x11, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x3d, 0x0a, 0x0f, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x2a, 0x33, 0x0a, 0x0e, 0x45,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a,
	0x05, 0x50, 0x4c, 0x41, 0x49, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x49, 0x47, 0x4e,
	0x41, 0x4c, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x52, 0x45, 0x4b, 0x45, 0x59, 0x10, 0x02,
	0x32, 0x50, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x41, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x12, 0x14, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01,
	0x30, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6a, 0x6f, 0x68, 0x6e, 0x6b, 0x68, 0x6b, 0x2f, 0x63, 0x6c, 0x69, 0x5f, 0x63, 0x68, 0x61,
	0x74, 0x5f, 0x61, 0x70, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_chat_chat_proto_rawDescOnce sync.Once
	file_proto_chat_chat_proto_rawDescData = file_proto_chat_chat_proto_rawDesc
)

func file_proto_chat_chat_proto_rawDescGZIP() []byte {
	file_proto_chat_chat_proto_rawDescOnce.Do(func() {
		file_proto_chat_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_chat_chat_proto_rawDescData)
	})
	return file_proto_chat_chat_proto_rawDescData
}

var file_proto_chat_chat_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_chat_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_chat_chat_proto_goTypes = []any{
	(EncryptionType)(0),     // 0: chat.EncryptionType
	(*MessageRequest)(nil),  // 1: chat.MessageRequest
	(*MessageResponse)(nil), // 2: chat.MessageResponse
}
var file_proto_chat_chat_proto_depIdxs = []int32{
	0, // 0: chat.MessageRequest.encryption_type:type_name -> chat.EncryptionType
	0, // 1: chat.MessageResponse.encryption_type:type_name -> chat.EncryptionType
	1, // 2: chat.ChatService.StreamMessages:input_type -> chat.MessageRequest
	2, // 3: chat.ChatService.StreamMessages:output_type -> chat.MessageResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_chat_chat_proto_init() }
func file_proto_chat_chat_proto_init() {
	if File_proto_chat_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_chat_chat_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chat_chat_proto_goTypes,
		DependencyIndexes: file_proto_chat_chat_proto_depIdxs,
		EnumInfos:         file_proto_chat_chat_proto_enumTypes,
		MessageInfos:      file_proto_chat_chat_proto_msgTypes,
	}.Build()
	File_proto_chat_chat_proto = out.File
	file_proto_chat_chat_proto_rawDesc = nil
	file_proto_chat_chat_proto_goTypes = nil
	file_proto_chat_chat_proto_depIdxs = nil
}

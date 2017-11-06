// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cockroach/pkg/storage/lease_status.proto

package storage

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import cockroach_roachpb2 "github.com/cockroachdb/cockroach/pkg/roachpb"
import cockroach_util_hlc "github.com/cockroachdb/cockroach/pkg/util/hlc"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type LeaseState int32

const (
	// ERROR indicates that the lease can't be used or acquired.
	LeaseState_ERROR LeaseState = 0
	// VALID indicates that the lease can be used.
	LeaseState_VALID LeaseState = 1
	// STASIS indicates that the lease has not expired, but can't be used.
	LeaseState_STASIS LeaseState = 2
	// EXPIRED indicates that the lease can't be used.
	LeaseState_EXPIRED LeaseState = 3
	// PROSCRIBED indicates that the lease's proposed timestamp is earlier than
	// allowed.
	LeaseState_PROSCRIBED LeaseState = 4
)

var LeaseState_name = map[int32]string{
	0: "ERROR",
	1: "VALID",
	2: "STASIS",
	3: "EXPIRED",
	4: "PROSCRIBED",
}
var LeaseState_value = map[string]int32{
	"ERROR":      0,
	"VALID":      1,
	"STASIS":     2,
	"EXPIRED":    3,
	"PROSCRIBED": 4,
}

func (x LeaseState) String() string {
	return proto.EnumName(LeaseState_name, int32(x))
}
func (LeaseState) EnumDescriptor() ([]byte, []int) { return fileDescriptorLeaseStatus, []int{0} }

// LeaseStatus holds the lease state, the timestamp at which the state
// is accurate, the lease and optionally the liveness if the lease is
// epoch-based.
type LeaseStatus struct {
	// Lease which this status describes.
	Lease cockroach_roachpb2.Lease `protobuf:"bytes,1,opt,name=lease" json:"lease"`
	// Timestamp that the lease was evaluated at.
	Timestamp cockroach_util_hlc.Timestamp `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp"`
	// State of the lease at timestamp.
	State LeaseState `protobuf:"varint,3,opt,name=state,proto3,enum=cockroach.storage.LeaseState" json:"state,omitempty"`
	// Liveness if this is an epoch-based lease.
	Liveness *Liveness `protobuf:"bytes,4,opt,name=liveness" json:"liveness,omitempty"`
}

func (m *LeaseStatus) Reset()                    { *m = LeaseStatus{} }
func (m *LeaseStatus) String() string            { return proto.CompactTextString(m) }
func (*LeaseStatus) ProtoMessage()               {}
func (*LeaseStatus) Descriptor() ([]byte, []int) { return fileDescriptorLeaseStatus, []int{0} }

func init() {
	proto.RegisterType((*LeaseStatus)(nil), "cockroach.storage.LeaseStatus")
	proto.RegisterEnum("cockroach.storage.LeaseState", LeaseState_name, LeaseState_value)
}
func (m *LeaseStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LeaseStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintLeaseStatus(dAtA, i, uint64(m.Lease.Size()))
	n1, err := m.Lease.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x12
	i++
	i = encodeVarintLeaseStatus(dAtA, i, uint64(m.Timestamp.Size()))
	n2, err := m.Timestamp.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if m.State != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintLeaseStatus(dAtA, i, uint64(m.State))
	}
	if m.Liveness != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintLeaseStatus(dAtA, i, uint64(m.Liveness.Size()))
		n3, err := m.Liveness.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func encodeVarintLeaseStatus(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *LeaseStatus) Size() (n int) {
	var l int
	_ = l
	l = m.Lease.Size()
	n += 1 + l + sovLeaseStatus(uint64(l))
	l = m.Timestamp.Size()
	n += 1 + l + sovLeaseStatus(uint64(l))
	if m.State != 0 {
		n += 1 + sovLeaseStatus(uint64(m.State))
	}
	if m.Liveness != nil {
		l = m.Liveness.Size()
		n += 1 + l + sovLeaseStatus(uint64(l))
	}
	return n
}

func sovLeaseStatus(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozLeaseStatus(x uint64) (n int) {
	return sovLeaseStatus(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LeaseStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLeaseStatus
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LeaseStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LeaseStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lease", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Lease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (LeaseState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liveness", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Liveness == nil {
				m.Liveness = &Liveness{}
			}
			if err := m.Liveness.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLeaseStatus(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipLeaseStatus(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLeaseStatus
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
					return 0, ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLeaseStatus
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthLeaseStatus
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowLeaseStatus
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipLeaseStatus(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthLeaseStatus = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLeaseStatus   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("cockroach/pkg/storage/lease_status.proto", fileDescriptorLeaseStatus) }

var fileDescriptorLeaseStatus = []byte{
	// 349 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x40, 0x33, 0xfd, 0xfd, 0x7a, 0x0b, 0x25, 0xdf, 0xe0, 0x22, 0x54, 0x8c, 0x55, 0x44, 0x8a,
	0x8b, 0x04, 0xac, 0xe0, 0xba, 0xb5, 0x59, 0x04, 0x0b, 0x2d, 0x93, 0x22, 0xe2, 0x46, 0xa6, 0x71,
	0x48, 0x4b, 0x53, 0x27, 0x34, 0x53, 0x9f, 0xc3, 0xc7, 0xea, 0xd2, 0xa5, 0x2b, 0xd1, 0xf8, 0x0a,
	0x3e, 0x80, 0x64, 0x32, 0x4d, 0x2d, 0xe8, 0xee, 0x12, 0xce, 0xb9, 0x73, 0x72, 0xa1, 0xed, 0x73,
	0x7f, 0xbe, 0xe4, 0xd4, 0x9f, 0xda, 0xd1, 0x3c, 0xb0, 0x63, 0xc1, 0x97, 0x34, 0x60, 0x76, 0xc8,
	0x68, 0xcc, 0xee, 0x63, 0x41, 0xc5, 0x2a, 0xb6, 0xa2, 0x25, 0x17, 0x1c, 0xff, 0xcf, 0x49, 0x4b,
	0x51, 0xcd, 0xd6, 0xae, 0x2c, 0xa7, 0x68, 0x62, 0x3f, 0x50, 0x41, 0x33, 0xa9, 0x79, 0xf2, 0xc7,
	0xfa, 0xd9, 0x13, 0x7b, 0x64, 0xb1, 0x5a, 0xdd, 0x3c, 0xdd, 0xa5, 0x56, 0x62, 0x16, 0xda, 0xd3,
	0xd0, 0xb7, 0xc5, 0x6c, 0xc1, 0x62, 0x41, 0x17, 0x91, 0xe2, 0xf6, 0x02, 0x1e, 0x70, 0x39, 0xda,
	0xe9, 0x94, 0x7d, 0x3d, 0xfe, 0x42, 0x50, 0x1f, 0xa4, 0xbd, 0x9e, 0xcc, 0xc5, 0x17, 0x50, 0x96,
	0xf9, 0x06, 0x6a, 0xa1, 0x76, 0xfd, 0xdc, 0xb0, 0xb6, 0xe1, 0xaa, 0xd0, 0x92, 0x78, 0xaf, 0xb4,
	0x7e, 0x3b, 0xd4, 0x48, 0x06, 0xe3, 0x2e, 0xd4, 0xf2, 0xe7, 0x8c, 0x82, 0x34, 0x0f, 0x7e, 0x98,
	0x69, 0x93, 0x35, 0x0d, 0x7d, 0x6b, 0xbc, 0x81, 0x94, 0xbe, 0xb5, 0x70, 0x07, 0xca, 0xe9, 0xc5,
	0x98, 0x51, 0x6c, 0xa1, 0x76, 0x63, 0x47, 0x57, 0x3f, 0x6e, 0xe5, 0x9d, 0x8c, 0x64, 0x2c, 0xbe,
	0x84, 0x7f, 0x9b, 0x6b, 0x18, 0x25, 0xf9, 0xec, 0xfe, 0x6f, 0x9e, 0x42, 0x48, 0x0e, 0x9f, 0x5d,
	0x03, 0x6c, 0xb7, 0xe1, 0x1a, 0x94, 0x1d, 0x42, 0x86, 0x44, 0xd7, 0xd2, 0xf1, 0xa6, 0x3b, 0x70,
	0xfb, 0x3a, 0xc2, 0x00, 0x15, 0x6f, 0xdc, 0xf5, 0x5c, 0x4f, 0x2f, 0xe0, 0x3a, 0x54, 0x9d, 0xdb,
	0x91, 0x4b, 0x9c, 0xbe, 0x5e, 0xc4, 0x0d, 0x80, 0x11, 0x19, 0x7a, 0x57, 0xc4, 0xed, 0x39, 0x7d,
	0xbd, 0xd4, 0x3b, 0x5a, 0x7f, 0x98, 0xda, 0x3a, 0x31, 0xd1, 0x4b, 0x62, 0xa2, 0xd7, 0xc4, 0x44,
	0xef, 0x89, 0x89, 0x9e, 0x3f, 0x4d, 0xed, 0xae, 0xaa, 0x12, 0x26, 0x15, 0x79, 0xed, 0xce, 0x77,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x14, 0x47, 0x83, 0x32, 0x02, 0x00, 0x00,
}
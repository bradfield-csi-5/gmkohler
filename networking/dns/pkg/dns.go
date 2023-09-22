package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
// The header contains the following fields:
//
//	0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
//	f  e  d  c  b  a  9  8  7  6  5  4  3  2  1  0
//
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                      ID                       |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    QDCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    ANCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    NSCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    ARCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type DnsHeader struct {
	Identification        uint16
	Flags                 DnsFlags
	QuestionCount         uint16
	AnswerCount           uint16
	NameServerCount       uint16
	AdditionalRecordCount uint16
}
type DnsMessage struct {
	Header     DnsHeader
	Questions  []ResourceRecord
	Answers    []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

func Query(rName ResourceName, rType ResourceType) *DnsMessage {
	flags := NewDnsFlagBuilder().
		SetOpCodeQuery().
		SetRecursionDesired().
		Build()
	return &DnsMessage{
		Header: DnsHeader{
			Identification: uint16(rand.Intn(0xffff)),
			Flags:          flags,
			QuestionCount:  1,
		},
		Questions: []ResourceRecord{*QuestionRecord(rName, rType)},
	}
}
func (dm *DnsMessage) Encode() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, dm.Header)
	if err != nil {
		return nil, fmt.Errorf("error writing DnsMessage.Header: %w", err)
	}
	for _, q := range dm.Questions {
		encodedQuestion, err := q.Encode()
		if err != nil {
			return nil, err
		}
		err = binary.Write(&buf, binary.BigEndian, encodedQuestion)
		if err != nil {
			return nil, fmt.Errorf(
				"DnsMessage.Encode(): error writing Question: %w",
				err,
			)
		}
	}
	for _, ans := range dm.Answers {
		encodedAnswer, err := ans.Encode()
		if err != nil {
			return nil, err
		}
		err = binary.Write(&buf, binary.BigEndian, encodedAnswer)
		if err != nil {
			return nil, fmt.Errorf(
				"DnsMessage.Encode(): error writing Answer: %w",
				err,
			)
		}
	}
	for _, auth := range dm.Authority {
		encodedAuthority, err := auth.Encode()
		if err != nil {
			return nil, err
		}
		err = binary.Write(&buf, binary.BigEndian, encodedAuthority)
		if err != nil {
			return nil, fmt.Errorf(
				"DnsMessage.Encode(): error writing Authority: %w",
				err,
			)
		}
	}
	for _, add := range dm.Additional {
		encodedAdditional, err := add.Encode()
		if err != nil {
			return nil, err
		}
		err = binary.Write(&buf, binary.BigEndian, encodedAdditional)
		if err != nil {
			return nil, fmt.Errorf(
				"DnsMessage.Encode(): error writing Authority: %w",
				err,
			)
		}
	}
	return buf.Bytes(), nil
}
func DecodeMessage(encoded []byte) (*DnsMessage, error) {
	r := bytes.NewReader(encoded)
	header := new(DnsHeader)
	err := binary.Read(r, binary.BigEndian, header)
	if err != nil {
		return nil, err
	}
	fmt.Printf("decoded header: %+v\n", *header)
	var questions []ResourceRecord
	for j := uint16(0); j < header.QuestionCount; j++ {
		rrp, err := decodeCommonRrComponents(r)
		if err != nil {
			return nil, err
		}
		questions = append(questions, *rrp)
	}
	var answers []ResourceRecord
	for j := uint16(0); j < header.AnswerCount; j++ {
		rrp, err := decodeResourceRecord(r)
		if err != nil {
			return nil, fmt.Errorf(
				"DecodeMessage(): error decoding answer: %w",
				err,
			)
		}
		answers = append(answers, *rrp)
	}
	var authorityRecords []ResourceRecord
	for j := uint16(0); j < header.NameServerCount; j++ {
		rrp, err := decodeResourceRecord(r)
		if err != nil {
			return nil, fmt.Errorf(
				"DecodeMessage(): error decoding authority record: %w",
				err,
			)
		}
		authorityRecords = append(authorityRecords, *rrp)
	}
	var additionalRecords []ResourceRecord
	for j := uint16(0); j < header.AdditionalRecordCount; j++ {
		rrp, err := decodeResourceRecord(r)
		if err != nil {
			return nil, err
		}
		additionalRecords = append(additionalRecords, *rrp)
	}
	return &DnsMessage{
		Header:     *header,
		Questions:  questions,
		Answers:    answers,
		Authority:  authorityRecords,
		Additional: additionalRecords,
	}, nil
}

func decodeCommonRrComponents(r *bytes.Reader) (*ResourceRecord, error) {
	name, err := decodeResourceName(r)
	if err != nil {
		return nil, fmt.Errorf(
			"decodeCommonRrComponents(): error decoding resource name: %w",
			err,
		)
	}
	var rType ResourceType
	err = binary.Read(r, binary.BigEndian, &rType)
	if err != nil {
		return nil, fmt.Errorf(
			"decodeCommonRrComponents(): error decoding resource type: %w",
			err,
		)
	}
	fmt.Printf("resource type: %v\n", rType)
	var rClass ResourceClass
	err = binary.Read(r, binary.BigEndian, &rClass)
	if err != nil {
		return nil, fmt.Errorf(
			"decodeCommonRrComponents(): error decoding resource class: %w",
			err,
		)
	}
	fmt.Printf("class: %+v\n", rClass)
	return &ResourceRecord{
		Name:  name,
		Type:  rType,
		Class: rClass,
	}, nil
}

func decodeResourceRecord(r *bytes.Reader) (*ResourceRecord, error) {
	rrp, err := decodeCommonRrComponents(r)
	if err != nil {
		return nil, fmt.Errorf("decodeResourceRecord(): %w", err)
	}
	var ttl uint32
	err = binary.Read(r, binary.BigEndian, &ttl)
	if err != nil {
		return nil, fmt.Errorf(
			"decodeResourceRecord(): error decoding ttl: %w",
			err,
		)
	}
	fmt.Printf("ttl: %d s\n", ttl)
	var rDataLength uint16
	err = binary.Read(r, binary.BigEndian, &rDataLength)
	if err != nil {
		return nil, fmt.Errorf(
			"decodeResourceRecord(): error decoding data length: %w",
			err,
		)
	}
	fmt.Printf("data length: %d\n", rDataLength)

	rrp.Ttl = ttl
	rrp.RecordDataLength = rDataLength
	rrp.recordData = decodeResourceData(r, rDataLength, rrp.Type)
	return rrp, nil
}

func decodeResourceData(
	r *bytes.Reader,
	rDataLength uint16,
	rt ResourceType,
) string {
	switch rt {
	case CName, Address: // interpret data as an IP address
		rData := make([]byte, rDataLength)
		if err := binary.Read(r, binary.BigEndian, rData); err != nil {
			panic(fmt.Sprintf(
				"decodeResourceRecord(): error decoding data: %v",
				err,
			))
		}
		strs := make([]string, len(rData))
		for j, b := range rData {
			strs[j] = strconv.Itoa(int(b))
		}
		return strings.Join(strs, ".")
	case NameServer: // interpret data as a domain name
		name, err := decodeResourceName(r)
		if err != nil {
			panic("error decoding name")
		}
		return string(name)
	case Soa:
		primaryNameServer, err := decodeResourceName(r)
		if err != nil {
			panic(fmt.Sprintf(
				"error decoding MNAME of SOA data: %v",
				err,
			))
		}
		responsibleAuthorityMailbox, err := decodeResourceName(r)
		if err != nil {
			panic(fmt.Sprintf(
				"error decoding RNAME of SOA data: %v",
				err,
			))
		}
		var serial uint32
		if err := binary.Read(r, binary.BigEndian, &serial); err != nil {
			panic(fmt.Sprintf(
				"error decoding SERIAL of SOA data: %v",
				err,
			))
		}
		var refresh uint32
		if err := binary.Read(r, binary.BigEndian, &refresh); err != nil {
			panic(fmt.Sprintf(
				"error decoding REFRESH of SOA data: %v",
				err,
			))
		}
		var retry uint32
		if err := binary.Read(r, binary.BigEndian, &retry); err != nil {
			panic(fmt.Sprintf(
				"error decoding RETRY of SOA data: %v",
				err,
			))
		}
		var expire uint32
		if err := binary.Read(r, binary.BigEndian, &expire); err != nil {
			panic(fmt.Sprintf(
				"error decoding EXPIRE of SOA data: %v",
				err,
			))
		}
		var minimum uint32
		if err := binary.Read(r, binary.BigEndian, &minimum); err != nil {
			panic(fmt.Sprintf(
				"error decoding MINIMUM of SOA data: %v",
				err,
			))
		}
		return fmt.Sprintf("Primary name server: %s, "+
			"Responsible authority's mailbox: %s, "+
			"Serial Number: %d, "+
			"Refresh Interval: %ds, "+
			"Retry Interval: %ds, "+
			"Expire limit: %ds, "+
			"Minimum TTL: %ds",
			primaryNameServer,
			responsibleAuthorityMailbox,
			serial,
			refresh,
			retry,
			expire,
			minimum,
		)
	default:
		rData := make([]byte, rDataLength)
		if _, err := r.Read(rData); err != nil {
			panic(fmt.Sprintf(
				"error reading data: %v",
				err,
			))
		}
		return string(rData)
	}
}

type ResourceName string

func (rn *ResourceName) Encode() ([]byte, error) {
	split := strings.Split(strings.ToLower(string(*rn)), ".")
	var encoded []byte
	for _, s := range split {
		b := []byte(s)
		for _, c := range b {
			if c > unicode.MaxASCII {
				return nil, fmt.Errorf("byte is not ASCII-compliant: %02x", c)
			}
		}
		encoded = append(encoded, byte(len(b)))
		encoded = append(encoded, b...)
	}
	encoded = append(encoded, 0x0)
	return encoded, nil
}
func decodeResourceName(r *bytes.Reader) (ResourceName, error) {
	var decoded []string
	for {
		l, err := r.ReadByte()
		if err != nil {
			return "", fmt.Errorf("error reading length byte: %w", err)
		}
		if l == 0x0 {
			break
		}
		if l>>6 == 0x03 {
			fmt.Printf("Found pointer-based name compression\n")
			restOfPointer, err := r.ReadByte()
			if err != nil {
				return "", err
			}
			p := uint16(l&0x3f)<<8 | uint16(restOfPointer)
			offset, _ := r.Seek(0, io.SeekCurrent)
			_, err = r.Seek(int64(p), io.SeekStart)
			if err != nil {
				return "", err
			}
			rName, err := decodeResourceName(r)
			if err != nil {
				return "", err
			}
			_, err = r.Seek(offset, io.SeekStart)
			if err != nil {
				return "", err
			}
			decoded = append(decoded, string(rName))
			return ResourceName(strings.Join(decoded, ".")), nil
		}
		s := make([]byte, l)
		_, err = r.Read(s)
		if err != nil {
			return "", fmt.Errorf("error reading string: %w", err)
		}
		decoded = append(decoded, string(s))
	}
	if len(decoded) > 0 {
		return ResourceName(strings.Join(decoded, ".")), nil
	} else {
		return "", nil
	}
}

type DnsFlagBuilder struct {
	f DnsFlags
}

func NewDnsFlagBuilder() *DnsFlagBuilder {
	return &DnsFlagBuilder{}
}
func (fb *DnsFlagBuilder) SetOpCodeQuery() *DnsFlagBuilder {
	fb.f &= ^(opCodeMask)                    // clear opcode
	fb.f |= DnsFlags(0b0000 << opCodeOffset) // set opcode to zero
	return fb
}
func (fb *DnsFlagBuilder) SetRecursionDesired() *DnsFlagBuilder {
	fb.f |= 0x100
	return fb
}
func (fb *DnsFlagBuilder) Build() DnsFlags {
	return fb.f
}

type DnsFlags uint16

func (df *DnsFlags) opCode() byte {
	return byte((*df & opCodeMask) >> opCodeOffset)
}

func (df *DnsFlags) Code() ResponseCode {
	return ResponseCode(*df & responseCodeMask)
}

type ResponseCode uint8

func (r ResponseCode) String() string {
	switch r {
	case Ok:
		return "OK"
	case FormatError:
		return "Format error"
	case ServerFailure:
		return "Server failure"
	case NameError:
		return "Name Error"
	case NotImplemented:
		return "Not Implemented"
	case Refused:
		return "Refused"
	default:
		return r.String()
	}
}

const (
	Ok ResponseCode = iota
	FormatError
	ServerFailure
	NameError
	NotImplemented
	Refused
)

// TYPE values: https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
type ResourceType uint16

const (
	Address    ResourceType = 0x01 // 'A'
	NameServer ResourceType = 0x02 // 'NS'
	CName      ResourceType = 0x05 // 'CNAME'
	Soa        ResourceType = 0x06
	MailServer ResourceType = 0x15 // 'MX'
)

func ResourceTypeFromString(s string) ResourceType {
	switch s {
	case "CNAME":
		return CName
	case "A":
		return Address
	case "NS":
		return NameServer
	case "MX":
		return MailServer
	default:
		panic(fmt.Sprintf("unrecognized ResourceType %s", s))
	}
}
func (rt ResourceType) String() string {
	switch rt {
	case CName:
		return "CNAME"
	case Address:
		return "A"
	case NameServer:
		return "NS"
	case Soa:
		return "SOA"
	case MailServer:
		return "MX"
	default:
		panic(fmt.Sprintf(
			"ResourceType.String(): unrecognized resource type %d",
			rt,
		))
	}
}

type ResourceClass uint16

func (rc ResourceClass) String() string {
	switch rc {
	case Internet:
		return "IN"
	case CsNet:
		return "CS"
	case Chaos:
		return "CH"
	case Hesiod:
		return "HS"
	default:
		panic(fmt.Sprintf(
			"ResourceClass.String(): unrecognized resource class %d",
			rc,
		))
	}

}

const (
	Internet ResourceClass = iota + 1
	CsNet
	Chaos
	Hesiod
)

const (
	opCodeOffset     DnsFlags = 0xa
	opCodeMask       DnsFlags = 0x78_00
	responseCodeMask DnsFlags = 0x0f
)

func (rr *ResourceRecord) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encodedName, err := rr.Name.Encode()

	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error encoding Name: %w",
			err,
		)
	}
	fmt.Printf("encoded name: %v\n", encodedName)
	err = binary.Write(&buf, binary.BigEndian, encodedName)
	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error writing Name: %w",
			err,
		)
	}
	err = binary.Write(&buf, binary.BigEndian, rr.Type)
	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error writing Type: %w",
			err,
		)
	}
	err = binary.Write(&buf, binary.BigEndian, rr.Class)
	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error writing Class: %w",
			err,
		)
	}
	err = binary.Write(&buf, binary.BigEndian, rr.Ttl)
	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error writing Ttl: %w",
			err,
		)
	}
	err = binary.Write(&buf, binary.BigEndian, rr.RecordDataLength)
	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error writing RecordDataLength: %w",
			err,
		)
	}
	err = binary.Write(&buf, binary.BigEndian, []byte(rr.recordData))
	if err != nil {
		return nil, fmt.Errorf(
			"ResourceRecord.Encode(): error writing recordData: %w",
			err,
		)
	}
	return buf.Bytes(), nil
}

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.3
type ResourceRecord struct {
	Name             ResourceName
	Type             ResourceType // use RECORD_* constants
	Class            ResourceClass
	Ttl              uint32
	RecordDataLength uint16
	recordData       string
}

func QuestionRecord(rName ResourceName, rType ResourceType) *ResourceRecord {
	return &ResourceRecord{
		Name:  rName,
		Type:  rType,
		Class: Internet,
	}
}

package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"tracio.com/sitemanagerservice/dwg/file/version"
	"tracio.com/sitemanagerservice/dwg/types"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type IDwgStreamReader interface {
	GetPosition() (int64, error)
	SetPosition(pos int64) (int64, error)
	GetStream() io.ReadSeeker
	IsEmpty() bool
	SetIsEmpty(isEmpty bool)

	ReadByte() (byte, error)
	ReadShort() (int16, error)
	ReadUShort() (uint16, error)
	ReadInt() (int32, error)
	ReadUInt() (uint32, error)
	ReadDouble() (float64, error)
	ReadBytes(length int) ([]byte, error)

	ReadBit() (bool, error)
	ReadBitAsShort() (int16, error)
	Read2Bits() (byte, error)
	ReadBitShort() (int16, error)
	ReadBitShortAsBool() (bool, error)
	ReadBitLong() (int32, error)
	ReadBitLongLong() (int32, error)
	ReadBitDouble() (float64, error)

	ReadRawChar() (int8, error)
	ReadRawLong() (int32, error)
	ReadModularChar() (uint32, error)
	ReadSignedModularChar() (int32, error)

	ReadVariableText() (string, error)
	ReadSentinel() ([]byte, error)

	ReadDateTime() (time.Time, error)
	ReadTimeSpan() (time.Duration, error)

	GetPositionInBits() (int32, error)
	SetPositionInBits(position int32) error

	SetPositionByFlag(position int32) (int32, error)

	Advance(offset int64) error
	ResetShift() (uint16, error)
}

type DwgStreamReader struct {
	IDwgStreamReader

	stream   io.ReadSeeker
	bitShift int
	lastByte byte
	endian   binary.ByteOrder
	isEmpty  bool
}

func (reader *DwgStreamReader) GetStream() io.ReadSeeker {
	return reader.stream
}

func (reader *DwgStreamReader) GetPosition() (int64, error) {
	return reader.stream.Seek(0, io.SeekCurrent)
}

func (reader *DwgStreamReader) SetPosition(pos int64) (int64, error) {
	err := error(nil)
	if pos, err = reader.stream.Seek(pos, io.SeekStart); err == nil {
		reader.bitShift = 0
	}
	return pos, err
}

func (reader *DwgStreamReader) IsEmpty() bool {
	return reader.isEmpty
}

func (reader *DwgStreamReader) SetIsEmpty(isEmpty bool) {
	reader.isEmpty = isEmpty
}

func (reader *DwgStreamReader) SetPositionByFlag(position int32) (int32, error) {
	var err error

	if err = reader.SetPositionInBits(position); err != nil {
		return 0, err
	}

	//String stream present bit (last bit in pre-handles section).
	var flag bool
	if flag, err = reader.ReadBit(); err != nil {
		return 0, err
	}

	startPositon := position
	if flag {
		//String stream present

		//If 1, then the “endbit” location should be decremented by 16 bytes
		var length, size int32
		length, size, err = reader.applyFlagToPosition(position)

		startPositon = length - size

		reader.SetPositionInBits(startPositon)
	} else {
		//Mark as empty
		reader.SetIsEmpty(true)
		//There is no information, set the position to the end
		reader.GetStream().Seek(0, io.SeekEnd)
	}

	return startPositon, nil
}

func (reader *DwgStreamReader) ReadBit() (value bool, err error) {
	if reader.bitShift == 0 {
		if err = reader.AdvanceByte(); err != nil {
			return false, err
		}
		value = (reader.lastByte & 128) != 0
		reader.bitShift = 1
	} else {
		value = ((reader.lastByte << reader.bitShift) & 128) != 0
		reader.bitShift = (reader.bitShift + 1) & 7
	}

	return value, nil
}

func (reader *DwgStreamReader) ReadBitAsShort() (value int16, err error) {
	var bValue bool
	if bValue, err = reader.ReadBit(); bValue {
		value = 1
	} else {
		value = 0
	}

	return value, err
}

func (reader *DwgStreamReader) Read2Bits() (value byte, err error) {
	if reader.bitShift == 0 {
		if err = reader.AdvanceByte(); err != nil {
			return value, err
		}
		value = byte(reader.lastByte >> 6)
		reader.bitShift = 2
	} else if reader.bitShift == 7 {
		lastValue := (reader.lastByte << 1) & 2
		if err = reader.AdvanceByte(); err != nil {
			return value, err
		}
		value = byte(lastValue | (reader.lastByte >> 7))
		reader.bitShift = 1
	} else {
		value = byte((int(reader.lastByte) >> (6 - reader.bitShift)) & 3)
		reader.bitShift = (reader.bitShift + 2) & 7
	}

	return value, nil
}

func (reader *DwgStreamReader) ReadBitShort() (value int16, err error) {
	var bValue byte
	if bValue, err = reader.Read2Bits(); err != nil {
		return 0, err
	}
	switch bValue {
	case 0:
		return reader.ReadShort()
	case 1:
		if reader.bitShift == 0 {
			if err = reader.AdvanceByte(); err != nil {
				return 0, err
			}
			value = int16(reader.lastByte)
		} else {
			bValue, err = reader.applyShiftToLasByte()
			value = int16(bValue)
		}
	case 2:
		value = 0
	case 3:
		value = 256
	}

	return value, err
}

func (reader *DwgStreamReader) ReadBitShortAsBool() (bool, error) {
	shortValue, err := reader.ReadBitShort()
	return shortValue != 0, err
}

func (reader *DwgStreamReader) ReadBitLong() (int32, error) {
	var value int32
	bValue, err := reader.Read2Bits()
	switch bValue {
	case 0:
		value, err = reader.ReadInt()
		break
	case 1:
		if reader.bitShift == 0 {
			err = reader.AdvanceByte()
			value = int32(reader.lastByte)
		} else {
			bValue, err = reader.applyShiftToLasByte()
			value = int32(bValue)
		}
		break
	case 2:
		value = 0.0
		break
	default:
		err = errors.New("Invalid ReadBitLong value")
	}

	return value, err

}

func (reader *DwgStreamReader) ReadBitLongLong() (value int32, err error) {
	var size byte
	size, err = reader.read3bits()

	for i := byte(0); i < size; i++ {
		var bValue byte
		bValue, err = reader.ReadByte()
		value += int32(bValue) << (i << 3)
	}

	return value, err
}

func (reader *DwgStreamReader) ReadBitDouble() (float64, error) {
	var value float64
	bValue, err := reader.Read2Bits()
	switch bValue {
	case 0:
		value, err = reader.ReadDouble()
		break
	case 1:
		value = 1.0
		break
	case 2:
		value = 0.0
		break
	default:
		err = errors.New("Invalid BitDouble value")
	}

	return value, err
}

func (reader *DwgStreamReader) ReadByteBase() (value byte, err error) {
	arr := make([]byte, 1)
	if _, err = reader.stream.Read(arr); err == nil {
		value = arr[0]
	}
	return value, err
}

func (reader *DwgStreamReader) ReadByte() (value byte, err error) {
	if reader.bitShift == 0 {
		reader.lastByte, err = reader.ReadByteBase()

		return reader.lastByte, err
	}

	value = uint8(reader.lastByte) << reader.bitShift
	reader.lastByte, err = reader.ReadByteBase()
	value |= reader.lastByte >> (8 - reader.bitShift)

	return value, err
}

func (reader *DwgStreamReader) ReadBytes(length int) ([]byte, error) {
	numArray := make([]byte, length)
	err := reader.applyShiftToArr(numArray)
	return numArray, err
}

func (reader *DwgStreamReader) ReadShort() (value int16, err error) {
	var data []byte
	if data, err = reader.ReadBytes(2); err != nil {
		return 0, err
	}

	buf := bytes.NewReader(data)
	if err = binary.Read(buf, reader.endian, &value); err != nil {
		return 0, err
	}

	return value, err
}

func (reader *DwgStreamReader) ReadUShort() (value uint16, err error) {
	var data []byte
	if data, err = reader.ReadBytes(2); err != nil {
		return 0, err
	}

	buf := bytes.NewReader(data)
	if err = binary.Read(buf, reader.endian, &value); err != nil {
		return 0, err
	}

	return value, err
}

func (reader *DwgStreamReader) ReadInt() (value int32, err error) {
	var data []byte
	if data, err = reader.ReadBytes(4); err != nil {
		return 0, err
	}

	buf := bytes.NewReader(data)
	if err = binary.Read(buf, reader.endian, &value); err != nil {
		return 0, err
	}

	return value, err
}

func (reader *DwgStreamReader) ReadRawChar() (value int8, err error) {
	var bValue byte
	if bValue, err = reader.ReadByte(); err != nil {
		value = int8(bValue)
	}
	return value, err
}

func (reader *DwgStreamReader) ReadRawLong() (value int32, err error) {
	return reader.ReadInt()
}

// func (reader *DwgStreamReader) ReadUInt() (value uint32, err error) {
// 	err = binary.Read(reader.stream, binary.LittleEndian, &value)
// 	return value, err
// }

func (reader *DwgStreamReader) ReadDouble() (value float64, err error) {
	var data []byte
	if data, err = reader.ReadBytes(8); err != nil {
		return 0, err
	}

	buf := bytes.NewReader(data)
	if err = binary.Read(buf, reader.endian, &value); err != nil {
		return 0, err
	}

	return value, err
}

func (reader *DwgStreamReader) read3bits() (value byte, err error) {
	value = 0

	var bValue bool
	if bValue, err = reader.ReadBit(); bValue {
		value = 1
	}
	value = value << 1
	if bValue, err = reader.ReadBit(); bValue {
		value |= 1
	}
	value = value << 1
	if bValue, err = reader.ReadBit(); bValue {
		value |= 1
	}

	return value, err
}

func (reader *DwgStreamReader) julianToDate(jdate int32, milliseconds int32) (value time.Time, err error) {
	unixTime := int64((float64(jdate) - 2440587.5) * float64(86400))
	value = time.Unix(unixTime, 0)
	value = value.Add(time.Duration(milliseconds) * time.Duration(time.Millisecond))

	return value, nil
}

func (reader *DwgStreamReader) ReadDateTime() (value time.Time, err error) {
	var jdate, milliseconds int32
	jdate, err = reader.ReadBitLong()
	milliseconds, err = reader.ReadBitLong()
	if err != nil {
		return value, err
	}

	return reader.julianToDate(jdate, milliseconds)
}

func (reader *DwgStreamReader) ReadTimeSpan() (value time.Duration, err error) {
	value = time.Duration(0)
	var hours, milliseconds int32
	hours, err = reader.ReadBitLong()
	milliseconds, err = reader.ReadBitLong()
	if err != nil {
		return value, err
	} else if hours < 0 || milliseconds < 0 {
		return value, nil
	}

	return time.Hour*time.Duration(hours) + time.Millisecond*time.Duration(milliseconds), nil
}

func (reader *DwgStreamReader) ReadModularChar() (value uint32, err error) {
	var lastByte byte
	if lastByte, err = reader.ReadByte(); err != nil {
		return value, err
	}

	// Remove the flag
	value = uint32(lastByte & 0b01111111)

	shift := 0
	for (lastByte & 0b10000000) != 0 {
		shift += 7
		if lastByte, err = reader.ReadByte(); err != nil {
			return value, err
		}
		value |= uint32(lastByte&0b01111111) << shift
	}

	return value, nil
}

// Modular characters are a method of storing compressed integer values. They are used in the object map to
// indicate both handle offsets and file location offsets.They consist of a stream of bytes, terminating when
// the high bit of the byte is 0.
func (reader *DwgStreamReader) ReadSignedModularChar() (value int32, err error) {
	if reader.bitShift == 0 {
		// No shift, read normal
		reader.AdvanceByte()

		//Check if the current byte
		if (reader.lastByte & 0b10000000) == 0 {
			//Drop the flags
			value = int32(reader.lastByte) & 0b00111111

			//Check the sign flag
			if (reader.lastByte & 0b01000000) > 0 {
				value = -value
			}
		} else {
			totalShift := 0
			sum := int32(reader.lastByte & 127)
			for true {
				//Shift to apply
				totalShift += 7
				reader.AdvanceByte()

				//Check if the highest byte is 0
				if (reader.lastByte & 0b10000000) != 0 {
					sum |= int32(reader.lastByte&127) << totalShift
				} else {
					break
				}
			}

			//Drop the flags at the las byte, and add it's value
			value = sum | (int32(reader.lastByte&0b00111111) << totalShift)

			//Check the sign flag
			if (reader.lastByte & 0b01000000) > 0 {
				value = -value
			}
		}
	} else {
		// Apply the shift to each byte
		var lastByte byte
		lastByte, _ = reader.applyShiftToLasByte()
		if (lastByte & 0b10000000) == 0 {
			// Drop the flags
			value = int32(lastByte) & 0b00111111

			// Check the sign flag
			if (lastByte & 0b01000000) > 0 {
				value = -value
			}
		} else {
			totalShift := 0
			sum := int32(lastByte) & 127
			var currByte byte
			for true {
				// Shift to apply
				totalShift += 7
				currByte, err = reader.applyShiftToLasByte()

				// Check if the highest byte is 0
				if (currByte & 0b10000000) != 0 {
					sum |= int32(currByte&127) << totalShift
				} else {
					break
				}
			}

			//Drop the flags at the las byte, and add it's value
			value = sum | (int32(currByte&0b00111111) << totalShift)

			//Check the sign flag
			if (currByte & 0b01000000) > 0 {
				value = -value
			}
		}
	}
	return value, nil
}

func (reader *DwgStreamReader) ReadString(length int, decoder *encoding.Decoder) (value string, err error) {
	if length == 0 {
		return "", nil
	}

	var data []byte
	if data, err = reader.ReadBytes(length); err != nil {
		return "", err
	}

	if decoder != nil {
		textReader := transform.NewReader(bytes.NewReader(data), decoder)
		data, err = ioutil.ReadAll(textReader)
	}

	return string(data), err
}

func (reader *DwgStreamReader) ReadVariableText() (value string, err error) {
	value = ""

	var length int16
	length, err = reader.ReadBitShort()
	if err == nil && length > 0 {
		value, err = reader.ReadString(int(length), nil)
		value = strings.Split(value, "\x00")[0]
	}

	return value, err
}

func (reader *DwgStreamReader) ReadSentinel() ([]byte, error) {
	return reader.ReadBytes(16)
}

func (reader *DwgStreamReader) GetPositionInBits() (int32, error) {
	position, err := reader.GetPosition()
	bitPosition := int32(position) * 8
	if reader.bitShift > 0 {
		bitPosition += int32(reader.bitShift) - 8
	}

	return bitPosition, err
}

func (reader *DwgStreamReader) SetPositionInBits(position int32) (err error) {
	if _, err = reader.SetPosition(int64(position) >> 3); err != nil {
		return err
	}
	reader.bitShift = int(position & 7)

	if reader.bitShift > 0 {
		err = reader.AdvanceByte()
	}

	return err
}

func (reader *DwgStreamReader) AdvanceByte() (err error) {
	reader.lastByte, err = reader.ReadByteBase()
	return err
}

func (reader *DwgStreamReader) Advance(offset int64) error {
	err := error(nil)
	if offset > 1 {
		var position int64
		if position, err = reader.GetPosition(); err != nil {
			return err
		}
		if position, err = reader.SetPosition(position + offset - 1); err != nil {
			return err
		}
	}

	_, err = reader.ReadByte()

	return err
}

func (reader *DwgStreamReader) ResetShift() (uint16, error) {
	reader.bitShift = 0

	reader.AdvanceByte()
	num := uint16(reader.lastByte)
	reader.AdvanceByte()

	num |= (uint16(reader.lastByte) << 8)

	return num, nil
}

func (reader *DwgStreamReader) applyFlagToPosition(lastPos int32) (length int32, strDataSize int32, err error) {
	//If 1, then the “endbit” location should be decremented by 16 bytes

	length = lastPos - 16
	reader.SetPositionInBits(length)

	//short should be read at location endbit – 128 (bits)
	var ushortValue uint16
	ushortValue, err = reader.ReadUShort()
	strDataSize = int32(ushortValue)

	//If this short has the 0x8000 bit set,
	//then decrement endbit by an additional 16 bytes,
	//strip the 0x8000 bit off of strDataSize, and read
	//the short at this new location, calling it hiSize.
	if (strDataSize & 0x8000) <= 0 {
		return length, strDataSize, err
	}

	length -= 16

	reader.SetPositionInBits(length)

	strDataSize &= 0x7FFF

	ushortValue, err = reader.ReadUShort()
	hiSize := int32(ushortValue)
	//Then set strDataSize to (strDataSize | (hiSize << 15))
	strDataSize += (hiSize & 0xFFFF) << 15

	//All unicode strings in this object are located in the “string stream”,
	//and should be read from this stream, even though the location of the
	//TV type fields in the object descriptions list these fields in among
	//the normal object data.

	return length, strDataSize, err
}

func (reader *DwgStreamReader) applyShiftToLasByte() (byte, error) {
	err := error(nil)
	value := reader.lastByte << reader.bitShift

	err = reader.AdvanceByte()

	return value | (reader.lastByte >> (8 - reader.bitShift)), err
}

func (reader *DwgStreamReader) applyShiftToArr(arr []byte) (err error) {
	if _, err = reader.stream.Read(arr); err != nil {
		return err
	}

	if reader.bitShift > 0 {
		shift := 8 - reader.bitShift
		for i := 0; i < len(arr); i++ {
			lastByteValue := byte(uint32(reader.lastByte) << reader.bitShift)
			reader.lastByte = arr[i]
			value := lastByteValue | (reader.lastByte >> shift)
			arr[i] = value
		}
	}

	return nil
}

type DwgStreamReaderAC21 struct {
	DwgStreamReader
}

func (reader *DwgStreamReaderAC21) ReadVariableText() (value string, err error) {
	value = ""

	var length int16
	length, err = reader.ReadBitShort()
	if err == nil && length > 0 {
		win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

		value, err = reader.ReadString(int(length)*2, win16le.NewDecoder())
		value = strings.Split(value, "\x00")[0]
	}

	return value, err
}

type DwgStreamReaderAC24 struct {
	DwgStreamReaderAC21
}

func (reader *DwgStreamReaderAC24) ReadObjectType() (objectType types.ObjectType, err error) {
	objectType = types.Object_UNUSED

	var pair byte
	if pair, err = reader.Read2Bits(); err != nil {
		return objectType, err
	}

	var bValue byte
	var sValue int16
	switch pair {
	//Read the following byte
	case 0:
		bValue, err = reader.ReadByte()
		objectType = types.ObjectType(bValue)
		break
	//Read following byte and add 0x1f0.
	case 1:
		bValue, err = reader.ReadByte()
		objectType = types.ObjectType(int16(0x1F0) + int16(bValue))
		break
	//Read the following two bytes (raw short)
	case 2:
		sValue, err = reader.ReadShort()
		objectType = types.ObjectType(sValue)
		break
	//The value 3 should never occur, but interpret the same as 2 nevertheless.
	case 3:
		sValue, err = reader.ReadShort()
		objectType = types.ObjectType(sValue)
		break
	}

	return objectType, err
}

func NewDwgStreamHandler(dwgVersion version.ACadVersion, stream io.ReadSeeker) (IDwgStreamReader, error) {
	var reader IDwgStreamReader
	reader = &DwgStreamReader{stream: stream, endian: binary.LittleEndian}

	switch dwgVersion {
	case version.AC1021:
		reader = &DwgStreamReaderAC21{
			DwgStreamReader{stream: stream, endian: binary.LittleEndian},
		}
	case version.AC1024, version.AC1027, version.AC1032:
		reader = &DwgStreamReaderAC24{
			DwgStreamReaderAC21{
				DwgStreamReader{stream: stream, endian: binary.LittleEndian},
			},
		}
	default:
		reader = nil
	}

	return reader, nil
}

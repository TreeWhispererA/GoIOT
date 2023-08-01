package stream

import (
	"encoding/binary"
	"io"
)

type DwgLZ77AC18Decompressor struct {
}

func copy(count int, src io.Reader, dst io.Writer) (byte, error) {
	var b byte
	err := error(nil)
	i := 0
	for true {
		if err = binary.Read(src, binary.LittleEndian, &b); err != nil {
			break
		}
		if i < count {
			if err = binary.Write(dst, binary.LittleEndian, b); err != nil {
				break
			}
		} else {
			break
		}
		i += 1
	}

	return b, err
}

func literalCount(code int, reader io.Reader) (int, error) {
	err := error(nil)
	lowbits := int(code & 0b1111)
	//0x00 : Set the running total to 0x0F, and read the next byte. From this point on, a 0x00 byte adds 0xFF to the running total, and a non-zero byte adds that value to the running total and terminates the process. Add 3 to the final result.
	var lastByte byte
	if lowbits == 0 {
		for true {
			if err = binary.Read(reader, binary.LittleEndian, &lastByte); err != nil {
				return 0, err
			} else if lastByte != 0 {
				break
			}
			lowbits += 0xFF
		}
		lowbits += int(0x0F + lastByte)
	}

	return lowbits, nil
}

func readCompressedBytes(opcode1 int, validBits int, compressed io.Reader) (int, error) {
	err := error(nil)
	compressedBytes := opcode1 & validBits

	if compressedBytes == 0 {
		lastByte := byte(0)
		for true {
			if err = binary.Read(compressed, binary.LittleEndian, &lastByte); err != nil {
				break
			} else if lastByte != 0 {
				break
			} else {
				compressedBytes += 0xFF
			}
		}

		compressedBytes += int(lastByte) + validBits
	}

	return compressedBytes + 2, err
}

func twoByteOffset(offset *int, addedValue int, reader io.Reader) (int, error) {
	var data [2]byte
	if err := binary.Read(reader, binary.LittleEndian, &data); err != nil {
		return 0, err
	}

	*offset |= int(data[0]) >> 2
	*offset |= int(data[1]) << 6
	*offset += addedValue

	return int(data[0]), nil
}

func (this DwgLZ77AC18Decompressor) Decompress(reader io.Reader, dst io.ReadWriteSeeker) error {
	err := error(nil)

	var opcode1 int
	var b byte
	if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
		return err
	} else {
		opcode1 = int(b)
	}

	if (opcode1 & 0xF0) == 0 {
		var count int
		if count, err = literalCount(int(opcode1), reader); err != nil {
			return err
		} else if b, err = copy(count+3, reader, dst); err != nil {
			return err
		}
		opcode1 = int(b)
	}

	for opcode1 != 0x11 {
		//0x00 – 0x0F : Not used, because this would be mistaken for a Literal Length in some situations.

		//Offset backwards from the current location in the decompressed data stream, where the “compressed” bytes should be copied from.
		compOffset := 0
		//Number of “compressed” bytes that are to be copied to this location from a previous location in the uncompressed data stream.
		compressedBytes := 0

		if opcode1 < 0x10 || opcode1 >= 0x40 {
			compressedBytes = int(opcode1>>4) - 1
			//Read the next byte(call it opcode2):
			var opcode2 byte
			if err = binary.Read(reader, binary.LittleEndian, &opcode2); err != nil {
				return err
			}
			compOffset = ((opcode1 >> 2 & 3) | (int(opcode2) << 2)) + 1
		} else if opcode1 < 0x20 {
			//0x12 – 0x1F
			if compressedBytes, err = readCompressedBytes(opcode1, 0b0111, reader); err != nil {
				return err
			}
			compOffset = (opcode1 & 8) << 11
			if opcode1, err = twoByteOffset(&compOffset, 0x4000, reader); err != nil {
				return err
			}
		} else if opcode1 >= 0x20 {
			//0x20
			if compressedBytes, err = readCompressedBytes(opcode1, 0b00011111, reader); err != nil {
				return err
			}
			if opcode1, err = twoByteOffset(&compOffset, 1, reader); err != nil {
				return err
			}
		}

		var position int64
		if position, err = dst.Seek(0, io.SeekCurrent); err != nil {
			return err
		}
		for i := int64(compressedBytes) + position; position < i; position++ {
			dst.Seek(position-int64(compOffset), io.SeekStart)
			binary.Read(dst, binary.LittleEndian, &b)
			dst.Seek(position, io.SeekStart)
			binary.Write(dst, binary.LittleEndian, b)
		}
		//Number of uncompressed or literal bytes to be copied from the input stream, following the addition of the compressed bytes.
		litCount := opcode1 & 3
		//0x00 : litCount is read as the next Literal Length (see format below)
		if litCount == 0 {
			if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
				return err
			}
			opcode1 = int(b)
			if (opcode1 & 0b11110000) == 0 {
				if litCount, err = literalCount(opcode1, reader); err != nil {
					return err
				}
				litCount += 3
			}
		}

		//Copy as literal
		if litCount > 0 {
			if b, err = copy(litCount, reader, dst); err != nil {
				return err
			}
			opcode1 = int(b)
		}
	}

	return nil
}

func DecompressLZ77AC18(reader io.Reader, decompressedSize int) (*MemoryStream, error) {
	buf := make([]byte, decompressedSize)
	dst := NewMemoryStream(buf)

	decompressor := DwgLZ77AC18Decompressor{}
	if err := decompressor.Decompress(reader, dst); err != nil {
		return nil, err
	}

	dst.Seek(0, io.SeekStart)

	return dst, nil
}

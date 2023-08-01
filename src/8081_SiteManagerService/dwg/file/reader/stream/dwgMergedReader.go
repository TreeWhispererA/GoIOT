package stream

type DwgMergedReader struct {
	IDwgStreamReader

	mainReader   IDwgStreamReader
	textReader   IDwgStreamReader
	handleReader IDwgStreamReader
}

func (reader DwgMergedReader) GetPosition() (int64, error) {
	return reader.mainReader.GetPosition()
}

func (reader DwgMergedReader) ReadByte() (byte, error) {
	return reader.mainReader.ReadByte()
}

func (reader DwgMergedReader) ReadShort() (int16, error) {
	return reader.mainReader.ReadShort()
}

func (reader DwgMergedReader) ReadBit() (bool, error) {
	return reader.mainReader.ReadBit()
}

func (reader DwgMergedReader) ReadBitAsShort() (int16, error) {
	return reader.mainReader.ReadBitAsShort()
}

func (reader DwgMergedReader) Read2Bits() (byte, error) {
	return reader.mainReader.Read2Bits()
}

func (reader DwgMergedReader) ReadBitShort() (int16, error) {
	return reader.mainReader.ReadBitShort()
}

func (reader DwgMergedReader) ReadBitShortAsBool() (bool, error) {
	return reader.mainReader.ReadBitShortAsBool()
}

func (reader DwgMergedReader) ReadBitLong() (int32, error) {
	return reader.mainReader.ReadBitLong()
}

func (reader DwgMergedReader) ReadBitLongLong() (int32, error) {
	return reader.mainReader.ReadBitLongLong()
}

func (reader DwgMergedReader) ReadBitDouble() (float64, error) {
	return reader.mainReader.ReadBitDouble()
}

func (reader DwgMergedReader) ReadSentinel() ([]byte, error) {
	return reader.mainReader.ReadSentinel()
}

func (reader DwgMergedReader) ReadVariableText() (string, error) {
	//Handle the text section if is empty
	if reader.textReader.IsEmpty() {
		return "", nil
	}

	return reader.textReader.ReadVariableText()
}

func (reader DwgMergedReader) ResetShift() (uint16, error) {
	return reader.mainReader.ResetShift()
}

func (reader DwgMergedReader) GetPositionInBits() (int32, error) {
	return reader.mainReader.GetPositionInBits()
}

func (reader DwgMergedReader) SetPositionInBits(position int32) error {
	return reader.mainReader.SetPositionInBits(position)
}

func NewDwgMergedReader(mainReader IDwgStreamReader, textReader IDwgStreamReader,
	handleReader IDwgStreamReader) IDwgStreamReader {
	return &DwgMergedReader{
		mainReader:   mainReader,
		textReader:   textReader,
		handleReader: handleReader,
	}
}

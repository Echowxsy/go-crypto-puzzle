package archive

import (
	"encoding/binary"
	"math/big"
)

type Archive struct {
	N  *big.Int
	A  *big.Int
	T  *big.Int
	Ck *big.Int
	Cm []byte
}

func ArchiveToByte(a Archive) []byte {
	writeInt32 := func(n int) []byte {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(n))
		return b
	}

	writeBigInt := func(n *big.Int) []byte {
		return n.Bytes()
	}

	nBytes := writeBigInt(a.N)
	nLength := writeInt32(len(nBytes))

	aBytes := writeBigInt(a.A)
	aLength := writeInt32(len(aBytes))

	tBytes := writeBigInt(a.T)
	tLength := writeInt32(len(tBytes))

	CkBytes := writeBigInt(a.Ck)
	CkLength := writeInt32(len(CkBytes))

	totalLength := len(nLength) + len(nBytes) + len(aLength) + len(aBytes) + len(tLength) + len(tBytes) + len(CkLength) + len(CkBytes) + len(a.Cm)
	result := make([]byte, 0, totalLength)
	result = append(result, nLength...)
	result = append(result, nBytes...)
	result = append(result, aLength...)
	result = append(result, aBytes...)
	result = append(result, tLength...)
	result = append(result, tBytes...)
	result = append(result, CkLength...)
	result = append(result, CkBytes...)
	result = append(result, a.Cm...)
	return result
}

func UnarchiveFromByte(archive []byte) Archive {
	readInt32 := func(data []byte) int {
		return int(binary.BigEndian.Uint32(data))
	}

	readBigInt := func(data []byte) *big.Int {
		return new(big.Int).SetBytes(data)
	}

	nLength := readInt32(archive[:4])
	n := readBigInt(archive[4 : 4+nLength])

	aLength := readInt32(archive[4+nLength : 4+nLength+4])
	a := readBigInt(archive[4+nLength+4 : 4+nLength+4+aLength])

	tLength := readInt32(archive[4+nLength+4+aLength : 4+nLength+4+aLength+4])
	t := readBigInt(archive[4+nLength+4+aLength+4 : 4+nLength+4+aLength+4+tLength])

	CkLength := readInt32(archive[4+nLength+4+aLength+4+tLength : 4+nLength+4+aLength+4+tLength+4])
	Ck := readBigInt(archive[4+nLength+4+aLength+4+tLength+4 : 4+nLength+4+aLength+4+tLength+4+CkLength])

	Cm := archive[4+nLength+4+aLength+4+tLength+4+CkLength:]

	return Archive{n, a, t, Ck, Cm}
}

// Package imohash implements a fast, constant-time hash for files. It is based atop
// murmurhash3 and uses file size and sample data to construct the hash.
//
// For more information, including important caveats on usage, consult https://github.com/kalafut/imohash.
package imohash

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/spaolacci/murmur3"
)

const Size = 16

// Files smaller than this will be hashed in their entirety.
const SampleThreshold = 10 * 1024 * 1024 * 1024
const SampleSize = 16 * 1024 * 1024

var emptyArray = [Size]byte{}

type ImoHash struct {
	hasher          murmur3.Hash128
	sampleSize      int
	sampleThreshold int
	bytesAdded      int
}

// New returns a new ImoHash using the default sample size
// and sample threshhold values.
func New() ImoHash {
	return NewCustom(SampleSize, SampleThreshold)
}

// New2 returns a new ImoHash using the custom sample size
// and sample threshhold values.
func New2(sampleSize, sampleThreshold int) ImoHash {
	return NewCustom(sampleSize, sampleThreshold)
}

// NewCustom returns a new ImoHash using the provided sample size
// and sample threshhold values. The entire file will be hashed
// (i.e. no sampling), if sampleSize < 1.
func NewCustom(sampleSize, sampleThreshold int) ImoHash {
	h := ImoHash{
		hasher:          murmur3.New128(),
		sampleSize:      sampleSize,
		sampleThreshold: sampleThreshold,
	}

	return h
}

// SumFile hashes a file using default sample parameters.
func SumFile(filename string) ([Size]byte, error) {
	imo := New()
	return imo.SumFile(filename)
}

// SumFile2 hashes a file using custom sample parameters.
func SumFile2(filename string, samplesize, samplethreshold int) ([Size]byte, error) {
	imo := New2(samplesize, samplethreshold)
	return imo.SumFile(filename)
}

// Sum hashes a byte slice using default sample parameters.
func Sum(data []byte) [Size]byte {
	imo := New()
	return imo.Sum(data)
}

// Sum hashes a byte slice using the ImoHash parameters.
func (imo *ImoHash) Sum(data []byte) [Size]byte {
	sr := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))

	return imo.hashCore(sr)
}

// SumFile hashes a file using using the ImoHash parameters.
func (imo *ImoHash) SumFile(filename string) ([Size]byte, error) {
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return emptyArray, err
	}

	fi, err := f.Stat()
	if err != nil {
		return emptyArray, err
	}
	sr := io.NewSectionReader(f, 0, fi.Size())
	return imo.hashCore(sr), nil
}

// hashCore hashes a SectionReader using the ImoHash parameters.
func (imo *ImoHash) hashCore(f *io.SectionReader) [Size]byte {
	var result [Size]byte

	imo.hasher.Reset()

	if f.Size() < int64(imo.sampleThreshold) || imo.sampleSize < 1 {
		buffer := make([]byte, f.Size())
		f.Read(buffer)
		imo.hasher.Write(buffer)
	} else {
		buffer := make([]byte, imo.sampleSize)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(f.Size()/2, 0)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(int64(-imo.sampleSize), 2)
		f.Read(buffer)
		imo.hasher.Write(buffer)
	}

	hash := imo.hasher.Sum(nil)

	binary.PutUvarint(hash, uint64(f.Size()))
	copy(result[:], hash)

	return result
}

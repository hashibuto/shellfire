package elf

import (
	"encoding/binary"
	"fmt"

	"github.com/hashibuto/shellfire/internal/binreader"
)

const (
	ElfHeaderSize32 = 0x34
	ElfHeaderSize64 = 0x40

	ProgramHeaderSize32 = 0x20
	ProgramHeaderSize64 = 0x40
)

type ElfHeader struct {
	IMagic      []byte // 4 bytes
	IClass      byte
	IData       byte
	IVersion    byte
	IOsAbi      byte
	IAbiVersion byte
	IPad        []byte // 7 bytes
	Type        uint16
	Machine     uint16
	Version     uint32
	Entry       uint64 // 32 or 64
	PhOffset    uint64 // 32 or 64
	ShOffset    uint64 // 32 or 64
	Flags       []byte
	HSize       uint16
	PhEntSize   uint16
	PhNum       uint16
	ShEntSize   uint16
	ShNum       uint16
	ShStrIndex  uint16
}

func LoadElfFile(data []byte) error {
	if len(data) < 5 {
		return fmt.Errorf("File too small")
	}
	var headerSize int
	arch := data[4]
	if arch == 1 {
		headerSize = ElfHeaderSize32
	} else if arch == 2 {
		headerSize = ElfHeaderSize64
	} else {
		return fmt.Errorf("Unknown architecture")
	}
	if len(data) < headerSize {
		return fmt.Errorf("File too small for header")
	}

	elfHeader := &ElfHeader{}
	err := elfHeader.Unmarshal(data[0:headerSize])
	if err != nil {
		return err
	}

	return nil
}

func (h *ElfHeader) ArchBits() int {
	if h.IClass == 1 {
		return 32
	}

	return 64
}

func (h *ElfHeader) ByteOrder() binary.ByteOrder {
	if h.IData == 1 {
		return binary.LittleEndian
	}

	return binary.BigEndian
}

// NewElfHeader returns a new ELF header object from a byte array
func NewElfHeader(data []byte) (*ElfHeader, error) {
	h := &ElfHeader{}

	b := data[5]
	var order binary.ByteOrder
	if b == 1 {
		order = binary.LittleEndian
	} else if b == 2 {
		order = binary.BigEndian
	} else {
		return nil, fmt.Errorf("Unknown endianness %d", b)
	}

	br := binreader.NewBinReader(data, order)
	h.IMagic, _ = br.ReadBytes(4)
	if h.IMagic[0] != 0x7f || h.IMagic[1] != 'E' || h.IMagic[2] != 'L' || h.IMagic[3] != 'F' {
		return nil, fmt.Errorf("Magic number incorrect")
	}
	h.IClass, _ = br.ReadByte()
	if h.IClass != 1 && h.IClass != 2 {
		return nil, fmt.Errorf("Unknown architecture")
	}
	h.IData, _ = br.ReadByte()
	if h.IData != 1 && h.IData != 2 {
		return nil, fmt.Errorf("Unknown byte order")
	}
	h.IVersion, _ = br.ReadByte()
	h.IOsAbi, _ = br.ReadByte()
	h.IAbiVersion, _ = br.ReadByte()
	h.IPad, _ = br.ReadBytes(7)
	h.Type, _ = br.ReadWord()
	h.Machine, _ = br.ReadWord()
	h.Version, _ = br.ReadDWord()
	if h.IClass == 1 {
		// 32 bit
		var eentry uint32
		var ephoff uint32
		var eshoff uint32
		eentry, _ = br.ReadDWord()
		ephoff, _ = br.ReadDWord()
		eshoff, _ = br.ReadDWord()
		h.Entry = uint64(eentry)
		h.PhOffset = uint64(ephoff)
		h.ShOffset = uint64(eshoff)
	} else {
		// 64 bit
		var eentry uint64
		var ephoff uint64
		var eshoff uint64
		eentry, _ = br.ReadQWord()
		ephoff, _ = br.ReadQWord()
		eshoff, _ = br.ReadQWord()
		h.Entry = eentry
		h.PhOffset = ephoff
		h.ShOffset = eshoff
	}
	h.Flags, _ = br.ReadBytes(4)
	h.HSize, _ = br.ReadWord()
	h.PhEntSize, _ = br.ReadWord()
	h.PhNum, _ = br.ReadWord()
	h.ShEntSize, _ = br.ReadWord()
	h.ShNum, _ = br.ReadWord()
	h.ShStrIndex, _ = br.ReadWord()

	return h, nil
}

func (h *ElfHeader) Unmarshal(data []byte) error {

	return nil
}

type ProgramHeader struct {
	Type     uint64
	Flags    []byte // 4 or 8 bytes (64 bit)
	Offset   uint64
	VAddr    uint64
	PAddr    uint64
	FileSize uint64
	MemSize  uint64
	Align    uint64
}

type programHeader32 struct {
	Type     uint32
	Flags    []byte // 4 bytes
	Offset   uint32
	VAddr    uint32
	PAddr    uint32
	FileSize uint32
	MemSize  uint32
	Align    uint32
}

func NewProgramHeader(elfHeader *ElfHeader, data []byte) (*ProgramHeader, error) {
	br := binreader.NewBinReader(data, elfHeader.ByteOrder())

	h := &ProgramHeader{}
	if elfHeader.ArchBits() == 32 {
		mh := &programHeader32{}
		mh.Type, _ = br.ReadDWord()
		mh.Offset, _ = br.ReadDWord()
		mh.VAddr, _ = br.ReadDWord()
		mh.PAddr, _ = br.ReadDWord()
		mh.FileSize, _ = br.ReadDWord()
		mh.MemSize, _ = br.ReadDWord()
		mh.Flags, _ = br.ReadBytes(4)
		mh.Align, _ = br.ReadDWord()
	} else {
		h.Type, _ = br.ReadQWord()
		h.Flags, _ = br.ReadBytes(8)
		h.Offset, _ = br.ReadQWord()
		h.VAddr, _ = br.ReadQWord()
		h.PAddr, _ = br.ReadQWord()
		h.FileSize, _ = br.ReadQWord()
		h.MemSize, _ = br.ReadQWord()
		h.Align, _ = br.ReadQWord()
	}

	return h, nil
}

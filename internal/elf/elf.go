package elf

import (
	"encoding/binary"
	"fmt"

	"github.com/hashibuto/shellfire/internal/binreader"
)

const (
	ElfHeaderSize32 = 52
	ElfHeaderSize64 = 64
)

type ElfHeader struct {
	EiMagic      []byte // 4 bytes
	EiClass      byte
	EiData       byte
	EiVersion    byte
	EiOsAbi      byte
	EiAbiVersion byte
	EiPad        []byte // 7 bytes
	EType        uint16
	EMachine     uint16
	EVersion     uint32
	EEntry       uint64 // 32 or 64
	EPhOffset    uint64 // 32 or 64
	EShOffset    uint64 // 32 or 64
	Flags        []byte
	EhSize       uint16
	EPhEntSize   uint16
	EPhNum       uint16
	EShEntSize   uint16
	EShNum       uint16
	EShStrIndex  uint16
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

func (h *ElfHeader) Bits() int {
	if h.EiClass == 1 {
		return 32
	}

	return 64
}

func (h *ElfHeader) ByteOrder() binary.ByteOrder {
	if h.EiData == 1 {
		return binary.LittleEndian
	}

	return binary.BigEndian
}

func (h *ElfHeader) Unmarshal(data []byte) error {
	b := data[5]
	var order binary.ByteOrder
	if b == 1 {
		order = binary.LittleEndian
	} else if b == 2 {
		order = binary.BigEndian
	} else {
		return fmt.Errorf("Unknown endianness %d", b)
	}

	br := binreader.NewBinReader(data, order)
	h.EiMagic, _ = br.ReadBytes(4)
	if h.EiMagic[0] != 0x7f || h.EiMagic[1] != 'E' || h.EiMagic[2] != 'L' || h.EiMagic[3] != 'F' {
		return fmt.Errorf("Magic number incorrect")
	}
	h.EiClass, _ = br.ReadByte()
	if h.EiClass != 1 && h.EiClass != 2 {
		return fmt.Errorf("Unknown architecture")
	}
	h.EiData, _ = br.ReadByte()
	if h.EiData != 1 && h.EiData != 2 {
		return fmt.Errorf("Unknown byte order")
	}
	h.EiVersion, _ = br.ReadByte()
	h.EiOsAbi, _ = br.ReadByte()
	h.EiAbiVersion, _ = br.ReadByte()
	h.EiPad, _ = br.ReadBytes(7)
	h.EType, _ = br.ReadWord()
	h.EMachine, _ = br.ReadWord()
	h.EVersion, _ = br.ReadDWord()
	if h.EiClass == 1 {
		// 32 bit
		var eentry uint32
		var ephoff uint32
		var eshoff uint32
		eentry, _ = br.ReadDWord()
		ephoff, _ = br.ReadDWord()
		eshoff, _ = br.ReadDWord()
		h.EEntry = uint64(eentry)
		h.EPhOffset = uint64(ephoff)
		h.EShOffset = uint64(eshoff)
	} else {
		// 64 bit
		var eentry uint64
		var ephoff uint64
		var eshoff uint64
		eentry, _ = br.ReadQWord()
		ephoff, _ = br.ReadQWord()
		eshoff, _ = br.ReadQWord()
		h.EEntry = eentry
		h.EPhOffset = ephoff
		h.EShOffset = eshoff
	}
	h.Flags, _ = br.ReadBytes(4)
	h.EhSize, _ = br.ReadWord()
	h.EPhEntSize, _ = br.ReadWord()
	h.EPhNum, _ = br.ReadWord()
	h.EShEntSize, _ = br.ReadWord()
	h.EShNum, _ = br.ReadWord()
	h.EShStrIndex, _ = br.ReadWord()
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

func (h *ProgramHeader) Unmarshal(elfHeader *ElfHeader, data []byte) error {

}

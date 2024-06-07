package efitools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const efivar_path string = "/sys/firmware/efi/efivars/"

type EfiVar struct {
	Name     string
	FullName string
	Data     []byte
}

type EfiVars []EfiVar

func (e EfiVars) Search(name string) (EfiVar, bool) {
	for _, efivar := range e {
		if efivar.Name == name {
			return efivar, true
		}
	}

	return EfiVar{}, false
}

const (
	HDD_MEDIA_DEVICE_PATH            uint16 = 0x0401
	CDROM_MEDIA_DEVICE_PATH          uint16 = 0x0402
	VENDOR_DEFINED_MEDIA_DEVICE_PATH uint16 = 0x0402
	FILE_PATH_MEDIA_DEVICE_PATH      uint16 = 0x0402
	PWIG_FIRMWARE_FILE               uint16 = 0x0406
	DEBUG_PORT_MSG_TYPE              uint16 = 0x030a
)

type EfiLoadOption struct {
	Attr         uint32
	FilePathLen  uint16
	Desc         string
	FilePathList []DevicePath
	OptData      []byte
}

type DevicePath interface {
	DevicePathFn() // temporary hack
	GetLength() uint16
}

type HDDMediaDevicePath struct {
	Type      byte
	SubType   byte
	Length    uint16
	PartNum   uint32 // partition number
	PartStart uint64 // starting LBA of partition
	PartSize  uint64
	PartSign  [2]uint64 // Partition signature
	PartFmt   byte
	SignType  byte
}

func (m HDDMediaDevicePath) DevicePathFn() {}

func (m HDDMediaDevicePath) GetLength() uint16 {
	return m.Length
}

type DebugPortMsgPath struct {
	Type       byte
	SubType    byte
	Length     uint16
	VendorGUID [2]uint64
}

func (d DebugPortMsgPath) DevicePathFn() {}
func (d DebugPortMsgPath) GetLength() uint16 {
	return d.Length
}

type PwigFirmwareFile struct {
	Type    byte
	SubType byte
	Length  uint16
	Data    []byte // TODO: refer to spec and implement later
}

func (p PwigFirmwareFile) DevicePathFn() {}

func (p PwigFirmwareFile) GetLength() uint16 {
	return p.Length
}

func GetEfiVars() EfiVars {
	dirents, err := os.ReadDir(efivar_path)
	check(err)

	var efivars []EfiVar
	for _, dirent := range dirents {
		var efivar EfiVar

		efivar.FullName = dirent.Name()
		efivar.Name = strings.Split(efivar.FullName, "-")[0]
		efivar.Data = readEfiVar(efivar.FullName)
		efivars = append(efivars, efivar)
	}

	return efivars
}

func readEfiVar(v string) []byte {
	fullpath := efivar_path + v
	data, err := os.ReadFile(fullpath)
	check(err)

	return data
}

func ParseLoadOption(data []byte) EfiLoadOption {
	var loadopt EfiLoadOption

	loadopt.Attr = binary.LittleEndian.Uint32(data[0:4])
	loadopt.FilePathLen = binary.LittleEndian.Uint16(data[4:6])
	data = data[8:]

	// Parse the description (null-terminated string).
	//The description is an array of 2 byte wide characters
	var offset int
	for i := 0; ; i += 2 {
		if data[i] == 0 && data[i+1] == 0 {
			offset += 2
			break
		}
		loadopt.Desc += string(data[i])
		offset += 2
	}
	data = data[offset:]

	for i := 0; i < int(loadopt.FilePathLen); i++ {
		devpath := parseDevicePath(data)
		loadopt.FilePathList = append(loadopt.FilePathList, devpath)
		offset = int(devpath.GetLength())
		data = data[offset:]
	}

	loadopt.OptData = data
	return loadopt
}

func parseDevicePath(data []byte) DevicePath {
	devtype := uint16(data[0])<<8 | uint16(data[1])

	switch devtype {
	case DEBUG_PORT_MSG_TYPE:
		var debugport DebugPortMsgPath
		readBytes(&debugport, data)
		return debugport
	case HDD_MEDIA_DEVICE_PATH:
		var hddpath HDDMediaDevicePath
		readBytes(&hddpath, data)
		return hddpath
	case PWIG_FIRMWARE_FILE:
		var pwigfw PwigFirmwareFile
		pwigfw.Type = data[0]
		pwigfw.SubType = data[1]
		pwigfw.Length = uint16(data[3])<<8 | uint16(data[2]) + 4 // n + 4 according to the spec
		pwigfw.Data = append(pwigfw.Data, data[4:pwigfw.Length-4]...)
		return pwigfw
	default:
		err("Unimplemented Device Path: type=%d, subtype=%d\n", data[0], data[1])
	}

	return DebugPortMsgPath{} // should not reach here
}

// @dst: a pointer to the destination data structure
func readBytes(dst any, src []byte) {
	srcbuf := bytes.NewReader(src)

	if err := binary.Read(srcbuf, binary.LittleEndian, dst); err != nil {
		panic(err)
	}
}

func sizeof(a any) uintptr {
	typ := reflect.TypeOf(a)
	return typ.Size()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func err(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
	os.Exit(1)
}

func usage(prog string) {
	fmt.Fprintf(os.Stderr, "usage: %s efivar\n", prog)
	fmt.Fprintf(os.Stderr, "example: %s VendorKeys\n", prog)
	os.Exit(1)
}

package fstree

/*
#include <dirent.h>
#include <stdlib.h>
#include <string.h>

struct dirent *do_readDir(DIR *dir, struct dirent *ent) {
	struct dirent *test = ent;
	int err = readdir_r(dir, ent, &test);

	if (err == 0 && test != NULL) {
		return ent;
	}

	return NULL; // final return on error on on EOF
}

*/
import "C"

import (
	"os"
	"unsafe"
)

// Dirent an interface to a reduced set of fields available via 'struct dirent'
type Dirent interface {
	Name() string     // Dirent's short name
	Type() DirentType // Dirent's type
	IsDir() bool      // Is directory
}

type DirentType uint8

const (
	DTBlockDevice = DirentType(C.DT_BLK)     // This is a block device.
	DTCharDevice  = DirentType(C.DT_CHR)     // This is a character device.
	DTDirectory   = DirentType(C.DT_DIR)     // This is a directory.
	DTFIFO        = DirentType(C.DT_FIFO)    // This is a named pipe (FIFO).
	DTSymLink     = DirentType(C.DT_LNK)     // This is a symbolic link.
	DTRegular     = DirentType(C.DT_REG)     // This is a regular file.
	DTSocket      = DirentType(C.DT_SOCK)    // This is a UNIX domain socket.
	DTUnknown     = DirentType(C.DT_UNKNOWN) // The file type is unknown.
)

func fileModeToDirentType(mode os.FileMode) DirentType {
	if mode.IsDir() {
		return DTDirectory
	} else if mode.IsRegular() {
		return DTRegular
	} else {
		if mode&os.ModeDevice != 0 {
			return DTBlockDevice
		}
		if mode&os.ModeCharDevice != 0 {
			return DTCharDevice
		}
		if mode&os.ModeNamedPipe != 0 {
			return DTFIFO
		}
		if mode&os.ModeSymlink != 0 {
			return DTSymLink
		}
		if mode&os.ModeSocket != 0 {
			return DTSocket
		}
	}
	return DTUnknown
}

func (nt DirentType) String() string {
	switch nt {
	case DTBlockDevice:
		return "BLK"
	case DTCharDevice:
		return "CHR"
	case DTDirectory:
		return "DIR"
	case DTFIFO:
		return "FIO"
	case DTSymLink:
		return "LNK"
	case DTRegular:
		return "REG"
	case DTSocket:
		return "SCK"
	default:
		return "UKN"
	}
}

type dirent struct {
	name string
	kind DirentType
}

func (n dirent) Name() string {
	return n.name
}

func (n dirent) Type() DirentType {
	return n.kind
}

func (n dirent) IsDir() bool {
	return n.kind == DTDirectory
}

func makeGoDirent(cde *C.struct_dirent) Dirent {
	return &dirent{
		name: C.GoString((*C.char)(&cde.d_name[0])),
		kind: DirentType(cde.d_type),
	}
}

func DirentFromFileInfo(info os.FileInfo) Dirent {
	return &dirent{
		name: info.Name(),                       // Set name
		kind: fileModeToDirentType(info.Mode()), // Set type
	}

}

// ReadDir a fast directory reading implementation (no lstat(3) calls)
func ReadDir(path string) ([]Dirent, error) {
	// Call into c readdir_r and get all the objects in the directory

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	dir, err := C.opendir(cpath)
	if err != nil {
		return nil, err
	}
	defer C.closedir(dir)

	// Get entries
	result := make([]Dirent, 0, 2) // at least to for . and ..
	dent := C.struct_dirent{}
	for {
		C.memset(unsafe.Pointer(&dent), 0, C.size_t(unsafe.Sizeof(dent)))
		if C.do_readDir(dir, &dent) == nil {
			break // done
		}
		result = append(result, makeGoDirent(&dent))
	}

	return result, nil
}

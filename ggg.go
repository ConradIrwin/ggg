package main

import (
	"github.com/ConradIrwin/go-dwarf"
	"debug/macho"
	"errors"
	"fmt"
	"log"
	"unsafe"
)

func GetDwarfFuncByName(name string) (*dwarf.Entry, error) {

	file, err := macho.Open("ggg")

    if err != nil {
		return nil, err
    }

	dw, err := dwarf.LoadFromMachO(file)

	if err != nil {
		return nil, err
	}

	reader := dw.Reader()

	for {
		entry, err := reader.Next()

		if err != nil {
			log.Fatal(err)
		}

		if entry == nil {
			break
		}

		if entry.Tag == dwarf.TagSubprogram && entry.Attribute(dwarf.AttrName).(string) == name {
				return entry, nil
		}
	}

	return nil, errors.New("Could not find runtime.getcallersp")
}

func main() {
	d, err := dwarf.LoadForSelf()
	if err != nil {
		panic(err)
	}

	s, err := SubprogramHere(Magic())

	if err != nil {
		panic(err)
	}

	pd, err := s.Locals()["d"].Location()
	if err != nil {
		panic(err)
	}
	ps, err := s.Locals()["s"].Location()
	if err != nil {
		panic(err)
	}
	pps, err := s.Locals()["ps"].Location()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", *(**dwarf.Data)(unsafe.Pointer(pd)) == d)
	fmt.Printf("%v\n", *(**Subprogram)(unsafe.Pointer(ps)) == s)
	fmt.Printf("%v %v\n", *(*uintptr)(unsafe.Pointer(pps)), ps)
}

package main

import (
	"github.com/ConradIrwin/go-dwarf"
	"errors"
	"runtime"
)

type Subprogram struct {
	CanonicalFrameAddress uintptr
	Entries []*dwarf.Entry
}

type Variable struct {
	Subprogram *Subprogram
	Entry *dwarf.Entry
}

func SubprogramHere(_sp magic) (*Subprogram, error) {

	sp := uintptr(_sp) + 8

	pc, _, _, ok := runtime.Caller(1)

	if !ok {
		return nil, errors.New("Caller was not OK")
	}

	return SubprogramForPC(pc, sp)
}

func SubprogramForPC(_pc uintptr, _sp uintptr) (*Subprogram, error) {

    dw, err := dwarf.LoadForSelf()
    if err != nil {
		return nil, err
    }

	cfa, err := dw.CanonicalFrameAddress(_pc, _sp)
    if err != nil {
		return nil, err
    }

    reader := dw.Reader()

	s := &Subprogram{CanonicalFrameAddress: cfa}
	pc := uint64(_pc)

	collecting := false

   for {
        entry, err := reader.Next()

        if err != nil {
			return nil, err
        }

        if entry == nil {
            break
        }

        if entry.Tag == dwarf.TagSubprogram {
            var begin, end uint64

            for _, f := range(entry.Field) {
                if f.Attr == dwarf.AttrLowpc {
                    begin = f.Val.(uint64)
                }
                if f.Attr == dwarf.AttrHighpc {
                    end = f.Val.(uint64)
                }
            }

            if begin <= pc && end >= pc {
                collecting = true
            } else if (collecting) {
				return s, nil
            }
        }

		if collecting {
			s.Entries = append(s.Entries, entry)
		}
	}


	return nil, errors.New("No program found");
}

func (s *Subprogram) Locals() (map[string]*Variable) {

	out := make(map[string]*Variable)

	for _, entry := range(s.Entries) {
		if entry.Tag == dwarf.TagVariable || entry.Tag == dwarf.TagFormalParameter {
			out[entry.Attribute(dwarf.AttrName).(string)] = &Variable{s, entry}
		}
	}

	return out
}

func (v *Variable) Location() (uintptr, error) {
	return v.Entry.Location(v.Subprogram.CanonicalFrameAddress)
}


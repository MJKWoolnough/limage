package xcf

import "io"

func (r *reader) Skip(n uint32) {
	r.rs.Seek(int64(n), io.SeekCurrent)
}

func (r *reader) SkipBoolProperty() {
	r.Skip(4)
}

func (r *reader) SkipByte() {
	r.Skip(1)
}

func (r *reader) SkipUint32() {
	r.Skip(4)
}

func (r *reader) SkipFloat32() {
	r.Skip(4)
}

func (r *reader) SkipString() {
	r.Skip(r.ReadUint32())
}

func (r *reader) SkipParasites(l uint32) {
	r.Skip(l)
}

func (r *reader) SkipParasite() {
	r.SkipString()         // name
	r.SkipUint32()         // flags
	r.Skip(r.ReadUint32()) // data[length]
}

func (r *reader) SkipPaths() {
	r.SkipUint32() // aIndex

	n := r.ReadUint32()

	for i := uint32(0); i < n; i++ {
		r.SkipString()       // name
		r.SkipBoolProperty() // linked
		r.SkipByte()         // state
		r.SkipBoolProperty() // closed

		np := r.ReadUint32()

		switch r.ReadUint32() { // version
		case 1:
		case 2:
			r.SkipUint32()
		case 3:
			r.SkipUint32()
			r.SkipUint32() // tattoo
		default:
			r.SetError(ErrUnknownPathsVersion)

			return
		}

		r.Skip(12 * np) // (control[4] + x[4] + y[4]) * np
	}
}

func (r *reader) SkipVectors() {
	if r.ReadUint32() != 1 { // version
		r.SetError(ErrUnknownVectorVersion)

		return
	}

	r.SkipUint32() // aIndex

	n := r.ReadUint32()

	for i := uint32(0); i < n; i++ {
		r.SkipString()       // name
		r.SkipUint32()       // tattoo
		r.SkipBoolProperty() // visible
		r.SkipBoolProperty() // linked

		m := r.ReadUint32()
		k := r.ReadUint32()

		for j := uint32(0); j < m; j++ {
			r.SkipParasite()
		}

		for j := uint32(0); j < k; j++ {
			if r.ReadUint32() != 1 { // stroke type
				r.SetError(ErrUnknownStrokeType)

				return
			}

			r.SkipBoolProperty() // closed

			nf := r.ReadUint32()
			np := r.ReadUint32()

			switch nf {
			case 2, 3, 4, 5, 6:
				/*
					2: Bool + Float32 + Float32 = 5
					3: +Float32 = 9
					4: +Float32 = 13
					5: +Float32 = 17
					6: +Float32 = 21
				*/
				r.Skip(np * (nf*4 + 1))
			default:
				r.SetError(ErrInvalidFloatsNumber)

				return
			}
		}
	}
}

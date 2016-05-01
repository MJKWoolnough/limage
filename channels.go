package xcf

import "os"

type channel struct{}

func (d *Decoder) readChannel() channel {
	width := d.r.ReadUint32()
	height := d.r.ReadUint32()
	name := d.r.ReadString()
Props:
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			break Props
		case propActiveChannel:
			//a := d.readActiveChannel()
			// no data, just set as active
		case propSelection:
			//s := d.readSelection()
			// no data, just set as selection
		case propOpacity:
			o := d.readOpacity()
			_ = o
		case propVisible:
			v := d.readBool()
			_ = v
		case propLinked:
			l := d.readBool()
			_ = l
		case propShowMasked:
			s := d.readBool()
			_ = s
		case propColor:
			r := d.r.ReadUint8()
			g := d.r.ReadUint8()
			b := d.r.ReadUint8()
			_, _, _ = r, g, b
		case propTattoo:
			t := d.readTattoo()
			_ = t
		case propParasites:
			p := d.readParasites(propLength)
			_ = p
		case propLockContent:
			l := d.readBool()
			_ = l
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
	hptr := d.r.ReadUint32() //
}

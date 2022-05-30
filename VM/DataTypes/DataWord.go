package DataTypes

import (
	"time"
)

type Datatype int

const (
	Boolean Datatype = 0
	Int32   Datatype = 1
	Int64   Datatype = 2
	Int256  Datatype = 3
	String  Datatype = 4
	Time    Datatype = 5
)

const ()

type DataWord struct {
	Data     ExtraBigInt
	Datatype Datatype
}

func (x DataWord) ToString() string {
	d := x.Data
	if x.Datatype >= 1 && x.Datatype <= 3 {
		return x.Data.toString()
	} else if x.Datatype == Boolean {
		if x.Data[0] == 1 {
			return "True"
		} else {
			return "False"
		}
	} else if x.Datatype == String {
		return string(x.Data.toByteArray())
	}

	t := time.Date(int(d[0]), time.Month(int(d[1])), int(d[2]), int(d[3]), int(d[4]), int(d[5]), 0, time.UTC)
	return t.GoString()
}

func NewData(datatype Datatype) *DataWord {

	var size int
	switch datatype {
	case Boolean:
		size = 1
	case Int32:
		size = 1
	case Int64:
		size = 2
	case Int256:
		size = 8
	case Time:
		size = 6
	default:
		size = 1
	}
	return &DataWord{
		Data:     NewExtraBigInt(size),
		Datatype: datatype,
	}
}

func NewString(arr []byte) *DataWord {
	var mod = 1
	if len(arr)%4 == 0 {
		mod = 0
	}
	data := &DataWord{
		Data:     NewExtraBigInt(len(arr)/4 + mod),
		Datatype: String,
	}
	data.Data.SetDataWord(arr)
	return data
}

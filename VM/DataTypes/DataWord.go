package DataTypes

import (
	"fmt"
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

type DataWord struct {
	Data     ExtraBigInt
	Datatype Datatype
}

func (x DataWord) ToString() string {
	d := x.Data
	if x.Datatype >= 1 && x.Datatype <= 3 {
		return x.Data.ToString()
	} else if x.Datatype == Boolean {
		if x.Data[0] == 1 {
			return "True"
		} else {
			return "False"
		}
	} else if x.Datatype == String {
		return string(x.Data.ToByteArray())
	}
	dataToByte := d.ToByteArray()
	var year = uint16(dataToByte[1])<<8 + uint16(dataToByte[0])
	t := time.Date(int(year), time.Month(int(d[2])), int(d[3]), int(d[4]), int(d[5]), int(d[6]), 0, time.UTC)
	return t.String()
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
		size = 2
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
	fmt.Println("HHHHH", string(arr))
	data := &DataWord{
		Data:     NewExtraBigInt(len(arr)/4 + mod),
		Datatype: String,
	}
	data.Data.SetDataWord(arr)
	return data
}

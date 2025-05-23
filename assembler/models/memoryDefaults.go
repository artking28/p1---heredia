package models

import (
	mgu "github.com/artking28/myGoUtils"
)

const NeanderPadding = 4   // 4 bytes
const JmpConstantsSize = 2 // 2 bytes

const (
	OneValue uint16 = iota + JmpConstantsSize
	ZeroValue
	MinusOneValue
	AlternateOneValue
	AcCache0Addr
	AcCache1Addr
	SiAddr
)

func GetLastConstant() uint16 {
	keys := mgu.MapKeys(GetBuiltinConstants())
	m := keys[0]
	for _, k := range keys {
		if k > m {
			m = k
		}
	}
	return m
}

func GetBuiltinConstants() map[uint16]int16 {
	return map[uint16]int16{
		OneValue + NeanderPadding - JmpConstantsSize:          1,
		ZeroValue + NeanderPadding - JmpConstantsSize:         0,
		MinusOneValue + NeanderPadding - JmpConstantsSize:     -1,
		AlternateOneValue + NeanderPadding - JmpConstantsSize: 1,
		AcCache0Addr + NeanderPadding - JmpConstantsSize:      0,
		AcCache1Addr + NeanderPadding - JmpConstantsSize:      0,
		SiAddr + NeanderPadding - JmpConstantsSize:            0,
	}
}

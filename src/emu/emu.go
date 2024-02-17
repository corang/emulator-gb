package emu

import (
	"github.com/corang/emulator-gb/src/cart"
	"github.com/corang/emulator-gb/src/common"
	"github.com/corang/emulator-gb/src/cpu"
)

func StartEmulator(filePath string) {

	_, err := cart.Load(filePath)
	if err != nil {
		panic(err)
	}

	cpu.Init()

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 		cpu.Step()
	// 	}
	// }

	common.Unimplemented()
}

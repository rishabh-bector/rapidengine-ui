package main

import (
	"rapidengine/material"

	"nuklear-golang/nk"
)

//   --------------------------------------------------
//   UI material editing window
//   --------------------------------------------------

var currentMaterial *material.Material

func leftMaterial() {
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	nk.NkLabel(ctx, "Materials", nk.TextAlignCentered|nk.TextAlignMiddle)

	nk.NkLayoutRowDynamic(ctx, 75, 1)
	if nk.NkButtonLabel(ctx, "Create Material") == 1 {
		createMaterial()
	}

	nk.NkLayoutRowDynamic(ctx, 100, 1)

	nk.NkLayoutRowDynamic(ctx, 1000, 1)
	nk.NkGroupBegin(ctx, "", nk.WindowBorder)
	for matName, _ := range engine.MaterialControl.Materials {
		nk.NkLayoutRowDynamic(ctx, 50, 1)
		nk.NkButtonLabel(ctx, matName)
	}
	nk.NkGroupEnd(ctx)
}

func rightMaterial() {
	if currentMaterial == nil {
		return
	}

}

func createMaterial() {
	engine.MaterialControl.NewPBRMaterial("ummm")
}

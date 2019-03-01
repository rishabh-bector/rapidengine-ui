package main

import (
	"fmt"
	"rapidengine/child"
	"rapidengine/lighting"
	"rapidengine/material"

	"nuklear-golang/nk"
)

//   --------------------------------------------------
//   UI material editing window
//   --------------------------------------------------

var currentMaterial material.MaterialUI

var materialViewChild *child.Child3D
var materialPointLight *lighting.PointLight

var col = nk.NkRgb(0, 0, 0)

func initMaterialView() {
	engine.TextureControl.NewTexture("./maps/col.jpg", "col", "mipmap")
	engine.TextureControl.NewTexture("./maps/nrm.jpg", "nrm", "mipmap")
	engine.TextureControl.NewTexture("./maps/disp.jpg", "disp", "mipmap")
	engine.TextureControl.NewTexture("./maps/rough.jpg", "rough", "mipmap")
	engine.TextureControl.NewTexture("./maps/ao.jpg", "ao", "mipmap")

	materialViewChild = engine.ChildControl.NewChild3D()
	currentMaterial = engine.MaterialControl.NewPBRMaterial("e")

	currentMaterial.AttachDiffuseMap(engine.TextureControl.GetTexture("col"))
	//currentMaterial.AttachNormalMap(engine.TextureControl.GetTexture("nrm"))
	//currentMaterial.AttachHeightMap(engine.TextureControl.GetTexture("disp"))
	//currentMaterial.AttachRoughnessMap(engine.TextureControl.GetTexture("rough"))
	//currentMaterial.AttachAOMap(engine.TextureControl.GetTexture("ao"))

	materialViewChild.AttachModel(engine.GeometryControl.LoadModel("../rapidengine/assets/obj/cube_smooth.obj", currentMaterial))
	/*materialViewChild.AttachModel(geometry.Model{
		Meshes:    []geometry.Mesh{geometry.NewPlane(200, 200, 100, nil, 1)},
		Materials: map[int]material.Material{0: currentMaterial},
	})
	materialViewChild.X = -100
	materialViewChild.Z = -100*/
	materialViewChild.Model.ComputeTangents()
	materialViewChild.Model.Meshes[0].ComputeTangents()
	materialViewChild.AttachMaterial(currentMaterial)

	materialPointLight = lighting.NewPointLight(
		[]float32{0.1, 0.1, 0.1},
		[]float32{500, 500, 500},
		[]float32{1, 1, 1},
		1.0, 0.2, 0.05,
	)

	materialPointLight.SetPosition([]float32{2, 5, 2})

	engine.LightControl.InstanceLight(materialPointLight, 0)
}

func leftMaterial() {
	nk.NkLayoutRowDynamic(ctx, topPanelHeight, 1)
	nk.NkLabel(ctx, "Materials", nk.TextAlignCentered|nk.TextAlignMiddle)

	nk.NkLayoutRowDynamic(ctx, 25, 1)
	if nk.NkButtonLabel(ctx, "Create Material") == 1 {
		createMaterial()
	}

	nk.NkLayoutRowDynamic(ctx, topPanelHeight, 1)

	nk.NkLayoutRowDynamic(ctx, 1000, 1)
	nk.NkGroupBegin(ctx, "", nk.WindowBorder)
	for matName, _ := range engine.MaterialControl.Materials {
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		nk.NkButtonLabel(ctx, matName)
	}
	nk.NkGroupEnd(ctx)
}

func rightMaterial() {
	if currentMaterial == nil {
		return
	}

	nk.NkLayoutRowDynamic(ctx, topPanelHeight, 1)
	nk.NkLabel(ctx, "Components", nk.TextAlignCentered|nk.TextAlignMiddle)

	ratio := []float32{0.2, 0.8}

	// Info
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)
	e := nk.NkGroupBegin(ctx, "Info", nk.WindowBorder|nk.WindowTitle)
	if e > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetScale()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetScale(), 1, 0.01)
		nk.NkGroupEnd(ctx)
	}

	// Diffuse
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)

	if nk.NkGroupBegin(ctx, "Diffuse", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetDiffuseScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetDiffuseScalar(), 1, 0.01)
		nk.NkGroupEnd(ctx)
	}

	// Normal
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)
	if nk.NkGroupBegin(ctx, "Normal", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetNormalScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetNormalScalar(), 1, 0.01)
		nk.NkGroupEnd(ctx)
	}

	// Height
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)
	if nk.NkGroupBegin(ctx, "Height", nk.WindowBorder|nk.WindowTitle) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/4, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetParallaxDisplacement()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetParallaxDisplacement(), 1, 0.01)

		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/4, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetVertexDisplacement()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetVertexDisplacement(), 10, 0.1)

		nk.NkGroupEnd(ctx)
	}

	// Metallic
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)
	if nk.NkGroupBegin(ctx, "Metallic", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetMetallicScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetMetallicScalar(), 1, 0.01)
		nk.NkGroupEnd(ctx)
	}

	// Roughness
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)
	if nk.NkGroupBegin(ctx, "Roughness", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetRoughnessScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetRoughnessScalar(), 1, 0.01)
		nk.NkGroupEnd(ctx)
	}

	// AO
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight, 1)
	if nk.NkGroupBegin(ctx, "Ambient Occlusion", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetAOScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetAOScalar(), 1, 0.01)
		nk.NkGroupEnd(ctx)
	}

	// Col
	nk.NkLayoutRowDynamic(ctx, 500, 1)
	c := nk.NkGroupBegin(ctx, "E", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar|nk.WindowMinimizable)
	if c > 0 {
		nk.NkLayoutRowDynamic(ctx, 300, 1)
		col = nk.NkColorPicker(ctx, col, nk.ColorFormatRGB)
		materialPointLight.Diffuse[0] = float32(col.R())
		materialPointLight.Diffuse[1] = float32(col.G())
		materialPointLight.Diffuse[2] = float32(col.B())
		nk.NkGroupEnd(ctx)
	}

	// Render child
	engine.Renderer.RenderChild(materialViewChild)
}

func createMaterial() {
	currentMaterial = engine.MaterialControl.NewPBRMaterial("ummm")
	materialViewChild.Model.Materials[0] = currentMaterial
}

package main

import (
	"fmt"
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/lighting"
	"rapidengine/material"

	"nuklear-golang/nk"
)

//   --------------------------------------------------
//   UI material editing window
//   --------------------------------------------------

var currentMaterial material.MaterialUI

var materialViewChild *child.Child3D

func initMaterialView() {
	engine.TextureControl.NewTexture("./maps/col.jpg", "col", "mipmap")
	engine.TextureControl.NewTexture("./maps/nrm.jpg", "nrm", "mipmap")
	engine.TextureControl.NewTexture("./maps/disp.jpg", "disp", "mipmap")
	engine.TextureControl.NewTexture("./maps/rough.jpg", "rough", "mipmap")
	engine.TextureControl.NewTexture("./maps/ao.jpg", "ao", "mipmap")

	materialViewChild = engine.ChildControl.NewChild3D()
	currentMaterial = engine.MaterialControl.NewPBRMaterial("e")

	currentMaterial.AttachDiffuseMap(engine.TextureControl.GetTexture("col"))
	currentMaterial.AttachNormalMap(engine.TextureControl.GetTexture("nrm"))
	currentMaterial.AttachHeightMap(engine.TextureControl.GetTexture("disp"))
	currentMaterial.AttachRoughnessMap(engine.TextureControl.GetTexture("rough"))
	currentMaterial.AttachAOMap(engine.TextureControl.GetTexture("ao"))

	materialViewChild.AttachModel(engine.GeometryControl.LoadModel("../rapidengine/assets/obj/sphere_uv.obj", currentMaterial))
	materialViewChild.AttachModel(geometry.Model{
		Meshes:    []geometry.Mesh{geometry.NewPlane(20, 20, 1000, nil, 1)},
		Materials: map[int]material.Material{0: currentMaterial},
	})
	materialViewChild.X = -10
	materialViewChild.Z = -10
	materialViewChild.Model.ComputeTangents()
	materialViewChild.AttachMaterial(currentMaterial)

	l := lighting.NewPointLight(
		[]float32{0.1, 0.1, 0.1},
		[]float32{100, 100, 100},
		[]float32{1, 1, 1},
		1.0, 0.2, 0.05,
	)

	l.SetPosition([]float32{0, 2, 0})

	engine.LightControl.InstanceLight(l, 0)
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
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	nk.NkGroupBegin(ctx, "Info", nk.WindowBorder|nk.WindowTitle)
	nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
	nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetScale()), nk.TextAlignCentered|nk.TextAlignMiddle)
	nk.NkSliderFloat(ctx, 0, currentMaterial.GetScale(), 1, 0.01)
	nk.NkGroupEnd(ctx)

	// Diffuse
	nk.NkLayoutRowDynamic(ctx, 100, 1)

	if nk.NkGroupBegin(ctx, "Diffuse", nk.WindowBorder|nk.WindowTitle|nk.WindowMinimizable|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetDiffuseScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetDiffuseScalar(), 1, 0.01)
	}
	nk.NkGroupEnd(ctx)

	// Normal
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	if nk.NkGroupBegin(ctx, "Normal", nk.WindowBorder|nk.WindowTitle|nk.WindowMinimizable|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetNormalScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetNormalScalar(), 1, 0.01)
	}
	nk.NkGroupEnd(ctx)

	// Height
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	if nk.NkGroupBegin(ctx, "Height", nk.WindowBorder|nk.WindowTitle|nk.WindowMinimizable) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetVertexDisplacement()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetVertexDisplacement(), 10, 0.1)

		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetParallaxDisplacement()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetParallaxDisplacement(), 10, 0.1)
	}
	nk.NkGroupEnd(ctx)

	// Metallic
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	if nk.NkGroupBegin(ctx, "Metallic", nk.WindowBorder|nk.WindowTitle|nk.WindowMinimizable|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetMetallicScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetMetallicScalar(), 1, 0.01)
	}
	nk.NkGroupEnd(ctx)

	// Roughness
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	if nk.NkGroupBegin(ctx, "Roughness", nk.WindowBorder|nk.WindowTitle|nk.WindowMinimizable|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetRoughnessScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetRoughnessScalar(), 1, 0.01)
	}
	nk.NkGroupEnd(ctx)

	// AO
	nk.NkLayoutRowDynamic(ctx, 100, 1)
	if nk.NkGroupBegin(ctx, "Ambient Occlusion", nk.WindowBorder|nk.WindowTitle|nk.WindowMinimizable|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetAOScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetAOScalar(), 1, 0.01)
	}
	nk.NkGroupEnd(ctx)

	// Render child
	engine.Renderer.RenderChild(materialViewChild)
}

func createMaterial() {
	currentMaterial = engine.MaterialControl.NewPBRMaterial("ummm")
	materialViewChild.Model.Materials[0] = currentMaterial
}

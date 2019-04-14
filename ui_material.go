package main

import (
	"fmt"
	"rapidengine/child"
	"rapidengine/geometry"
	"rapidengine/lighting"
	"rapidengine/material"
	"sort"

	"nuklear-golang/nk"
)

//   --------------------------------------------------
//   UI material editing window
//   --------------------------------------------------

var currentMaterial material.MaterialUI
var allMaterials []string

var materialViewChild *child.Child3D
var materialBoxChild *child.Child3D
var materialPointLight *lighting.PointLight

var col = nk.NkRgb(0, 0, 0)
var col2 = nk.NkRgb(0, 0, 0)

var roughORsmooth int32

func initMaterialView() {

	//   --------------------------------------------------
	//   View Children
	//   --------------------------------------------------

	materialViewChild = engine.ChildControl.NewChild3D()
	materialViewChild.AttachModel(geometry.NewModel(geometry.NewPlane(10, 10, 100, nil, 1), currentMaterial))
	materialViewChild.Model.Meshes[0].ComputeTangents()
	materialViewChild.X -= 5
	materialViewChild.Z -= 5

	boxMat := engine.MaterialControl.NewPBRMaterial("box")
	engine.TextureControl.NewTexture("./maps/concrete/col.jpg", "concrete_col", "mipmap")
	engine.TextureControl.NewTexture("./maps/concrete/nrm.jpg", "concrete_nrm", "mipmap")
	engine.TextureControl.NewTexture("./maps/concrete/rough.jpg", "concrete_rough", "mipmap")
	boxMat.AlbedoMap = engine.TextureControl.GetTexture("concrete_col")
	boxMat.NormalMap = engine.TextureControl.GetTexture("concrete_nrm")
	boxMat.RoughnessMap = engine.TextureControl.GetTexture("concrete_rough")

	materialBoxChild = engine.ChildControl.NewChild3D()
	materialBoxChild.AttachModel(engine.GeometryControl.LoadModel("../rapidengine/assets/obj/viewer.obj", boxMat))
	materialBoxChild.Model.Meshes[0].ComputeTangents()

	//   --------------------------------------------------
	//   Rocks
	//   --------------------------------------------------

	engine.TextureControl.NewTexture("./maps/rocks1/col.jpg", "col", "mipmap")
	engine.TextureControl.NewTexture("./maps/rocks1/nrm.jpg", "nrm", "mipmap")
	engine.TextureControl.NewTexture("./maps/rocks1/disp.jpg", "disp", "mipmap")
	engine.TextureControl.NewTexture("./maps/rocks1/rough.jpg", "rough", "mipmap")
	engine.TextureControl.NewTexture("./maps/rocks1/ao.jpg", "ao", "mipmap")

	createMaterial("Rocks")

	currentMaterial.AttachDiffuseMap(engine.TextureControl.GetTexture("col"))
	currentMaterial.AttachNormalMap(engine.TextureControl.GetTexture("nrm"))
	currentMaterial.AttachHeightMap(engine.TextureControl.GetTexture("disp"))
	currentMaterial.AttachRoughnessMap(engine.TextureControl.GetTexture("rough"))
	currentMaterial.AttachAOMap(engine.TextureControl.GetTexture("ao"))

	materialPointLight = lighting.NewPointLight(
		[]float32{0.1, 0.1, 0.1},
		[]float32{0.5, 0.5, 0.5},
		[]float32{0, 0, 0},
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
		createMaterial("e")
	}

	nk.NkLayoutRowDynamic(ctx, topPanelHeight, 1)

	nk.NkLayoutRowDynamic(ctx, 1000, 1)
	nk.NkGroupBegin(ctx, "", nk.WindowBorder)

	for _, matName := range allMaterials {
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		if nk.NkButtonLabel(ctx, matName) == 1 {
			selectMaterial(matName)
		}
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
	nk.NkLayoutRowDynamic(ctx, componentGroupHeight*2, 1)
	if nk.NkGroupBegin(ctx, "Roughness", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {
		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("%v", *currentMaterial.GetRoughnessScalar()), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, 0, currentMaterial.GetRoughnessScalar(), 1, 0.01)

		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight, 1, ratio)
		roughORsmooth = nk.NkCheckLabel(ctx, "Rough/Gloss", roughORsmooth)
		if roughORsmooth == 1 {
			currentMaterial.SetRoughOrSmooth(true)
		} else {
			currentMaterial.SetRoughOrSmooth(false)
		}

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
	nk.NkLayoutRowDynamic(ctx, 300, 1)
	c := nk.NkGroupBegin(ctx, "E", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar|nk.WindowMinimizable)
	if c > 0 {
		nk.NkLayoutRowDynamic(ctx, 200, 1)
		col = nk.NkColorPicker(ctx, col, nk.ColorFormatRGB)
		materialPointLight.Diffuse[0] = float32(col.R())
		materialPointLight.Diffuse[1] = float32(col.G())
		materialPointLight.Diffuse[2] = float32(col.B())
		nk.NkGroupEnd(ctx)
	}

	// Col
	nk.NkLayoutRowDynamic(ctx, 300, 1)
	if nk.NkGroupBegin(ctx, "F", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar|nk.WindowMinimizable) > 0 {
		nk.NkLayoutRowDynamic(ctx, 200, 1)
		col2 = nk.NkColorPicker(ctx, col2, nk.ColorFormatRGB)
		//engine.LightControl.DirLight[0].Diffuse[0] = float32(col2.R()) * 100
		//engine.LightControl.DirLight[0].Diffuse[1] = float32(col2.G()) * 100
		//engine.LightControl.DirLight[0].Diffuse[2] = float32(col2.B()) * 100
		nk.NkGroupEnd(ctx)
	}

	nk.NkLayoutRowDynamic(ctx, componentGroupHeight*3, 1)
	if nk.NkGroupBegin(ctx, "Dirlight", nk.WindowBorder|nk.WindowTitle|nk.WindowNoScrollbar) > 0 {

		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("X: %v", engine.LightControl.DirLight[0].Direction[0]), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, -1, &engine.LightControl.DirLight[0].Direction[0], 1, 0.01)

		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("Y: %v", engine.LightControl.DirLight[0].Direction[1]), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, -1, &engine.LightControl.DirLight[0].Direction[1], 1, 0.01)

		nk.NkLayoutRow(ctx, nk.Dynamic, componentGroupHeight/2, 2, ratio)
		nk.NkLabel(ctx, fmt.Sprintf("Z: %v", engine.LightControl.DirLight[0].Direction[2]), nk.TextAlignCentered|nk.TextAlignMiddle)
		nk.NkSliderFloat(ctx, -1, &engine.LightControl.DirLight[0].Direction[2], 1, 0.01)

		nk.NkGroupEnd(ctx)
	}

	// Render child
	engine.Renderer.RenderChild(materialBoxChild)
	engine.Renderer.RenderChild(materialViewChild)
}

func createMaterial(name string) {
	currentMaterial = engine.MaterialControl.NewPBRMaterial(name)
	materialViewChild.Model.Materials[0] = currentMaterial
	updateMatList()
}

func selectMaterial(mat string) {
	currentMaterial = engine.MaterialControl.Materials[mat].(material.MaterialUI)
	materialViewChild.Model.Materials[0] = currentMaterial
}

func updateMatList() {
	allMaterials = []string{}
	for matName := range engine.MaterialControl.Materials {
		allMaterials = append(allMaterials, matName)
	}
	sort.Strings(allMaterials)
}

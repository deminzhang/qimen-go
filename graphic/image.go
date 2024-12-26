package graphic

import (
	"github.com/deminzhang/qimen-go/asset"
	"github.com/deminzhang/qimen-go/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func NewRectImage(size int) *ebiten.Image {
	e := ebiten.NewImage(size, size)
	vector.DrawFilledRect(e, 0, 0, float32(size), float32(size), color.White, true)
	return e
}
func NewCircleImage(size int) *ebiten.Image {
	e := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(e, float32(size/2), float32(size/2), float32(size/2), color.White, true)
	return e
}

// 太极
func NewTaiJiImage(size int) *ebiten.Image {
	halfSize := float32(size / 2)
	tj := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(tj, halfSize, halfSize, halfSize, color.White, true)                  //阳大
	vector.DrawFilledRect(tj, halfSize, 0, halfSize, float32(size), color.Black, false)           //阴半
	vector.DrawFilledCircle(tj, halfSize, float32(size/4), float32(size/4), color.White, true)    //太阳
	vector.DrawFilledCircle(tj, halfSize, float32(size*3/4), float32(size/4), color.Black, true)  //太阴
	vector.DrawFilledCircle(tj, halfSize, float32(size*2/8), float32(size/16), color.Black, true) //少阴
	vector.DrawFilledCircle(tj, halfSize, float32(size*6/8), float32(size/16), color.White, true) //少阳

	cover := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(cover, halfSize, halfSize, halfSize, color.Black, true) //mask
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceIn
	cover.DrawImage(tj, op)
	return cover
}

// 八卦
func NewBaGuaImage(gua string, size int) *ebiten.Image {
	bg := ebiten.NewImage(size, size)
	h := float32(size) / 5
	ll := float32(size) * 2 / 5
	drawUpS := func() {
		vector.DrawFilledRect(bg, 0, 0, float32(size), h, color.White, true)
	}
	drawMidS := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*2/5, float32(size), h, color.White, false)
	}
	drawDownS := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*4/5, float32(size), h, color.White, false)
	}
	drawUpL := func() {
		vector.DrawFilledRect(bg, 0, 0, ll, h, color.White, false)
		vector.DrawFilledRect(bg, ll+h, 0, ll, h, color.White, false)
	}
	drawMidL := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*2/5, ll, float32(size)/5, color.White, false)
		vector.DrawFilledRect(bg, ll+h, float32(size)*2/5, ll, h, color.White, false)
	}
	drawDownL := func() {
		vector.DrawFilledRect(bg, 0, float32(size)*4/5, ll, float32(size)/5, color.White, false)
		vector.DrawFilledRect(bg, ll+h, float32(size)*4/5, ll, h, color.White, false)
	}
	switch gua {
	case "乾":
		drawUpS()
		drawMidS()
		drawDownS()
	case "坤":
		drawUpL()
		drawMidL()
		drawDownL()
	case "艮":
		drawUpS()
		drawMidL()
		drawDownL()
	case "兑":
		drawUpL()
		drawMidS()
		drawDownS()
	case "震":
		drawUpL()
		drawMidL()
		drawDownS()
	case "巽":
		drawUpS()
		drawMidS()
		drawDownL()
	case "坎":
		drawUpL()
		drawMidS()
		drawDownL()
	case "离":
		drawUpS()
		drawMidL()
		drawDownS()
	default:
		bg.DrawImage(NewTaiJiImage(size), nil)
	}
	return bg
}

// 六十四卦
func New64GuaImage(up, down string, size int) *ebiten.Image {
	img := ebiten.NewImage(size, size*2+size/5)
	img.DrawImage(NewBaGuaImage(up, size), nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(size+size/5))
	img.DrawImage(NewBaGuaImage(down, size), op)
	return img
}

// 旗帜
func NewFlagImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	vector.StrokeLine(img, 4, 0, float32(4), float32(size), 1, color.White, true)
	vector.StrokeLine(img, 4, 0, float32(size), float32(size)/2, 1, color.White, true)
	vector.StrokeLine(img, 4, float32(size)/2, float32(size), float32(size)/2, 1, color.White, true)
	return img
}

// 地球
func NewEarthImage(size int) *ebiten.Image {
	e := ebiten.NewImage(size, size)
	hs := float32(size / 2)
	vector.DrawFilledCircle(e, hs, hs, float32(size/4), color.White, true)
	return e
}

// 太阳
func NewSunImage(size int) *ebiten.Image {
	sun := ebiten.NewImage(size, size)
	hs := float32(size / 2)
	vector.DrawFilledCircle(sun, hs, hs, float32(size/4), color.White, true)
	for i := 0; i < 8; i++ {
		y, x := util.CalRadiansPos(hs, hs, hs, float32(i*45))
		vector.StrokeLine(sun, hs, hs, x, y, 0.5, color.White, true)
	}
	return sun
}

// 月亮
func NewMoonImage(size int) *ebiten.Image {
	moon := ebiten.NewImage(size, size)
	hs := float32(size / 2)
	vector.DrawFilledCircle(moon, hs, hs, float32(size/3), color.White, true)

	cover := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(cover, float32(size/4), hs, float32(size/3), color.Black, true)
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOut
	cover.DrawImage(moon, op)
	return cover
}

// 火星
func NewMarsImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, "火", ft, 0, size, color.White)
	return img
}

// 木星
func NewJupiterImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, "木", ft, 0, size, color.White)
	return img
}

// 土星
func NewSaturnImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, "土", ft, 0, size, color.White)
	return img
}

// 水星
func NewMercuryImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, "水", ft, 0, size, color.White)
	return img
}

// 金星
func NewVenusImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, "金", ft, 0, size, color.White)
	return img
}

func NewTextImage(txt string, size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, txt, ft, 0, size, color.White)
	return img
}

func NewCampImage(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	c := color.White
	u := float32(size) / 8
	u9 := u * 7
	sm := float32(size - 1)
	hf := float32(size) / 2
	vector.StrokeLine(img, u, sm, float32(size)-u, float32(size-1), 1, c, true)  //地基
	vector.StrokeLine(img, u, sm, u, u9, 1, c, true)                             //右底
	vector.StrokeLine(img, float32(size)-u, sm, float32(size)-u, u9, 1, c, true) //左底
	vector.StrokeLine(img, float32(1), sm, hf, float32(size)/6, 1, c, true)      //左外墙
	vector.StrokeLine(img, sm, sm, hf, float32(size)/6, 1, c, true)              //右外墙
	vector.StrokeLine(img, hf, float32(1), hf, float32(size)/6, 1, c, true)      //外尖
	vector.StrokeLine(img, float32(1), sm, float32(1), u9, 1, c, true)           //左钉
	vector.StrokeLine(img, sm, sm, sm, u9, 1, c, true)                           //右钉
	vector.StrokeLine(img, u*2, sm, hf, float32(size)/2, 1, c, true)             //左门
	vector.StrokeLine(img, sm-u*2, sm, hf, float32(size)/2, 1, c, true)          //右门
	vector.StrokeLine(img, hf, float32(size)/3, hf, float32(size)/2, 1, c, true) //门尖
	return img
}

func NewArmyImage(name string, size, action int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	ft, _ := asset.GetDefaultFontFace(float64(size))
	text.Draw(img, name, ft, 0, size*7/8, color.White)
	switch action {
	case 1:
		vector.StrokeLine(img, 1, float32(size)*3/4, float32(size), float32(size)/4, 1, color.White, true) //横兵
	default:
		vector.StrokeLine(img, float32(size)/5, float32(1), float32(size)/5, float32(size), 1, color.White, true) //立兵
	}
	return img
}

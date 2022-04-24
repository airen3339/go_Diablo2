package engine

import (
	"fmt"
	"game/tools"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//Draw Login Update
func (g *Game) ChangeScenceLoginUpdate() {
	//切换场景判定
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if mouseX > 286 && mouseX < 503 && mouseY > 150 && mouseY < 218 {
			g.status.ChangeScenceFlg = true
			//切换场景
			g.ChangeScene("select")
			g.status.ChangeScenceFlg = false
		}
	}
	//音乐控制
	if !g.status.MusicIsPlay {
		g.status.MusicIsPlay = true
		g.music.PlayMusic("Act0-Intro.mp3", tools.MUSICMP3)
	}
	g.count++
	//Change Frame
	if g.count > 2 {
		counts++
		g.count = 0
		if counts >= 30 {
			counts = 0
		}
	}
}

//Draw Login Scence
func (g *Game) ChangeScenceLoginDraw(screen *ebiten.Image) {
	//Draw UI
	g.ui.DrawUI(screen)

	//Draw Logo Left
	name := "logoFireLeft_" + strconv.Itoa(counts) + ".png"
	left, _, _ := g.ui.GetAnimator("logo", name)
	opLo := &ebiten.DrawImageOptions{}
	opLo.Filter = ebiten.FilterLinear
	opLo.GeoM.Translate(220, 0)
	opLo.CompositeMode = ebiten.CompositeModeLighter
	opLo.GeoM.Scale(1, 0.7)
	screen.DrawImage(left, opLo)
	//Draw Logo Right
	name = "logoFireRight_" + strconv.Itoa(counts) + ".png"
	right, _, _ := g.ui.GetAnimator("logo", name)
	opRo := &ebiten.DrawImageOptions{}
	opRo.Filter = ebiten.FilterLinear
	opRo.GeoM.Translate(float64(220+right.Bounds().Max.X), 0)
	opRo.CompositeMode = ebiten.CompositeModeLighter
	opRo.GeoM.Scale(1, 0.7)
	screen.DrawImage(right, opRo)
	if g.status.DisPlayDebugInfo {
		//Draw Debug
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nmouse position %d,%d",
			int64(ebiten.CurrentFPS()), mouseX, mouseY))
	}
	//Change Frame
	if g.count > frameSpeed {
		counts++
		g.count = 0
		if counts >= 30 {
			counts = 0
		}
	}
}
package engine

import (
	"fmt"
	"game/status"
	"game/tools"
	"runtime"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var clearFlg = false

//Draw OpenDoor Scence
func (g *Game) ChangeScenceOpenDoorDraw(screen *ebiten.Image) {
	if !status.Config.DoorCountFlg {
		counts = 0
		status.Config.DoorCountFlg = true
	}
	//Draw Open Door
	name := "loading_" + strconv.Itoa(counts)
	door, _, _ := g.ui.GetAnimator(tools.PlistN, name)
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(268, 120)
	screen.DrawImage(door, op)
	//Draw Debug
	if status.Config.DisPlayDebugInfo {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nmouse position %d,%d", int64(ebiten.CurrentFPS()), mouseX, mouseY))
	}
}

//Draw OpenDoor Update
func (g *Game) ChangeScenceOpenDoorUpdate() {
	if !clearFlg {
		clearFlg = true
		go func() {
			g.ui.GCSelectBGImage()
			g.ui.ClearSlice(0)
			g.music.CloseBGMusic()
			runtime.GC()
		}()
	}
	g.count++
	//Change Frame
	if g.count > 10 && counts != 9 {
		counts++
		g.count = 0
	}

	// Change Scence
	if counts == 9 && !status.Config.LoadingFlg {
		status.Config.LoadingFlg = true
		w := sync.WaitGroup{}
		w.Add(1)
		go func() {
			//close music
			status.Config.MusicIsPlay = false
			g.ChangeScene("game")
			runtime.GC()
			w.Done()
		}()
		w.Wait()
		status.Config.ChangeScenceFlg = false
	}
}

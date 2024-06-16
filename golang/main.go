package main

import (
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	TileSize = 30;
	FPS = 30;
	Sleep = 1000 / FPS;
)

const (
	TileAir int = iota
	TileFlux
	TileUnbreakable
	TilePlayer
	TileStone
	TileFallingStone
	TileBox
	TileFallingBox
	TileKey1
	TileLock1
	TileKey2
	TileLock2
)

type Input interface {
	isRight() bool
	isLeft() bool
	isUp() bool
	isDown() bool
	handle()
}

type Right struct {}

func (r Right) isRight() bool {
	return true
}

func (r Right) isLeft() bool {
	return false
}

func (r Right) isUp() bool {
	return false
}

func (r Right) isDown() bool {
	return false
}

func (r Right) handle() {
	moveHorizontal(1)
}

type Left struct {}

func (l Left) isRight() bool {
	return false
}

func (l Left) isLeft() bool {
	return true
}

func (l Left) isUp() bool {
	return false
}

func (l Left) isDown() bool {
	return false
}

func (l Left) handle() {
	moveHorizontal(-1)
}

type Up struct {}

func (u Up) isRight() bool {
	return false
}

func (u Up) isLeft() bool {
	return false
}

func (u Up) isUp() bool {
	return true
}

func (u Up) isDown() bool {
	return false
}

func (u Up) handle() {
	moveVertical(-1)
}

type Down struct {}

func (d Down) isRight() bool {
	return false
}

func (d Down) isLeft() bool {
	return false
}

func (d Down) isUp() bool {
	return false
}

func (d Down) isDown() bool {
	return true
}

func (d Down) handle() {
	moveVertical(1)
}

var (
	playerx int = 1
	playery int = 1
	gameMap [][]int = [][]int{
		{2, 2, 2, 2, 2, 2, 2, 2},
		{2, 3, 0, 1, 1, 2, 0, 2},
		{2, 4, 2, 6, 1, 2, 0, 2},
		{2, 8, 4, 1, 1, 2, 0, 2},
		{2, 4, 1, 1, 1, 9, 0, 2},
		{2, 2, 2, 2, 2, 2, 2, 2},
	}
	inputs []Input = []Input{}
	inputMutex = sync.Mutex{}
)

func main() {
	go gameLoop()

	keyboard.Open()
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyEsc {
			break
		} else if char == 'w' || key == keyboard.KeyArrowUp {
			inputMutex.Lock()
			inputs = append(inputs, Up{})
			inputMutex.Unlock()
		} else if char == 's' || key == keyboard.KeyArrowDown {
			inputMutex.Lock()
			inputs = append(inputs, Down{})
			inputMutex.Unlock()
		} else if char == 'a' || key == keyboard.KeyArrowLeft {
			inputMutex.Lock()
			inputs = append(inputs, Left{})
			inputMutex.Unlock()
		} else if char == 'd' || key == keyboard.KeyArrowRight {
			inputMutex.Lock()
			inputs = append(inputs, Right{})
			inputMutex.Unlock()
		}
	}
}

func remove(tile int) {
	for y := 0; y < len(gameMap); y++ {
		for x := 0; x < len(gameMap[y]); x++ {
			if gameMap[y][x] == tile {
				gameMap[y][x] = TileAir
			}
		}
	}
}

func moveToTile(newx int, newy int) {
	gameMap[playery][playerx] = TileAir
	gameMap[newy][newx] = TilePlayer
	playerx = newx
	playery = newy
}

func moveHorizontal(dx int) {
	if gameMap[playery][playerx+dx] == TileFlux || gameMap[playery][playerx+dx] == TileAir {
		moveToTile(playerx+dx, playery)
	} else if (gameMap[playery][playerx+dx] == TileStone || gameMap[playery][playerx+dx] == TileBox) && gameMap[playery][playerx+dx+dx] == TileAir && gameMap[playery+1][playerx+dx] != TileAir {
		gameMap[playery][playerx+dx+dx] = gameMap[playery][playerx+dx]
		moveToTile(playerx+dx, playery)
	} else if gameMap[playery][playerx+dx] == TileKey1 {
		remove(TileLock1)
		moveToTile(playerx+dx, playery)
	} else if gameMap[playery][playerx+dx] == TileKey2 {
		remove(TileLock2)
		moveToTile(playerx+dx, playery)
	}
}

func moveVertical(dy int) {
	if gameMap[playery+dy][playerx] == TileFlux || gameMap[playery+dy][playerx] == TileAir {
		moveToTile(playerx, playery+dy)
	} else if gameMap[playery+dy][playerx] == TileKey1 {
		remove(TileLock1)
		moveToTile(playerx, playery+dy)
	} else if gameMap[playery+dy][playerx] == TileKey2 {
		remove(TileLock2)
		moveToTile(playerx, playery+dy)
	}
}

func update() {
	handleInputs()
	updateMap()
}

func handleInputs() {
	for len(inputs) > 0 {
		inputMutex.Lock()
		current := inputs[len(inputs)-1]
		inputs = inputs[:len(inputs)-1]
		inputMutex.Unlock()
		current.handle()
	}
}

func updateMap() {
	for y := len(gameMap) - 1; y >= 0; y-- {
		for x := 0; x < len(gameMap[y]); x++ {
			updateTile(x, y)
		}
	}
}

func updateTile(x, y int) {
	if (gameMap[y][x] == TileStone || gameMap[y][x] == TileFallingStone) && gameMap[y+1][x] == TileAir {
		gameMap[y+1][x] = TileFallingStone
		gameMap[y][x] = TileAir
	} else if (gameMap[y][x] == TileBox || gameMap[y][x] == TileFallingBox) && gameMap[y+1][x] == TileAir {
		gameMap[y+1][x] = TileFallingBox
		gameMap[y][x] = TileAir
	} else if gameMap[y][x] == TileFallingStone {
		gameMap[y][x] = TileStone
	} else if gameMap[y][x] == TileFallingBox {
		gameMap[y][x] = TileBox
	}
}

func createGraphics() CanvasRenderingContext2D {
	canvas := GetElementById("GameCanvas")
	g := canvas.GetContext("2d")
	g.ClearRect(0, 0, canvas.Width, canvas.Height);
	return g
}

func draw() {
	g := createGraphics()
	drawMap(g)
	drawPlayer(g)
}

func drawMap(g CanvasRenderingContext2D) {
	for y := 0; y < len(gameMap); y++ {
		for x := 0; x < len(gameMap[y]); x++ {
			if gameMap[y][x] == TileFlux {
				g.FillStyle = "#ccffcc";
			} else if gameMap[y][x] == TileUnbreakable {
				g.FillStyle = "#999999";
			} else if gameMap[y][x] == TileStone || gameMap[y][x] == TileFallingStone {
				g.FillStyle = "#0000cc";
			} else if gameMap[y][x] == TileBox || gameMap[y][x] == TileFallingBox {
				g.FillStyle = "#8b4513";
			} else if gameMap[y][x] == TileKey1 || gameMap[y][x] == TileLock1 {
				g.FillStyle = "#ffcc00";
			} else if gameMap[y][x] == TileKey2 || gameMap[y][x] == TileLock2 {
				g.FillStyle = "#00ccff";
			}

			if gameMap[y][x] != TileAir && gameMap[y][x] != TilePlayer {
				g.FillRect(x * TileSize, y * TileSize, TileSize, TileSize);
			}
		}
	}
}

func drawPlayer(g CanvasRenderingContext2D) {
	g.FillStyle = "#ff0000";
	g.FillRect(playerx * TileSize, playery * TileSize, TileSize, TileSize);
}

func gameLoop() {
	for {
		before := time.Now()
		update()
		draw()
		after := time.Now()
		frameTime := after.Sub(before).Milliseconds()
		sleep := Sleep - frameTime
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

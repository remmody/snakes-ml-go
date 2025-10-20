package game

import (
	"fmt"
	"image/color"

	"snakes-ml/internal/snake"
	"snakes-ml/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Renderer struct {
	screenWidth  int
	screenHeight int
}

func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		screenWidth:  width,
		screenHeight: height,
	}
}

func (r *Renderer) DrawSnake(screen *ebiten.Image, s *snake.Snake) {
	// ✅ Адаптировано под 1280x720
	gridX, gridY := 30, 150
	cellSize := 25

	maxWidth := r.screenWidth - gridX*2
	maxHeight := r.screenHeight - gridY - 80

	cellWidth := maxWidth / s.Width()
	cellHeight := maxHeight / s.Height()

	if cellWidth < cellHeight {
		cellSize = cellWidth
	} else {
		cellSize = cellHeight
	}

	if cellSize < 8 {
		cellSize = 8
	}
	if cellSize > 30 {
		cellSize = 30
	}

	totalWidth := s.Width() * cellSize
	// totalHeight := s.Height() * cellSize
	gridX = (r.screenWidth - totalWidth) / 2

	// Grid
	for y := 0; y < s.Height(); y++ {
		for x := 0; x < s.Width(); x++ {
			posX := float32(gridX + x*cellSize)
			posY := float32(gridY + y*cellSize)
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 1, ui.Grid, false)
		}
	}

	// Obstacles (Yellow)
	for _, obs := range s.Obstacles() {
		posX := float32(gridX + obs.X*cellSize)
		posY := float32(gridY + obs.Y*cellSize)
		vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), ui.Obstacle, false)
		vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, ui.ObstacleBorder, false)
	}

	// Food
	food := s.Food()
	posX := float32(gridX + food.X*cellSize)
	posY := float32(gridY + food.Y*cellSize)
	vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), ui.Food, false)
	vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, ui.FoodBorder, false)

	// Snake
	for i, segment := range s.Body() {
		posX := float32(gridX + segment.X*cellSize)
		posY := float32(gridY + segment.Y*cellSize)

		var col color.RGBA
		if i == 0 {
			col = ui.SnakeHead
		} else {
			intensity := uint8(200 - i*2)
			if intensity < 100 {
				intensity = 100
			}
			col = color.RGBA{50, intensity, 50, 255}
		}

		vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), col, false)

		if i == 0 {
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, color.RGBA{50, 200, 50, 255}, false)
		}
	}

	// Info
	occupancy := s.GetOccupancy() * 100
	infoText := fmt.Sprintf("Length: %d | Steps: %d | Map: %dx%d | Occupancy: %.1f%% | Obstacles: %d",
		s.Length(), s.Steps(), s.Width(), s.Height(), occupancy, len(s.Obstacles()))

	vector.FillRect(screen, float32(gridX-5), float32(gridY-40), float32(totalWidth+10), 35, ui.TextBg, false)
	ebitenutil.DebugPrintAt(screen, infoText, gridX, gridY-35)
}

func (r *Renderer) DrawProgressBar(screen *ebiten.Image, progress float64, episode, maxEpisodes int) {
	barX := float32(10)
	barY := float32(r.screenHeight - 50)
	barWidth := float32(r.screenWidth - 20)
	barHeight := float32(30)

	vector.FillRect(screen, barX, barY, barWidth, barHeight, color.RGBA{40, 40, 50, 255}, false)
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 2, color.RGBA{100, 100, 120, 255}, false)

	fillWidth := barWidth * float32(progress)
	if fillWidth > 0 {
		r := uint8(100 - progress*50)
		g := uint8(200 - progress*50)
		b := uint8(100 + progress*100)
		vector.FillRect(screen, barX+2, barY+2, fillWidth-4, barHeight-4, color.RGBA{r, g, b, 255}, false)
	}

	progressText := fmt.Sprintf("Training Progress: %.1f%% (%d/%d episodes)", progress*100, episode, maxEpisodes)
	ebitenutil.DebugPrintAt(screen, progressText, int(barX)+10, int(barY)+10)
}

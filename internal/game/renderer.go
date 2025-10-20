package game

import (
	"fmt"
	"image/color"

	"snakes-ml/config"
	"snakes-ml/internal/snake"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Renderer handles game rendering
type Renderer struct {
	screenWidth  int
	screenHeight int
}

// NewRenderer creates new renderer
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		screenWidth:  width,
		screenHeight: height,
	}
}

// DrawSnake renders snake and game field
func (r *Renderer) DrawSnake(screen *ebiten.Image, s *snake.Snake) {
	gridX := config.GridStartX
	gridY := config.GridStartY
	cellSize := config.CellSizeInit

	// Calculate adaptive cell size
	maxWidth := r.screenWidth - gridX*2
	maxHeight := r.screenHeight - gridY - config.GridPadding

	cellWidth := maxWidth / s.Width()
	cellHeight := maxHeight / s.Height()

	if cellWidth < cellHeight {
		cellSize = cellWidth
	} else {
		cellSize = cellHeight
	}

	if cellSize < config.CellSizeMin {
		cellSize = config.CellSizeMin
	}
	if cellSize > config.CellSizeMax {
		cellSize = config.CellSizeMax
	}

	totalWidth := s.Width() * cellSize
	// ✅ ИСПРАВЛЕНО: используем totalHeight
	totalHeight := s.Height() * cellSize
	gridX = (r.screenWidth - totalWidth) / 2

	// Draw grid
	for y := 0; y < s.Height(); y++ {
		for x := 0; x < s.Width(); x++ {
			posX := float32(gridX + x*cellSize)
			posY := float32(gridY + y*cellSize)
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 1, config.ColorGrid, false)
		}
	}

	// Draw obstacles (yellow)
	for _, obs := range s.Obstacles() {
		posX := float32(gridX + obs.X*cellSize)
		posY := float32(gridY + obs.Y*cellSize)
		vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), config.ColorObstacle, false)
		vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, config.ColorObstacleBorder, false)
	}

	// Draw food (red)
	food := s.Food()
	posX := float32(gridX + food.X*cellSize)
	posY := float32(gridY + food.Y*cellSize)
	vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), config.ColorFood, false)
	vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, config.ColorFoodBorder, false)

	// Draw snake
	for i, segment := range s.Body() {
		posX := float32(gridX + segment.X*cellSize)
		posY := float32(gridY + segment.Y*cellSize)

		var col color.RGBA
		if i == 0 {
			col = config.ColorSnakeHead
		} else {
			intensity := uint8(200 - i*2)
			// ✅ ИСПРАВЛЕНО: правильное сравнение uint8 с uint8
			if intensity < config.ColorSnakeBodyMin {
				intensity = config.ColorSnakeBodyMin
			}
			col = color.RGBA{50, intensity, 50, 255}
		}

		vector.FillRect(screen, posX, posY, float32(cellSize), float32(cellSize), col, false)

		if i == 0 {
			vector.StrokeRect(screen, posX, posY, float32(cellSize), float32(cellSize), 2, config.ColorSnakeHeadBorder, false)
		}
	}

	// Draw info
	occupancy := s.GetOccupancy() * 100
	infoText := fmt.Sprintf("Length: %d | Steps: %d | Map: %dx%d | Occupancy: %.1f%% | Obstacles: %d",
		s.Length(), s.Steps(), s.Width(), s.Height(), occupancy, len(s.Obstacles()))

	// ✅ ИСПРАВЛЕНО: используем totalHeight для правильного позиционирования
	_ = totalHeight // Помечаем как используемую переменную
	vector.FillRect(screen, float32(gridX-5), float32(gridY-40), float32(totalWidth+10), 35, config.ColorTextBg, false)
	ebitenutil.DebugPrintAt(screen, infoText, gridX, gridY-35)
}

// DrawProgressBar renders training progress bar
func (r *Renderer) DrawProgressBar(screen *ebiten.Image, progress float64, current, total int) {
	barX := float32(10)
	barY := float32(r.screenHeight - config.ProgressBarMargin)
	barWidth := float32(r.screenWidth - 20)
	barHeight := float32(config.ProgressBarHeight)

	vector.FillRect(screen, barX, barY, barWidth, barHeight, config.ColorProgressBg, false)
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 2, config.ColorProgressBorder, false)

	fillWidth := barWidth * float32(progress)
	if fillWidth > 0 {
		r := uint8(100 - progress*50)
		g := uint8(200 - progress*50)
		b := uint8(100 + progress*100)
		vector.FillRect(screen, barX+2, barY+2, fillWidth-4, barHeight-4, color.RGBA{r, g, b, 255}, false)
	}

	progressText := fmt.Sprintf("Training Progress: %.1f%% (Generation %d/%d)", progress*100, current, total)
	ebitenutil.DebugPrintAt(screen, progressText, int(barX)+10, int(barY)+10)
}

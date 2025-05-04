package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	Font             font.Face
	countAa          int
	buttons          []string
	fadeAlpha        float64 // от 0 (полная темнота) до 1 (полная видимость)
	fadeSpeed        float64 // скорость появления
	isExiting        bool    // проверка состояние выхода из игры
	heroX            int     // X героя
	heroY            int     // Y героя
	level            int     // текущий уровень
	isflip           bool    // развернут ли герой (по умолчанию нет - смотрит вправо)
	count            int     // счётчик кадров анимации
	frameWidth       int     // ширина кадра героя
	frameHeight      int     // высота кадра героя
	currentAnimation string  // текущая анимация героя
	animspeed        int     // скорость анимации
	frameCounter     int     // число кадров на анимацию
	frameOX          int
	frameOY          int
	speedY           int
	gravity          float32
}

// ====================
// Логические функции (Aa-Kk)
// ====================

func (core *Game) Aa() { // Главное меню
	buttons := 3

	if core.countAa < 0 {
		core.countAa = buttons - 1
	}
	if core.countAa >= buttons {
		core.countAa = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		core.countAa--
		if core.countAa < 0 {
			core.countAa = buttons - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		core.countAa++
		if core.countAa >= buttons {
			core.countAa = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch core.countAa {
		case 0:
			core.level = 0
		case 1:
			// Загрузка сохранения
		case 2:
			core.isExiting = true
		}
	}
}

func (core *Game) Bb() { // Границы игрового мира
	if core.heroX < 2 {
		core.heroX = 3
	}
	if core.heroX > 1280 {
		core.heroX = 1277
	}
	if core.heroY > 700 {
		core.heroY = 700
	}
}
func (core *Game) Cc() { // Управление героем
	if core.isflip {
		core.currentAnimation = "/idleR.png"
	} else if !core.isflip {
		core.currentAnimation = "/idle.png"
	}
	core.frameCounter = 2
	core.frameWidth = 18
	core.frameHeight = 36
	core.animspeed = 5
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		core.isflip = false
		core.currentAnimation = "/walk1.png"
		core.heroX += 5
		core.frameCounter = 8
		core.frameWidth = 32
		core.frameHeight = 36
		core.animspeed = 6
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		core.isflip = true
		core.currentAnimation = "/walk1R.png"
		core.heroX -= 5
		core.frameCounter = 8
		core.frameWidth = 32
		core.frameHeight = 36
		core.animspeed = 6
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		core.gravity = 0.7
		core.speedY = -15
		core.currentAnimation = "/jump.png"
		core.frameCounter = 8
		core.frameWidth = 18
		core.frameHeight = 31
		core.animspeed = 6
		core.heroY -= core.speedY
		floatHY := float32(core.heroY)
		floatHY += core.gravity
	}
	floatHY := float32(core.heroY)
	floatHY += core.gravity
}
func Dd() {

}
func Ee() {

}
func Ff() {

}
func Gg() {

}
func Hh() {

}
func Ii() {

}
func Jj() font.Face { // Загрузка шрифта
	fontBytes, err := ioutil.ReadFile("font/pixel.ttf")
	// 2. Парсим шрифт
	tt, err := opentype.Parse(fontBytes)
	// 3. Создаём Face
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face

}
func Kk() {

}

// ====================
// Графические функции (gB1-gB10)
// ====================

func (core *Game) gB1(screen *ebiten.Image) { // Графика меню
	screen.Fill(color.RGBA{0, 49, 83, 255})

	for i, btn := range core.buttons {
		y := 200 + i*200

		if i == core.countAa {
			text.Draw(screen, btn, core.Font, 100, y, color.White)
		} else {
			text.Draw(screen, btn, core.Font, 100, y, color.Gray{100})
		}
	}

	text.Draw(screen, "Nex (Alpha v1)", core.Font, 100, 100, color.White)

	overlay := ebiten.NewImage(1280, 720)
	overlay.Fill(color.RGBA{0, 0, 0, uint8(255 * (1 - core.fadeAlpha))})
	screen.DrawImage(overlay, nil)
}
func (core *Game) gB2(screen *ebiten.Image) { // Графика уровней
	screen.Fill(color.RGBA{0, 49, 83, 255})
	vector.DrawFilledRect(screen, 0, 700, 1280, 20, color.White, false)
	text.Draw(screen, "Тестирование...", core.Font, 100, 100, color.White)
}
func (core *Game) gB3(screen *ebiten.Image) { // Графика героя
	core.gB2(screen)
	var err error
	runnerImage, _, err := ebitenutil.NewImageFromFile("spr/hero" + core.currentAnimation)
	if err != nil {
		log.Fatal("Ошибка загрузки спрайта:", err)
	}
	op := &ebiten.DrawImageOptions{}
	scale := 1.0
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-float64(core.frameWidth)/2, -float64(core.frameHeight)/2)
	op.GeoM.Translate(float64(core.heroX), float64(core.heroY))
	i := (core.count / core.animspeed) % core.frameCounter
	sx, sy := core.frameOX+i*core.frameWidth, core.frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+core.frameWidth, sy+core.frameHeight)).(*ebiten.Image), op)
}
func gB4() {

}
func gB5() {

}
func gB6() {

}
func gB7() {

}
func gB8() {

}
func gB9() {

}
func gB10() {

}

func (core *Game) Update() error {
	if core.isExiting {
		// Затемнение при выходе
		core.fadeAlpha -= core.fadeSpeed
		if core.fadeAlpha <= 0 {
			os.Exit(0) // завершаем игру когда экран полностью чёрный
		}
	} else if core.fadeAlpha < 1 {
		// Появление при старте
		core.fadeAlpha += core.fadeSpeed
	}
	core.Aa()
	switch core.level {
	case 0:
		core.Cc()
		core.Bb()
	}
	core.count++
	return nil
}

func (core *Game) Draw(screen *ebiten.Image) {
	core.gB1(screen)
	switch core.level {
	case 0:
		core.gB3(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 720
}

func main() {
	game := &Game{
		heroX:     100,
		heroY:     600,
		level:     1000,
		fadeAlpha: 0,
		fadeSpeed: 0.05,
		countAa:   0,
		buttons:   []string{"Новая игра", "Продолжить", "Выход"},
		Font:      Jj(),
	}

	// Настройка окна
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Тесты")

	// Запуск игры
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

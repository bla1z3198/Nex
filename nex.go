package main

import (
	"fmt"
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
	heroX            float32 // X героя
	heroY            float32 // Y героя
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
	multiplier       int       // множитель анимации
	plt              []float32 // массив платформ
	W                float32   // ширина платформы
	H                float32   // высота платформы
	ground           bool      // на земле или нет?
	px               float32
	py               float32
}

// Логические функции (Aa-Jj)
func (core *Game) Aa() { // Главное меню (v3-final)
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
			os.Exit(0)
		}
	}
}

func (core *Game) Bb() { // Границы игрового мира (v5-final)
	if core.heroX <= 2 {
		core.heroX = 10
	}
	if core.heroX > float32(1280-core.frameWidth) {
		core.heroX = 1270 - 20
	}
	if core.heroY > float32(690-core.frameHeight) {
		core.heroY = float32(690 - core.frameHeight)
	}
	if core.heroY <= 0 {
		core.heroY = 1
	}
}
func (core *Game) Cc() { // Управление героем (v4-final)
	if core.isflip {
		core.currentAnimation = "/idleR.png"
		core.multiplier = 2
	} else if !core.isflip {
		core.currentAnimation = "/idle.png"
		core.multiplier = 2
	}
	core.frameCounter = 2
	core.frameWidth = 18
	core.frameHeight = 30
	core.animspeed = 5
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		core.isflip = false
		core.currentAnimation = "/walk.png"
		core.heroX += 7
		core.frameCounter = 4
		core.frameWidth = 32
		core.frameHeight = 30
		core.animspeed = 6
		core.multiplier = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		core.isflip = true
		core.currentAnimation = "/walkR.png"
		core.heroX -= 7
		core.frameCounter = 4
		core.frameWidth = 32
		core.frameHeight = 30
		core.animspeed = 6
		core.multiplier = 0
	}
}
func (core *Game) Dd() {
	// Инициализация основных понятий
	core.ground = true                                                   // Персонаж по сути всегда на поверхности (ставим false лишь в прыжке)
	core.plt = []float32{700.00, 220.00, 600.00, 420.00, 800.00, 570.00} // Объявляем массив координат платформ
	core.heroY += 10                                                     // Постоянная гравитации (по сути как g на земле (g ~ 10))
	// Работаем с массивом в цикле
	for i := 0; i < len(core.plt); i += 2 {
		core.px = core.plt[i]   // Подбираем каждый элемент - т.е. координата X
		core.py = core.plt[i+1] // Каждый 2-ой элемент - т.е. координата Y
		// Коллизия AABB - по пересечению фигур (к примеру кадра героя и прямоугольнкика платформы)
		if core.heroX+float32(core.frameWidth) > core.px && // Правее левого края
			core.heroX < core.px+core.W && // Левее правого края
			core.heroY+float32(core.frameHeight) >= core.py && // Ниже верха платформы
			core.heroY+float32(core.frameHeight) <= core.py+15 { // Выше низа платформы
			core.heroY = core.py - float32(core.frameHeight) - 10 // Ставим на платформу
			if core.heroY == core.py-float32(core.frameHeight)-10 && core.heroY == 670 {
				core.ground = true
			} else {
				core.ground = false
			}
			break // Выходим после первой коллизии (нет смысла прокручивать цикл бесконечно)
		}
	}
	if core.heroY == core.py-float32(core.frameHeight)-10 || core.heroY == 670 { // Ставим false когда герой в воздухе
		core.ground = true
	} else {
		core.ground = false
	}
	// Блок прыжка
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && core.ground { // Прыгаем только если стоим на чём-то
		core.heroY -= 20 // Начальный толчок
		core.count = 0
		core.ground = false // Прыгнули и оказались в воздухе
		if core.count < 100 {
			core.animspeed = 2
			core.heroY -= (200 + float32(core.count)) // Плавное замедление
			core.count++
			core.ground = false
		}
	}
}

func Jj() font.Face { // Загрузка шрифта (v1-final)
	fontBytes, err := ioutil.ReadFile("font/pixel.ttf")
	// 2. Парсим шрифт
	tt, err := opentype.Parse(fontBytes)
	// 3. Создаём Face
	face, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face

}

// Графические функции (gB1-gB5)
func (core *Game) gB1(screen *ebiten.Image) { // Графика меню (v1-final)
	screen.Fill(color.RGBA{0, 49, 83, 255})

	for i, btn := range core.buttons {
		y := 200 + i*200

		if i == core.countAa {
			text.Draw(screen, btn, core.Font, 100, y, color.White)
		} else {
			text.Draw(screen, btn, core.Font, 100, y, color.Gray{100})
		}
	}
	text.Draw(screen, "Nex (Alpha v5)", core.Font, 100, 100, color.White)
}
func (core *Game) gB2(screen *ebiten.Image) { // Графика уровней
	screen.Fill(color.RGBA{0, 49, 83, 255})
	vector.DrawFilledRect(screen, 0, 710, 1280, 20, color.White, false)
	text.Draw(screen, "WAD - Controls", core.Font, 100, 100, color.White)
}
func (core *Game) gB3(screen *ebiten.Image) { // Графика игрового мира (v5-final)
	core.gB2(screen) // Уровень
	core.gB4(screen) // Платформы
	core.gB5(screen) // Нэт Граф
	var err error
	runnerImage, _, err := ebitenutil.NewImageFromFile("spr/hero2" + core.currentAnimation)
	if err != nil {
		log.Fatal("Ошибка загрузки спрайта:", err)
	}
	op := &ebiten.DrawImageOptions{} // Герой
	scale := 1.8
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-float64(core.frameWidth)/2, -float64(core.frameHeight)/2)
	op.GeoM.Translate(float64(core.heroX), float64(core.heroY))
	i := (core.count / (core.animspeed * (core.multiplier + 1))) % core.frameCounter
	sx, sy := core.frameOX+i*core.frameWidth, core.frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+core.frameWidth, sy+core.frameHeight)).(*ebiten.Image), op)
}
func (core *Game) gB4(screen *ebiten.Image) {
	core.W = 250
	core.H = 10
	for i := 0; i < len(core.plt); i += 2 {
		if i+1 >= len(core.plt) {
			break // Если элементов не хватает для пары (X, Y)
		}
		x := core.plt[i]
		y := core.plt[i+1]
		vector.DrawFilledRect(screen, float32(x), float32(y), core.W, core.H, color.White, false)
	}
}
func (core *Game) gB5(screen *ebiten.Image) {
	fps := fmt.Sprintf("FPS: %.1f | Player: (%.2f, %.2f)",
		ebiten.ActualFPS(),
		float32(core.heroX),
		core.heroY,
	)
	text.Draw(screen, fps, core.Font, 500, 90, color.White)
}
func (core *Game) Update() error {
	ebiten.SetTPS(70)
	core.Aa() // Загружаем меню
	switch core.level {
	case 0:
		core.Cc() // Включаем управление
		core.Bb() // Устанавливаем границы мира
		core.Dd() // Подключаем прыжки и платформы
	}
	core.count++
	fmt.Println(core.ground)
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
		frameHeight: 30,
		ground:      true,
		heroX:       100,
		heroY:       670,
		level:       1000,
		countAa:     0,
		buttons:     []string{"Новая игра", "Продолжить", "Выход"},
		Font:        Jj(),
	}

	// Настройка окна
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Nex")

	// Запуск игры
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

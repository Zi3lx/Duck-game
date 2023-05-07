package mypackage

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var waitForTxt int = 0
var wiatFor2Txt int = 0
var superMove int = 0
var takeMousePos = 0
var rX, rY int

func Music() {
	//Tutorial https://github.com/faiface/beep/blob/v1.1.0/examples/tutorial/1-hello-beep/a/main.go

	f, err := os.Open("grafika/dofMP3.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func ScreenShotBg() *image.RGBA {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("grafika/screenshot.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
	return img
}

func CreateWindow() *pixelgl.Window {
	x, y := robotgo.GetScreenSize() // Full screen
	cfg := pixelgl.WindowConfig{
		Title:  "De Deadly Duck",
		Bounds: pixel.R(0, 0, float64(x), float64(y)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	return win
}

func CreateText() *text.Text {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(500, 500), basicAtlas)
	basicTxt.Color = colornames.Red

	return basicTxt
}

func MakeText(endTxt *text.Text, darthTxt *text.Text, skillTxt *text.Text) {
	fmt.Fprintln(endTxt, "You Lost!") // Text to print
	fmt.Fprintln(darthTxt, "Darth Duck has arrived!")
	fmt.Fprintln(skillTxt, "FORCE PUSH")
}

func OpenFile(fileName string) image.Image {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func DrawBG(win *pixelgl.Window, bgIMG *pixel.PictureData) {
	background := pixel.NewSprite(bgIMG, bgIMG.Bounds())
	background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func InitDuck(duckIMG image.Image) (*pixel.Sprite, pixel.Rect) {
	duckPic := pixel.PictureDataFromImage(duckIMG) // open duck image
	duckSprite := pixel.NewSprite(duckPic, duckPic.Bounds())
	duckRect := duckSprite.Picture().Bounds()
	return duckSprite, duckRect
}

func CheckIfMouseOnDuck(distance float64, endTxt *text.Text, win *pixelgl.Window, bgIMG *pixel.PictureData) {
	if distance < 40 {
		endTxt.Draw(win, pixel.IM.Scaled(endTxt.Orig, 4))
		win.Update()
		time.Sleep(time.Second * 3)
		win.SetClosed(true)
	}
}

func GetRandomNumber(speed *float64) {
	for {
		number := rand.Intn(100)
		if number < 5 && waitForTxt == 0 { // 5% chance to use killer move
			superMove = 1
			wiatFor2Txt = 1
			*speed = 75
			break
		}
		time.Sleep(time.Second * 1)
	}
	time.Sleep(time.Second * 2)
	wiatFor2Txt = 0
}

func DrawAll(win *pixelgl.Window, bgIMG *pixel.PictureData, duckSprite *pixel.Sprite, duckRect pixel.Rect, duckPos pixel.Vec, darthTxt *text.Text, skillTxt *text.Text) {

	DrawBG(win, bgIMG)

	duckSprite.Draw(win, pixel.IM.Moved(duckPos)) // Draw duck
	duckRect = duckSprite.Picture().Bounds().Moved(duckPos)

	if waitForTxt == 1 {
		darthTxt.Draw(win, pixel.IM.Scaled(darthTxt.Orig, 4)) // Draw text when darth duck arrives
	}

	if superMove == 1 {
		var x, y int

		if takeMousePos == 0 {
			takeMousePos = 1
			x, y = robotgo.GetScreenSize()
			rX = rand.Intn(x-20) + 10
			rY = rand.Intn(y-20) + 10
			fmt.Println(rX, rY)
		}
		if wiatFor2Txt == 1 {
			skillTxt.Draw(win, pixel.IM.Scaled(skillTxt.Orig, 4)) // Draw super move name
		}
		robotgo.Move(rX, rY) // Move the mouse to random position and prevent it from moveing
	}
}

func CreateGame() {
	//---------------------------- GRAFIKA ---------------------------------
	img := ScreenShotBg()
	bgIMG := pixel.PictureDataFromImage(img) // Open bg image

	duckIMG := OpenFile("grafika/duck.png")
	darthDuckIMG := OpenFile("grafika/duckDarth.png")

	win := CreateWindow()

	endTxt := CreateText()
	darthTxt := CreateText()
	skillTxt := CreateText()
	MakeText(endTxt, darthTxt, skillTxt) // Init text

	// Duck
	duckSprite, duckRect := InitDuck(duckIMG)
	duckPos := pixel.V(100, 100)

	//-------------------------- KONIEC GRAFIKI -----------------------------

	speed := 150.0
	last := time.Now() // Smoth movement variables

	//---------------------------- MAIN GAME LOOP ---------------------------
	for !win.Closed() {
		mPos := win.MousePosition()

		DrawAll(win, bgIMG, duckSprite, duckRect, duckPos, darthTxt, skillTxt)

		//Change Sprites after pressing space and play music and start calculating chances for super move
		if win.Pressed(pixelgl.KeySpace) {
			duckSprite, duckRect = InitDuck(darthDuckIMG)
			go Music()
			go GetRandomNumber(&speed)
			go func() {
				waitForTxt = 1
				time.Sleep(time.Second * 2)
				waitForTxt = 0
			}()
		}

		// Time since last frame
		dt := time.Since(last).Seconds() // Duration time
		last = time.Now()

		// Move duck function
		direction := pixel.V(mPos.X-duckPos.X, mPos.Y-duckPos.Y)
		distance := direction.Len() // Calculate distance from one point to another
		if distance > 0 {
			amount := speed * dt
			if amount > distance {
				amount = distance
			}
			direction = direction.Unit().Scaled(amount) // Direction.Unit powoduje ze wektor bedzie miał długość 1, dzieki temu zawsze bedzie kierunkiem
			duckPos = duckPos.Add(direction)            // A SCALED(amount) daje długość tego wektora
		}

		CheckIfMouseOnDuck(distance, endTxt, win, bgIMG)

		// Help print to see positions
		fmt.Printf(" %f,%f	%f,%f	%f,%f\n", duckPos.X+duckRect.W(), (duckPos.Y)+duckRect.H(), mPos.X, mPos.Y, duckPos.X, duckPos.Y)
		win.Update()
		time.Sleep(time.Millisecond * 100)
	}
}

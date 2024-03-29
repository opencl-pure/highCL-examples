package main

import (
	"bytes"
	"github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/png"
	"os"
	"runtime"
)

// Game states
const (
	Logo = iota
	Title
	GamePlay
	Ending
)

func init() {
	rl.SetCallbackFunc(main)
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.SetConfigFlags(rl.FlagVsyncHint)

	rl.InitWindow(screenWidth, screenHeight, "Android example")

	rl.InitAudioDevice()

	currentScreen := Logo

	texture := rl.LoadTexture("raylib_logo.png") // Load texture (placed on assets folder)
	fx := rl.LoadSound("coin.wav")               // Load WAV audio file (placed on assets folder)
	ambient := rl.LoadMusicStream("ambient.ogg") // Load music

	rl.PlayMusicStream(ambient)

	go doSpecific()
	framesCounter := 0 // Used to count frames

	//rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(ambient)
		if runtime.GOOS == "android" && rl.IsKeyDown(rl.KeyBack) {
			break
		}

		switch currentScreen {
		case Logo:
			framesCounter++ // Count frames

			// Wait for 4 seconds (240 frames) before jumping to Title screen
			if framesCounter > 240 {
				currentScreen = Title
			}
		case Title:
			// Press enter to change to GamePlay screen
			if rl.IsGestureDetected(rl.GestureTap) {
				rl.PlaySound(fx)
				currentScreen = GamePlay
				n++
			}
		case GamePlay:
			// Press enter to change to Ending screen
			if rl.IsGestureDetected(rl.GestureTap) {
				rl.PlaySound(fx)
				currentScreen = Ending
				n++
			}
		case Ending:
			// Press enter to return to Title screen
			if rl.IsGestureDetected(rl.GestureTap) {
				rl.PlaySound(fx)
				currentScreen = Title
				n++
			}
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		switch currentScreen {
		case Logo:
			rl.DrawText("LOGO SCREEN", 20, 20, 40, rl.LightGray)
			rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, rl.White)
			rl.DrawText("WAIT for 4 SECONDS...", 290, 400, 20, rl.Gray)
		case Title:
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Green)
			rl.DrawText("TITLE SCREEN", 20, 20, 40, rl.DarkGreen)
			rl.DrawText("TAP SCREEN to JUMP to GAMEPLAY SCREEN", 160, 220, 20, rl.DarkGreen)
		case GamePlay:
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Purple)
			rl.DrawText("GAMEPLAY SCREEN", 20, 20, 40, rl.Maroon)
			rl.DrawText("TAP SCREEN to JUMP to ENDING SCREEN", 170, 220, 20, rl.Maroon)
		case Ending:
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Blue)
			rl.DrawText("ENDING SCREEN", 20, 20, 40, rl.DarkBlue)
			rl.DrawText("TAP SCREEN to RETURN to TITLE SCREEN", 160, 220, 20, rl.DarkBlue)
		default:
			break
		}
		if len(images) >= 1 && len(textures) == 0 {
			for _, img := range images {
				textures = append(textures, rl.LoadTextureFromImage(rl.NewImageFromImage(img)))
			}
		}
		if len(textures) >= 1 {
			rl.DrawTexture(textures[n%len(textures)], 0, 0, rl.White)
		}

		rl.EndDrawing()
	}

	rl.UnloadSound(fx)            // Unload sound data
	rl.UnloadMusicStream(ambient) // Unload music stream data
	rl.CloseAudioDevice()         // Close audio device (music streaming is automatically stopped)
	rl.UnloadTexture(texture)     // Unload texture data
	rl.CloseWindow()              // Close window

	os.Exit(0)
}

var kernels = []string{
	"mandelbrot_basic",
	"mandelbrot_blue_red_black",
	"mandelbrot_pseudo_random_colors",
	"julia",
	"julia_set",
	"julia_basic",
	"sierpinski_triangle",
	"sierpinski_triangle2",
}
var images = make([]image.Image, 0, len(kernels))
var textures = make([]rl.Texture2D, 0, len(kernels))
var n = 0

func loadKernelImages() {
	p := rl.HomeDir() + "/outputs/"
	var imgs = make([]image.Image, 0, len(kernels))
	for _, kernel := range kernels {
		file, err := os.ReadFile(p + kernel + "_fractal.png")
		if err != nil {
			rl.TraceLog(3, err.Error())
			continue
		}
		decode, err := png.Decode(bytes.NewReader(file))
		if err != nil {
			rl.TraceLog(3, err.Error())
			continue
		}
		imgs = append(imgs, decode)
	}
	images = append(images, imgs...)

}

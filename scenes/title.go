package scenes

import (
	"bytes"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/SpaceCamp/assets"
	"github.com/Noofbiz/SpaceCamp/systems"
)

type TitleScene struct {
	files []string
}

func (*TitleScene) Type() string { return "Title Scene" }

func (s *TitleScene) Preload() {
	s.files = []string{
		"title/title.png",
		"title/PressStart.ttf",
		"title/parsec.mp3",
		"title/cursor.png",
		"title/move.wav",
		"ship/idle.png",
		"title/earth.png",
		"title/moon.png",
		"title/stars.png",
	}
	common.AddShader(sShader)
	common.AddShader(wShader)
	common.AddShader(starShader)
	for _, file := range s.files {
		data, err := assets.Asset(file)
		if err != nil {
			log.Fatalf("Unable to locate asset with URL: %v\n", file)
		}
		err = engo.Files.LoadReaderData(file, bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Unable to load asset with URL: %v\n At %v", file, s.Type())
		}
	}
	engo.Input.RegisterButton("up", engo.KeyW)
	engo.Input.RegisterButton("down", engo.KeyS)
	engo.Input.RegisterButton("left", engo.KeyA)
	engo.Input.RegisterButton("right", engo.KeyD)
	engo.Input.RegisterButton("A", engo.KeyJ)
	engo.Input.RegisterButton("B", engo.KeyK)
	engo.Input.RegisterButton("X", engo.KeyL)
	engo.Input.RegisterButton("Y", engo.KeySemicolon)
	engo.Input.RegisterButton("FullScreen", engo.KeyFour, engo.KeyF4)
	engo.Input.RegisterButton("Exit", engo.KeyEscape)
}

func (s *TitleScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var animatable *common.Animationable
	var notanimatable *common.NotAnimationable
	w.AddSystemInterface(&common.AnimationSystem{}, animatable, notanimatable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	w.AddSystemInterface(&common.AudioSystem{}, audioable, notaudioable)

	w.AddSystem(&systems.FullScreenSystem{})
	w.AddSystem(&systems.ExitSystem{})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	var curSys systems.CursorSystem
	curSys.ClickSoundURL = s.files[4]
	curSys.CursorURL = "title/cursor.png"
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	var selectable *systems.SceneSwitchAble
	var notselectable *systems.NotSceneSwitchAble
	w.AddSystemInterface(&systems.SceneSwitchSystem{}, selectable, notselectable)

	common.SetBackground(color.RGBA{R: 0x43, G: 0x46, B: 0x4b, A: 0xff})

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(s.files[2])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	click := audio{BasicEntity: ecs.NewBasic()}
	clickPlayer, _ := common.LoadedPlayer(s.files[4])
	click.AudioComponent = common.AudioComponent{Player: clickPlayer}
	w.AddEntity(&click)

	bg := animation{BasicEntity: ecs.NewBasic()}
	bgSS := common.NewSpritesheetWithBorderFromFile(s.files[8], 340, 180, 1, 1)
	bg.Drawable = bgSS.Drawable(0)
	bg.Scale = engo.Point{X: 2.0, Y: 2.0}
	bg.AnimationComponent = common.NewAnimationComponent(bgSS.Drawables(), 0.3)
	bg.AddDefaultAnimation(&common.Animation{Name: "twinkle", Frames: []int{0, 1}})
	w.AddEntity(&bg)

	earth := sprite{BasicEntity: ecs.NewBasic()}
	earthTex, err := common.LoadedSprite(s.files[6])
	if err != nil {
		log.Fatalf("Title Scene Setup. %v texture was not found. \nError was: %v\n", s.files[6], err)
	}
	earth.RenderComponent.Drawable = earthTex
	earth.SetZIndex(1)
	w.AddEntity(&earth)

	moon := animation{BasicEntity: ecs.NewBasic()}
	moonSS := common.NewSpritesheetWithBorderFromFile(s.files[7], 64, 64, 1, 1)
	moon.Drawable = moonSS.Drawable(0)
	moon.Position = engo.Point{X: 460, Y: 20}
	moon.Scale = engo.Point{X: 2, Y: 2}
	moon.SetZIndex(1)
	moon.AnimationComponent = common.NewAnimationComponent(moonSS.Drawables(), 0.6)
	moon.AddDefaultAnimation(&common.Animation{Name: "float", Frames: []int{0, 1, 2, 3, 4, 5, 6}})
	w.AddEntity(&moon)

	ship := animation{BasicEntity: ecs.NewBasic()}
	shipSS := common.NewSpritesheetWithBorderFromFile(s.files[5], 128, 64, 1, 1)
	ship.Drawable = shipSS.Drawable(0)
	ship.Position = engo.Point{X: 125, Y: 150}
	ship.Scale = engo.Point{X: 4, Y: 4}
	ship.SetZIndex(2)
	ship.AnimationComponent = common.NewAnimationComponent(shipSS.Drawables(), 0.3)
	ship.AddDefaultAnimation(&common.Animation{Name: "idle", Frames: []int{0, 1, 2, 3, 4, 5, 6}})
	w.AddEntity(&ship)

	selFont := &common.Font{
		Size: 20,
		FG:   color.White,
		URL:  s.files[1],
	}
	selFont.CreatePreloaded()

	titleText := sprite{BasicEntity: ecs.NewBasic()}
	titleText.Drawable = common.Text{
		Text: "Space Camp!",
		Font: selFont,
	}
	titleText.SetZIndex(2)
	titleText.SetCenter(engo.Point{X: 10, Y: 25})
	titleText.Scale = engo.Point{X: 1.5, Y: 1.5}
	w.AddEntity(&titleText)

	startText := selectionsceneswitch{BasicEntity: ecs.NewBasic()}
	startText.Drawable = common.Text{
		Text: "start",
		Font: selFont,
	}
	startText.SetZIndex(2)
	startText.SetCenter(engo.Point{X: 66, Y: 75})
	startText.Selected = true
	startText.To = "New Game Scene"
	startText.NewWorld = true
	w.AddEntity(&startText)

	optionsText := selection{BasicEntity: ecs.NewBasic()}
	optionsText.Drawable = common.Text{
		Text: "options",
		Font: selFont,
	}
	optionsText.SetZIndex(2)
	optionsText.SetCenter(engo.Point{X: 66, Y: 100})
	w.AddEntity(&optionsText)

	creditsText := selection{BasicEntity: ecs.NewBasic()}
	creditsText.Drawable = common.Text{
		Text: "credits",
		Font: selFont,
	}
	creditsText.SetZIndex(2)
	creditsText.SetCenter(engo.Point{X: 66, Y: 125})
	w.AddEntity(&creditsText)

	fsText := sprite{BasicEntity: ecs.NewBasic()}
	fsText.Drawable = common.Text{
		Text: "press F4 to enable full screen!",
		Font: selFont,
	}
	fsText.SetZIndex(2)
	fsText.SetCenter(engo.Point{X: 25, Y: 150})
	fsText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&fsText)
}

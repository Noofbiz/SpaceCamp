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

	"github.com/Noofbiz/pixelshader"
)

type TakeOffScene struct {
	files []string
}

func (s *TakeOffScene) Type() string { return "Take Off Scene" }

func (s *TakeOffScene) Preload() {
	s.files = []string{
		"title/PressStart.ttf",
		"title/cursor.png",
		"title/move.wav",
		"start/dots.png",
		"start/logs.png",
		"takeoff/bg.mp3",
		"start/log.wav",
	}
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
	engo.Input.RegisterButton("Exit", engo.KeyEscape)
}

func (s *TakeOffScene) Setup(u engo.Updater) {
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

	w.AddSystem(&systems.LogFinishedSetSceneSystem{})
	w.AddSystem(&systems.ExitSystem{})
	w.AddSystem(&common.FPSSystem{Display: true})

	w.AddSystem(&systems.CombatLogSystem{
		BackgroundURL: s.files[4],
		DotURL:        s.files[3],
		FontURL:       s.files[0],
		LineDelay:     0.3,
		LetterDelay:   0.1,
	})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	var curSys systems.CursorSystem
	curSys.ClickSoundURL = s.files[2]
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	common.SetBackground(color.RGBA{R: 0x43, G: 0x46, B: 0x4b, A: 0xff})

	selFont := &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  s.files[0],
	}
	selFont.CreatePreloaded()

	logSnd := audio{BasicEntity: ecs.NewBasic()}
	logPlayer, _ := common.LoadedPlayer(s.files[6])
	logSnd.AudioComponent = common.AudioComponent{Player: logPlayer}
	logSnd.AudioComponent.Player.SetVolume(0.25)
	w.AddEntity(&logSnd)

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(s.files[5])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	click := audio{BasicEntity: ecs.NewBasic()}
	clickPlayer, _ := common.LoadedPlayer(s.files[2])
	click.AudioComponent = common.AudioComponent{Player: clickPlayer}
	w.AddEntity(&click)

	bg := sprite{BasicEntity: ecs.NewBasic()}
	bg.Drawable = pixelshader.PixelRegion{}
	bg.SetShader(starShader)
	w.AddEntity(&bg)

	msgs := []string{
		"Your ship is cleared for takeoff!",
		"10!",
		"9!",
		"8!",
		"7, 6, 5",
		"4, 3",
		"2 1",
		"BLAST OFF!!!",
	}

	for _, msg := range msgs {
		engo.Mailbox.Dispatch(systems.CombatLogMessage{
			Msg:  msg,
			Fnt:  selFont,
			Clip: logPlayer,
		})
	}
}

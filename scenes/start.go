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

type NewGameScene struct{}

func (*NewGameScene) Type() string { return "New Game Scene" }

var NGFiles = []string{
	"title/PressStart.ttf",
	"title/cursor.png",
	"start/starting.mp3",
	"start/dots.png",
	"start/logs.png",
	"start/log.wav",
	"start/jobSelect.png",
	"start/atk.png",
	"start/def.png",
	"start/spd.png",
	"start/spc.png",
	// "start/spc_chef.png",
	// "start/spc_mechanic.png",
	// "start/spc_medic.png",
	// "start/spc_pilot.png",
	// "start/spc_security.png",
	// "start/spc_exo.png",
}

func (s *NewGameScene) Preload() {
	for _, file := range NGFiles {
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
	engo.Input.RegisterButton("FullScreen", engo.KeyFour)
	engo.Input.RegisterButton("Exit", engo.KeyEscape)
}

func (*NewGameScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	w.AddSystemInterface(&common.AudioSystem{}, audioable, notaudioable)

	w.AddSystem(&systems.CombatLogSystem{
		BackgroundURL: NGFiles[4],
		DotURL:        NGFiles[3],
		FontURL:       NGFiles[0],
		LineDelay:     1,
		LetterDelay:   0.1,
	})

	w.AddSystem(&systems.FullScreenSystem{})
	w.AddSystem(&systems.ExitSystem{})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	var curSys systems.CursorSystem
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	selFont := &common.Font{
		Size: 64,
		FG:   color.White,
		URL:  NGFiles[0],
	}
	selFont.CreatePreloaded()

	var jobselectable *systems.JobSelectAble
	var notjobselectable *systems.NotJobSelectAble
	var jobSys systems.JobSelectSystem
	jobSys.Fnt = selFont
	w.AddSystemInterface(&jobSys, jobselectable, notjobselectable)

	common.SetBackground(color.RGBA{R: 0x43, G: 0x46, B: 0x4b, A: 0xff})

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(NGFiles[2])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	logPlayer, _ := common.LoadedPlayer(NGFiles[5])

	bg := sprite{BasicEntity: ecs.NewBasic()}
	bg.Drawable = pixelshader.PixelRegion{}
	bg.SetShader(sShader)
	w.AddEntity(&bg)

	chefText := selection{BasicEntity: ecs.NewBasic()}
	chefText.Drawable = common.Text{
		Text: "Chef",
		Font: selFont,
	}
	chefText.SetZIndex(2)
	chefText.SetCenter(engo.Point{X: 130, Y: 120})
	chefText.Scale = engo.Point{X: 0.25, Y: 0.25}
	chefText.Selected = true
	w.AddEntity(&chefText)

	mechanicText := selection{BasicEntity: ecs.NewBasic()}
	mechanicText.Drawable = common.Text{
		Text: "Mecha",
		Font: selFont,
	}
	mechanicText.SetZIndex(2)
	mechanicText.SetCenter(engo.Point{X: 130, Y: 160})
	mechanicText.Scale = engo.Point{X: 0.25, Y: 0.25}
	w.AddEntity(&mechanicText)

	medicText := selection{BasicEntity: ecs.NewBasic()}
	medicText.Drawable = common.Text{
		Text: "Medic",
		Font: selFont,
	}
	medicText.SetZIndex(2)
	medicText.SetCenter(engo.Point{X: 130, Y: 200})
	medicText.Scale = engo.Point{X: 0.25, Y: 0.25}
	w.AddEntity(&medicText)

	pilotText := selection{BasicEntity: ecs.NewBasic()}
	pilotText.Drawable = common.Text{
		Text: "Pilot",
		Font: selFont,
	}
	pilotText.SetZIndex(2)
	pilotText.SetCenter(engo.Point{X: 130, Y: 240})
	pilotText.Scale = engo.Point{X: 0.25, Y: 0.25}
	w.AddEntity(&pilotText)

	securityText := selection{BasicEntity: ecs.NewBasic()}
	securityText.Drawable = common.Text{
		Text: "Defen",
		Font: selFont,
	}
	securityText.SetZIndex(2)
	securityText.SetCenter(engo.Point{X: 130, Y: 280})
	securityText.Scale = engo.Point{X: 0.25, Y: 0.25}
	w.AddEntity(&securityText)

	exoText := selection{BasicEntity: ecs.NewBasic()}
	exoText.Drawable = common.Text{
		Text: "ExoBio",
		Font: selFont,
	}
	exoText.SetZIndex(2)
	exoText.SetCenter(engo.Point{X: 130, Y: 320})
	exoText.Scale = engo.Point{X: 0.25, Y: 0.25}
	w.AddEntity(&exoText)

	msgs := []string{
		"You've been selected to lead",
		"OriginX's first mission to Mars!",
		"CONGRATULATIONS!",
		"Um...",
		"What was your job again?",
	}

	for _, msg := range msgs {
		engo.Mailbox.Dispatch(systems.CombatLogMessage{
			Msg:  msg,
			Fnt:  selFont,
			Clip: logPlayer,
		})
	}
}

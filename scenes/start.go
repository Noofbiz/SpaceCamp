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

type NewGameScene struct {
	files []string
}

func (*NewGameScene) Type() string { return "New Game Scene" }

func (s *NewGameScene) Preload() {
	s.files = []string{
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
		"start/spc_chef.png",
		"start/spc_mechanic.png",
		"start/spc_medic.png",
		"start/spc_pilot.png",
		"start/spc_security.png",
		"start/spc_exo.png",
		"start/af.png",
		"start/co.png",
		"start/dl.png",
		"start/sp.png",
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

func (s *NewGameScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

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
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	selFont := &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  s.files[0],
	}
	selFont.CreatePreloaded()

	acceptSys := systems.AcceptSystem{Fnt: selFont}
	w.AddSystem(&acceptSys)

	nameSys := systems.NameSelectSystem{Fnt: selFont, MaxLen: 12}
	w.AddSystem(&nameSys)

	var jobselectable *systems.JobSelectAble
	var notjobselectable *systems.NotJobSelectAble
	var jobSys systems.JobSelectSystem
	jobSys.Fnt = selFont
	w.AddSystemInterface(&jobSys, jobselectable, notjobselectable)

	common.SetBackground(color.RGBA{R: 0x43, G: 0x46, B: 0x4b, A: 0xff})

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(s.files[2])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	logSnd := audio{BasicEntity: ecs.NewBasic()}
	logPlayer, _ := common.LoadedPlayer(s.files[5])
	logSnd.AudioComponent = common.AudioComponent{Player: logPlayer}
	logSnd.AudioComponent.Player.SetVolume(0.25)
	w.AddEntity(&logSnd)
	jobSys.LogSnd = logPlayer
	acceptSys.LogSnd = logPlayer
	nameSys.LogSnd = logPlayer

	bg := sprite{BasicEntity: ecs.NewBasic()}
	bg.Drawable = pixelshader.PixelRegion{}
	bg.SetShader(sShader)
	w.AddEntity(&bg)

	chef := character{BasicEntity: ecs.NewBasic()}
	chef.Atk = 3
	chef.Def = 2
	chef.Spd = 4
	chef.Job = "Chef"
	chef.SpecialName = "EGG:"
	chef.SpecialURL = s.files[11]
	w.AddEntity(&chef)

	mechanic := character{BasicEntity: ecs.NewBasic()}
	mechanic.Atk = 3
	mechanic.Def = 3
	mechanic.Spd = 1
	mechanic.Job = "Mecha"
	mechanic.SpecialName = "FIX:"
	mechanic.SpecialURL = s.files[12]
	w.AddEntity(&mechanic)

	medic := character{BasicEntity: ecs.NewBasic()}
	medic.Atk = 2
	medic.Def = 4
	medic.Spd = 2
	medic.Job = "Medic"
	medic.SpecialName = "AID:"
	medic.SpecialURL = s.files[13]
	w.AddEntity(&medic)

	pilot := character{BasicEntity: ecs.NewBasic()}
	pilot.Atk = 4
	pilot.Def = 3
	pilot.Spd = 2
	pilot.Job = "Pilot"
	pilot.SpecialName = "FLY:"
	pilot.SpecialURL = s.files[14]
	w.AddEntity(&pilot)

	security := character{BasicEntity: ecs.NewBasic()}
	security.Atk = 4
	security.Def = 4
	security.Spd = 1
	security.Job = "Defen"
	security.SpecialName = "GUN:"
	security.SpecialURL = s.files[15]
	w.AddEntity(&security)

	exo := character{BasicEntity: ecs.NewBasic()}
	exo.Atk = 2
	exo.Def = 2
	exo.Spd = 3
	exo.Job = "ExoBio"
	exo.SpecialName = "SCI:"
	exo.SpecialURL = s.files[16]
	w.AddEntity(&exo)

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

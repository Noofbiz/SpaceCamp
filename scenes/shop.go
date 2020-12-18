package scenes

import (
	"bytes"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/SpaceCamp/assets"
	"github.com/Noofbiz/SpaceCamp/savedata"
	"github.com/Noofbiz/SpaceCamp/systems"
	"github.com/Noofbiz/pixelshader"
)

type ShopScene struct {
	files []string
}

func (*ShopScene) Type() string { return "Shop Scene" }

func (s *ShopScene) Preload() {
	s.files = []string{
		"title/PressStart.ttf",
		"title/cursor.png",
		"shop/bg.mp3",
		"shop/dots.png",
		"shop/log.png",
		"shop/keeper.wav",
		"shop/logo.png",
		"shop/itemSelect.png",
		"shop/keeper.png",
		"title/move.wav",
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

func (s *ShopScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	savedata.CurrentSave.S.Money = 9990

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
	curSys.ClickSoundURL = s.files[9]
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	w.AddSystem(&systems.ShopKeeperSystem{KeeperURL: s.files[8]})

	common.SetBackground(color.RGBA{R: 0x43, G: 0x46, B: 0x4b, A: 0xff})

	selFont := &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  s.files[0],
	}
	selFont.CreatePreloaded()

	var shopitemselectable *systems.ShopItemSelectAble
	var notshopitemselectable *systems.NotShopItemSelectAble
	var shopItemSelectSys systems.ShopItemSelectSystem
	shopItemSelectSys.URL = s.files[7]
	shopItemSelectSys.Fnt = selFont
	w.AddSystemInterface(&shopItemSelectSys, shopitemselectable, notshopitemselectable)

	var shoppurchaseable *systems.ShopPurchaseAble
	var notshoppurchaseable *systems.NotShopPurchaseAble
	var shopPurchaseSys systems.ShopPurchaseSystem
	shopPurchaseSys.Fnt = selFont
	w.AddSystemInterface(&shopPurchaseSys, shoppurchaseable, notshoppurchaseable)

	logSnd := audio{BasicEntity: ecs.NewBasic()}
	logPlayer, _ := common.LoadedPlayer(s.files[5])
	logSnd.AudioComponent = common.AudioComponent{Player: logPlayer}
	logSnd.AudioComponent.Player.SetVolume(0.25)
	w.AddEntity(&logSnd)
	shopPurchaseSys.Snd = logPlayer

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(s.files[2])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	click := audio{BasicEntity: ecs.NewBasic()}
	clickPlayer, _ := common.LoadedPlayer(s.files[9])
	click.AudioComponent = common.AudioComponent{Player: clickPlayer}
	w.AddEntity(&click)

	bg := sprite{BasicEntity: ecs.NewBasic()}
	bg.Drawable = pixelshader.PixelRegion{}
	bg.SetShader(wShader)
	w.AddEntity(&bg)

	logo := sprite{BasicEntity: ecs.NewBasic()}
	logo.Drawable, _ = common.LoadedSprite(s.files[6])
	logo.SetZIndex(1)
	logo.Position = engo.Point{X: 0, Y: 275}
	w.AddEntity(&logo)

	itemList := [][]string{
		[]string{
			"O2",
			"O'Hare Air",
			"Fresh, crisp\nBottled",
			"200",
		},
		[]string{
			"Noodle",
			"Insta Noodle",
			"Boiled\nFlavored\nVeg?",
			"150",
		},
		[]string{
			"FrzDry",
			"Space Rations",
			"Novelty!\nBalanced!\nDessert!",
			"200",
		},
		[]string{
			"Pie",
			"Chkn Pot Pie",
			"'Fresh' veg\nChkn!\nV Hot!",
			"300",
		},
		[]string{
			"Water",
			"Hydrogendiox",
			"It's wet!",
			"100",
		},
		[]string{
			"Soda",
			"Pepsa Cola",
			"Sugar\nAnd\nCarbon!",
			"250",
		},
		[]string{
			"Beer",
			"Root beer!",
			"Great with\nFreeze-dried\nIce-cream",
			"300",
		},
		[]string{
			"Cheap",
			"Cheap'O'Fuel",
			"Makes ships\nGoooooooo",
			"250",
		},
		[]string{
			"Rocket",
			"Rocket Fuel",
			"In space\nuse space\nfuel!",
			"400",
		},
		[]string{
			"DiLi",
			"Dilithium",
			"Crystals!\nMove space\naround you!",
			"600",
		},
		[]string{
			"Shield",
			"EM Shield",
			"Repels beams\nand\nasteroids!",
			"500",
		},
		[]string{
			"Fixit",
			"Hull Repair",
			"Fixes almost\nany breach!",
			"200",
		},
		[]string{
			"Medkit",
			"Heal Ray",
			"Instantly\nheal most\ninjuries",
			"100",
		},
		[]string{
			"Bolts",
			"ChargedBolts",
			"Works on\naliens and\nasteroids!",
			"500",
		},
		[]string{
			"3DPntr",
			"Rep rap!",
			"Creates 1\nnew item\ndaily!",
			"750",
		},
		[]string{
			"AutoP",
			"Auto Pilot",
			"Flies ship\nfor you!",
			"900",
		},
		[]string{
			"SArm",
			"Salvage Arm",
			"Pull in\nspace\nscraps!",
			"900",
		},
		[]string{
			"Leave",
			"Leave Shop",
			"This one\nexits the\nstore!",
			"0",
		},
	}
	costList := [][]int{
		[]int{200, 99},
		[]int{150, 99},
		[]int{200, 99},
		[]int{300, 99},
		[]int{100, 99},
		[]int{250, 99},
		[]int{300, 99},
		[]int{250, 99},
		[]int{400, 99},
		[]int{600, 99},
		[]int{500, 99},
		[]int{200, 99},
		[]int{100, 99},
		[]int{500, 99},
		[]int{750, 1},
		[]int{900, 1},
		[]int{900, 1},
		[]int{0, 0},
	}
	for i := 0; i < len(itemList); i++ {
		item := shopItem{BasicEntity: ecs.NewBasic()}
		item.Name = itemList[i][0]
		item.FullName = itemList[i][1]
		item.Desc = itemList[i][2]
		item.Cost = itemList[i][3]
		item.Price = costList[i][0]
		item.MaxQuantity = costList[i][1]
		w.AddEntity(&item)
	}

	msgs := []string{
		"Welcome to the SpaceForce Surplus!",
		"Good enough for NASA!",
		"What can I get you?",
	}

	for _, msg := range msgs {
		engo.Mailbox.Dispatch(systems.CombatLogMessage{
			Msg:  msg,
			Fnt:  selFont,
			Clip: logPlayer,
		})
	}
}

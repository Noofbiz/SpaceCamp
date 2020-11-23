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

var wShader = &pixelshader.PixelShader{FragShader: `
#ifdef GL_ES
#define LOWP lowp
precision highp float;
#else
#define LOWP
#endif
uniform vec2 u_resolution;  // Canvas size (width,height)
uniform vec2 u_mouse;       // mouse position in screen pixels
uniform float u_time;       // Time in seconds since load

#define PI 3.14159265359
#define PI2 6.28318530718

vec4 rocket(vec2 pos){
    vec4 col = vec4(0.0);

    // Clip (because otherwise a sine is repeated)
    if(pos.x < -0.5 || pos.x > 0.5){
        return col;
    }

    if(
      // Base parabolic shape
      pos.y + 0.02 * cos(12.0 * pos.y + 0.1) * pos.y < 0.5 - pow(3.88 * pos.x, 2.0) && pos.y > -0.1
      ||
        // Lower rectangle
       ( pos.y < 0.0 && pos.y > -0.2
            &&
                // Lower left arc
                (pos.x > -0.1 || distance(pos, vec2(-0.1,-0.1)) < 0.10)
                // Lower right arc
            &&     (pos.x < 0.1  || distance(pos, vec2(0.1,-0.1)) < 0.10)
       )
      )
    {
        // Window
        if (
            distance(pos, vec2(0.0,0.2)) < 0.05
        )
        {
            col.rgb += vec3(0.1,0.1,0.1);
            col.a = 1.0;
        }
        // Rest
        else
        {
            col.rgb += vec3(1.0,1.0,1.0);
            col.a = 1.0;
        }
    }

    else if (
        pos.y < -0.4 + 0.5 * cos(4.5 * pos.x)
        &&
        pos.y > -0.5 + 0.3 * cos(3.0 * pos.x)
    )
    {
        col.rgb += vec3(1.0,0.1,0.2);
        col.a = 1.0;
    }

    // Propeller
    else if (pos.x < 0.1 && pos.y < 0.0 && pos.x > -0.1 && pos.y > -0.3)
    {
        col.rgb += vec3(0.3,0.3,0.3) + 0.3 * cos(pos.x * 10.0 + 1.0);
        col.a = 1.0;
    }


    return col;
}

mat2 rotation(float angle){
    mat2 r = mat2(cos(angle), -sin(angle), sin(angle), cos(angle));
    return r;
}

vec4 smoke(vec2 pos){
    vec4 col = vec4(0.0);

    // Density
    float d = 0.0;

    pos.y += 0.08;

    if(pos.y > 0.0){
    	return col;
    }

    pos.x += 0.003 * cos(20.0 * pos.y + 4.0 * u_time * PI2);
    float dd = distance(pos,vec2(0.0,0.0));
    if(dd > 1.0){
    	pos *= 2.2 * pow(1.0 - dd, 2.0);
    }

    pos *= 1.9;

    d += cos(pos.x * 10.0);
	d += cos(pos.x * 20.0);
	d += cos(pos.x * 40.0);

    d += 0.3 * cos(pos.y * 6.0 + 8.0 * u_time * PI2) - 1.4;
	d += 0.3 * cos(pos.y * 50.0 + 4.0 * u_time * PI2) ;
	d += 0.3 * cos(pos.y * 10.0 + 2.0 * u_time * PI2);

    if(distance(pos.x, 0.0) < 0.05){
    	d *= 0.2 - distance(pos.x, 0.0);
    } else {
    	d *= 0.0;
    }
    if( d < 0.0){
    	d = 0.0;
    }

    float dy = distance(pos.y, 0.0);

    if(dy < 0.3){
        float fac = 1.0 / 0.3 * dy;
    	col.r += 50.0 * pow(1.0 - fac,2.0) * d;
        col.g += 10.0 * pow(1.0 - fac,4.0) * d;
        col.a += 20.0 * (1.0 - fac) * d;
    }

    col.rgb += d * 10.0;
    col.a += d;

    return col;
}


vec4 alpha_over(vec4 a, vec4 b){
	return a * a.a + (1.0 - a.a) * b;
}

void main(){
    vec2 pos = gl_FragCoord.xy / u_resolution.xy;
    pos += vec2(-2.0 * u_mouse.x / u_resolution.x, 2.0 * (u_mouse.y / u_resolution.y - 0.5));

    vec4 col = 0.4 * vec4(0.3, 0.5, 0.7, 0.0) - 0.2 * cos(u_time * 0.3 + pos.y + 0.35 * pos.x);


    vec2 rocket_pos = pos * rotation(0.5 + 0.02 * cos(u_time * PI2) + 0.02 * cos(2.0 * u_time * PI2));
    rocket_pos *= 3.9;
    col = alpha_over(rocket(rocket_pos),col);

    vec2 smoke_pos = pos * rotation(0.5);
    col = alpha_over(smoke(smoke_pos),col);

    col.a = 1.0;

    gl_FragColor = col;
}
`}

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

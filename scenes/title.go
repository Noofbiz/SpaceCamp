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

var sShader = &pixelshader.PixelShader{FragShader: `
  #ifdef GL_ES
  #define LOWP lowp
  precision mediump float;
  #else
  #define LOWP
  #endif
  uniform vec2 u_resolution;  // Canvas size (width,height)
  uniform vec2 u_mouse;       // mouse position in screen pixels
  uniform float u_time;       // Time in seconds since load

	// Star Nest by Pablo Roman Andrioli

	// This content is under the MIT License.

	#define iterations 17
	#define formuparam 0.53

	#define volsteps 20
	#define stepsize 0.1

	#define zoom   0.800
	#define tile   0.850
	#define speed  0.010

	#define brightness 0.0006
	#define darkmatter 0.400
	#define distfading 0.730
	#define saturation 0.850


	void main()
	{
		//get coords and direction
		vec2 uv=gl_FragCoord.xy/u_resolution.xy-.5;
		uv.y*=u_resolution.y/u_resolution.x;
		vec3 dir=vec3(uv*zoom,1.);
		float time=u_time*speed+.25;

		//mouse rotation
		float a1=.5+u_mouse.x/u_resolution.x*.005;
		float a2=.8+u_mouse.y/u_resolution.y*.005;
		mat2 rot1=mat2(cos(a1),sin(a1),-sin(a1),cos(a1));
		mat2 rot2=mat2(cos(a2),sin(a2),-sin(a2),cos(a2));
		dir.xz*=rot1;
		dir.xy*=rot2;
		vec3 from=vec3(1.,.5,0.5);
		from+=vec3(time*2.,time,-2.);
		from.xz*=rot1;
		from.xy*=rot2;

		//volumetric rendering
		float s=0.1,fade=1.;
		vec3 v=vec3(0.);
		for (int r=0; r<volsteps; r++) {
			vec3 p=from+s*dir*.5;
			p = abs(vec3(tile)-mod(p,vec3(tile*2.))); // tiling fold
			float pa,a=pa=0.;
			for (int i=0; i<iterations; i++) {
				p=abs(p)/dot(p,p)-formuparam; // the magic formula
				a+=abs(length(p)-pa); // absolute sum of average change
				pa=length(p);
			}
			float dm=max(0.,darkmatter-a*a*.001); //dark matter
			a*=a*a; // add contrast
			if (r>6) fade*=1.-dm; // dark matter, don't render near
			//v+=vec3(dm,dm*.5,0.);
			v+=fade;
			v+=vec3(s,s*s,s*s*s*s)*a*brightness*fade; // coloring based on distance
			fade*=distfading; // distance fading
			s+=stepsize;
		}
		v=mix(vec3(length(v)),v,saturation); //color adjust
		gl_FragColor = vec4(v*.01,1.);

	}
  `}

type TitleScene struct{}

func (*TitleScene) Type() string { return "Title Scene" }

var files = []string{
	"title/title.png",
	"title/PressStart.ttf",
	"title/parsec.mp3",
	"title/cursor.png",
}

func (s *TitleScene) Preload() {
	common.AddShader(sShader)
	for _, file := range files {
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

func (*TitleScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	w.AddSystemInterface(&common.AudioSystem{}, audioable, notaudioable)

	w.AddSystem(&systems.FullScreenSystem{})
	w.AddSystem(&systems.ExitSystem{})

	var cursorable *systems.CursorAble
	var notcursorable *systems.NotCursorAble
	var curSys systems.CursorSystem
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	var sceneswitchable *systems.SceneSwitchAble
	var notsceneswitchable *systems.NotSceneSwitchAble
	w.AddSystemInterface(&systems.SceneSwitchSystem{}, sceneswitchable, notsceneswitchable)

	common.SetBackground(color.RGBA{R: 0xd8, G: 0xee, B: 0xff, A: 0xff})

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer(files[2])
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	bg := sprite{BasicEntity: ecs.NewBasic()}
	bg.Drawable = pixelshader.PixelRegion{}
	bg.SetShader(sShader)
	w.AddEntity(&bg)

	scene := sprite{BasicEntity: ecs.NewBasic()}
	sceneTex, err := common.LoadedSprite(files[0])
	if err != nil {
		log.Fatalf("Title Scene Setup. %v texture was not found. \nError was: %v\n", files[0], err)
	}
	scene.RenderComponent.Drawable = sceneTex
	w.AddEntity(&scene)

	selFont := &common.Font{
		Size: 20,
		FG:   color.White,
		URL:  files[1],
	}
	selFont.CreatePreloaded()

	titleText := sprite{BasicEntity: ecs.NewBasic()}
	titleText.Drawable = common.Text{
		Text: "Space Camp!",
		Font: selFont,
	}
	titleText.SetZIndex(1)
	titleText.SetCenter(engo.Point{X: 10, Y: 25})
	titleText.Scale = engo.Point{X: 1.5, Y: 1.5}
	w.AddEntity(&titleText)

	startText := selectionsceneswitch{BasicEntity: ecs.NewBasic()}
	startText.Drawable = common.Text{
		Text: "start",
		Font: selFont,
	}
	startText.SetZIndex(1)
	startText.SetCenter(engo.Point{X: 66, Y: 75})
	startText.Selected = true
	startText.To = "intro battle"
	startText.NewWorld = true
	w.AddEntity(&startText)

	optionsText := selection{BasicEntity: ecs.NewBasic()}
	optionsText.Drawable = common.Text{
		Text: "options",
		Font: selFont,
	}
	optionsText.SetZIndex(1)
	optionsText.SetCenter(engo.Point{X: 66, Y: 100})
	w.AddEntity(&optionsText)

	creditsText := selection{BasicEntity: ecs.NewBasic()}
	creditsText.Drawable = common.Text{
		Text: "credits",
		Font: selFont,
	}
	creditsText.SetZIndex(1)
	creditsText.SetCenter(engo.Point{X: 66, Y: 125})
	w.AddEntity(&creditsText)

	fsText := sprite{BasicEntity: ecs.NewBasic()}
	fsText.Drawable = common.Text{
		Text: "press 4 to enable full screen!",
		Font: selFont,
	}
	fsText.SetZIndex(1)
	fsText.SetCenter(engo.Point{X: 25, Y: 150})
	fsText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&fsText)
}

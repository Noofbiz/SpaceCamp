package scenes

import "github.com/Noofbiz/pixelshader"

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

	#define volsteps 10
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

var starShader = &pixelshader.PixelShader{FragShader: `
#ifdef GL_ES
precision mediump float;
#endif

#extension GL_OES_standard_derivatives : enable

uniform vec2 u_resolution;  // Canvas size (width,height)
uniform vec2 u_mouse;       // mouse position in screen pixels
uniform float u_time;       // Time in seconds since load

//////////////////////////////////////////////////
// Xavier Benech
// Galaxy Trip
// Inspired by "Star Tunnel" shader from P_Malin
// https://www.shadertoy.com/view/MdlXWr
// License Creative Commons Attribution-NonCommercial-ShareAlike 3.0 Unported License.
//

// Increase pass count for a denser effect
#define PASS_COUNT 4

float fBrightness = 2.5;

// Number of angular segments
float fSteps = 121.0;

float fParticleSize = 0.015;
float fParticleLength = 0.5 / 60.0;

// Min and Max star position radius. Min must be present to prevent stars too near camera
float fMinDist = 0.8;
float fMaxDist = 5.0;

float fRepeatMin = 1.0;
float fRepeatMax = 2.0;

// fog density
float fDepthFade = 0.8;

float Random(float x)
{
	return fract(sin(x * 123.456) * 23.4567 + sin(x * 345.678) * 45.6789 + sin(x * 456.789) * 56.789);
}

vec3 GetParticleColour( const in vec3 vParticlePos, const in float fParticleSize, const in vec3 vRayDir )
{
	vec2 vNormDir = normalize(vRayDir.xy);
	float d1 = dot(vParticlePos.xy, vNormDir.xy) / length(vRayDir.xy);
	vec3 vClosest2d = vRayDir * d1;

	vec3 vClampedPos = vParticlePos;

	vClampedPos.z = clamp(vClosest2d.z, vParticlePos.z - fParticleLength, vParticlePos.z + fParticleLength);

	float d = dot(vClampedPos, vRayDir);

	vec3 vClosestPos = vRayDir * d;

	vec3 vDeltaPos = vClampedPos - vClosestPos;

	float fClosestDist = length(vDeltaPos) / fParticleSize;
	float fShade = clamp(1.0 - fClosestDist, 0.0, 1.0);

	if (d<3.0)
	{
		fClosestDist = max(abs(vDeltaPos.x),abs(vDeltaPos.y)) / fParticleSize;
		float f = clamp(1.0 - 0.8*fClosestDist, 0.0, 1.0);
		fShade += f*f*f*f;
		fShade *= fShade;
	}

	fShade = fShade * exp2(-d * fDepthFade) * fBrightness;
	return vec3(fShade);
}

vec3 GetParticlePos( const in vec3 vRayDir, const in float fZPos, const in float fSeed )
{
	float fAngle = atan(vRayDir.x, vRayDir.y);
	float fAngleFraction = fract(fAngle / (3.14 * 2.0));

	float fSegment = floor(fAngleFraction * fSteps + fSeed) + 0.5 - fSeed;
	float fParticleAngle = fSegment / fSteps * (3.14 * 2.0);

	float fSegmentPos = fSegment / fSteps;
	float fRadius = fMinDist + Random(fSegmentPos + fSeed) * (fMaxDist - fMinDist);

	float tunnelZ = vRayDir.z / length(vRayDir.xy / fRadius);

	tunnelZ += fZPos;

	float fRepeat = fRepeatMin + Random(fSegmentPos + 0.1 + fSeed) * (fRepeatMax - fRepeatMin);

	float fParticleZ = (ceil(tunnelZ / fRepeat) - 0.5) * fRepeat - fZPos;

	return vec3( sin(fParticleAngle) * fRadius, cos(fParticleAngle) * fRadius, fParticleZ );
}

vec3 Starfield( const in vec3 vRayDir, const in float fZPos, const in float fSeed )
{
	vec3 vParticlePos = GetParticlePos(vRayDir, fZPos, fSeed);

	return GetParticleColour(vParticlePos, fParticleSize, vRayDir);
}

vec3 RotateX( const in vec3 vPos, const in float fAngle )
{
    float s = sin(fAngle); float c = cos(fAngle);
    return vec3( vPos.x, c * vPos.y + s * vPos.z, -s * vPos.y + c * vPos.z);
}

vec3 RotateY( const in vec3 vPos, const in float fAngle )
{
    float s = sin(fAngle); float c = cos(fAngle);
    return vec3( c * vPos.x + s * vPos.z, vPos.y, -s * vPos.x + c * vPos.z);
}

vec3 RotateZ( const in vec3 vPos, const in float fAngle )
{
    float s = sin(fAngle); float c = cos(fAngle);
    return vec3( c * vPos.x + s * vPos.y, -s * vPos.x + c * vPos.y, vPos.z);
}

// Simplex Noise by IQ
vec2 hash( vec2 p )
{
	p = vec2( dot(p,vec2(127.1,311.7)),
			  dot(p,vec2(269.5,183.3)) );

	return -1.0 + 2.0*fract(sin(p)*43758.5453123);
}

float noise( in vec2 p )
{
    const float K1 = 0.366025404; // (sqrt(3)-1)/2;
    const float K2 = 0.211324865; // (3-sqrt(3))/6;

	vec2 i = floor( p + (p.x+p.y)*K1 );

    vec2 a = p - i + (i.x+i.y)*K2;
    vec2 o = (a.x>a.y) ? vec2(1.0,0.0) : vec2(0.0,1.0); //vec2 of = 0.5 + 0.5*vec2(sign(a.x-a.y), sign(a.y-a.x));
    vec2 b = a - o + K2;
	   vec2 c = a - 1.0 + 2.0*K2;

    vec3 h = max( 0.5-vec3(dot(a,a), dot(b,b), dot(c,c) ), 0.0 );

	vec3 n = h*h*h*h*vec3( dot(a,hash(i+0.0)), dot(b,hash(i+o)), dot(c,hash(i+1.0)));

    return dot( n, vec3(70.0) );

}

const mat2 m = mat2( 0.80,  0.60, -0.60,  0.80 );

float fbm4( in vec2 p )
{
    float f = 0.0;
    f += 0.5000*noise( p ); p = m*p*2.02;
    f += 0.2500*noise( p ); p = m*p*2.03;
    f += 0.1250*noise( p ); p = m*p*2.01;
    f += 0.0625*noise( p );
    return f;
}

float marble(in vec2 p)
{
	return cos(p.x+fbm4(p));
}

float dowarp ( in vec2 q, out vec2 a, out vec2 b )
{
	float ang=0.;
	ang = 1.2345 * sin (33.33); //0.015*iTime);
	mat2 m1 = mat2(cos(ang), -sin(ang), sin(ang), cos(ang));
	ang = 0.2345 * sin (66.66); //0.021*iTime);
	mat2 m2 = mat2(cos(ang), -sin(ang), sin(ang), cos(ang));

	a = vec2( marble(m1*q), marble(m2*q+vec2(1.12,0.654)) );

	ang = 0.543 * cos (13.33); //0.011*iTime);
	m1 = mat2(cos(ang), -sin(ang), sin(ang), cos(ang));
	ang = 1.128 * cos (53.33); //0.018*iTime);
	m2 = mat2(cos(ang), -sin(ang), sin(ang), cos(ang));

	b = vec2( marble( m2*(q + a)), marble( m1*(q + a) ) );

	return marble( q + b +vec2(0.32,1.654));
}

// -----------------------------------------------

void main()
{
	vec4 fragCoord = gl_FragCoord;
	vec2 iResolution = u_resolution;
 	vec3 iMouse = vec3(u_mouse, 0.0);
	float iTime = u_time;

	vec2 uv = fragCoord.xy / iResolution.xy;
	vec2 q = 2.*uv-1.;
	q.y *= iResolution.y/iResolution.x;

	// camera
	vec3 rd = normalize(vec3( q.x, q.y, 1. ));
	vec3 euler = vec3(
		sin(iTime * 0.2) * 0.625,
		cos(iTime * 0.1) * 0.625,
		iTime * 0.1 + sin(iTime * 0.3) * 0.5);

	if(iMouse.z > 0.0)
	{
		euler.x = -((iMouse.y / iResolution.y) * 2.0 - 1.0);
		euler.y = -((iMouse.x / iResolution.x) * 2.0 - 1.0);
		euler.z = 0.0;
	}
	rd = RotateX(rd, euler.x);
	rd = RotateY(rd, euler.y);
	rd = RotateZ(rd, euler.z);

	// Nebulae Background
	float pi = 3.141592654;
	q.x = 0.5 + atan(rd.z, rd.x)/(2.*pi);
	q.y = 0.5 - asin(rd.y)/pi + 0.512 + 0.001*iTime;
	q *= 2.34;

	vec2 wa = vec2(0.);
	vec2 wb = vec2(0.);
	float f = dowarp(q, wa, wb);
	f = 0.5+0.5*f;

	vec3 col = vec3(f);
	float wc = 0.;
	wc = f;
	col = vec3(wc, wc*wc, wc*wc*wc);
	wc = abs(wa.x);
	col -= vec3(wc*wc, wc, wc*wc*wc);
	wc = abs(wb.x);
	col += vec3(wc*wc*wc, wc*wc, wc);
	col *= 0.7;
	col.x = pow(col.x, 2.18);
	col.z = pow(col.z, 1.88);
	col = smoothstep(0., 1., col);
	col = 0.5 - (1.4*col-0.7)*(1.4*col-0.7);
	col = 0.75*sqrt(col);
	col *= 1. - 0.5*fbm4(8.*q);
	col = clamp(col, 0., 1.);

	// StarField
	float fShade = 0.0;
	float a = 0.2;
	float b = 10.0;
	float c = 1.0;
	float fZPos = 5.0;// + iTime * c + sin(iTime * a) * b;
	float fSpeed = 0.; //c + a * b * cos(a * iTime);

	fParticleLength = 0.25 * fSpeed / 60.0;

	float fSeed = 0.0;

	vec3 vResult = vec3(0.);

	vec3 red = vec3(0.7,0.4,0.3);
	vec3 blue = vec3(0.3,0.4,0.7);
	vec3 tint = vec3(0.);
	float ti = 1./float(PASS_COUNT-1);
	float t = 0.;
	for(int i=0; i<PASS_COUNT; i++)
	{
		tint = mix(red,blue,t);
		vResult += 1.1*tint*Starfield(rd, fZPos, fSeed);
		t += ti;
		fSeed += 1.234;
		rd = RotateX(rd, 0.25*euler.x);
	}

	col += sqrt(vResult);

	// Vignetting
	vec2 r = -1.0 + 2.0*(uv);
	float vb = max(abs(r.x), abs(r.y));
	col *= (0.15 + 0.85*(1.0-exp(-(1.0-vb)*30.0)));
	gl_FragColor = vec4( col, 1.0 );
}
`}

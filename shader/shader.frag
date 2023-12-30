#version 410

in vec2 texturePos;
in float colorMult;

uniform vec4 ourColor;
uniform sampler2D tex;

out vec4 color;
void main() {
    color = mix(texture(tex,texturePos), ourColor, colorMult*0.01);
}
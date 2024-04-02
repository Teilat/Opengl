#version 410

in vec2 texturePos;
in float colorMult; // коофицент отдаления от камеры

uniform sampler2D tex;

out vec4 color;
void main() {
    color = texture2D(tex,texturePos);
    if (color.w < 1) discard; // if alpha < 1 skip

    if (texturePos == vec2(0,0)) {
        color = vec4(1,1,1,1)*(colorMult*0.06);
    }
}
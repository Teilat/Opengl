#version 410

in vec2 texturePos;
in vec3 verColor;

uniform vec4 ourColor;
uniform sampler2D ourTexture;

out vec4 color;
void main() {
    //color = ourColor;
    //color = vec4(verColor,1);
    color = texture(ourTexture,texturePos);
}
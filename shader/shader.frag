#version 410

in vec2 texturePos;

uniform vec4 ourColor;
uniform sampler2D ourTexture;

out vec4 color;
void main() {
    //color = ourColor;
    color = texture(ourTexture,texturePos);
}
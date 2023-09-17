#version 410

uniform vec4 ourColor;

out vec4 fragColor;
void main() {
    fragColor = vec4(ourColor);
}
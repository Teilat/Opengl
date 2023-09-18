#version 410

uniform vec4 ourColor;
in vec3 color;

out vec4 fragColor;
void main() {
    fragColor = vec4(color,1);
}
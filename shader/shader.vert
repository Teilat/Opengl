#version 410

layout(location = 0) in vec2 vertexPosition;

void main() {
	gl_Position = vec4(vertexPosition.x, vertexPosition.y, 0, 0.8);
}
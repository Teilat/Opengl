#version 410

layout(location = 0) in vec2 vertexPosition;
layout(location = 1) in vec3 vertexColor;

out vec3 color;

void main() {
	color = vertexColor;
	gl_Position = vec4(vertexPosition.x, vertexPosition.y, 0, 0.8);
}
#version 410

layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec3 vertexColor;
layout(location = 2) in vec2 texturePosition;

out vec2 texturePos;
out vec3 verColor;
void main() {
	texturePos = texturePosition;
	verColor = vertexColor;
	gl_Position = vec4(vertexPosition.x, vertexPosition.y, vertexPosition.z, 1);
}

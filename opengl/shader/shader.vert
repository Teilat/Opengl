#version 410

layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 texturePosition;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

out vec2 texturePos;
out float colorMult;
void main() {
	texturePos = texturePosition;
	gl_Position = projection * camera * model * vec4(vertexPosition, 1);
	colorMult = distance(vec3(camera), vec3(gl_Position));
}
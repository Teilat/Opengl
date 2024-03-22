#version 410

layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 texturePosition;

uniform mat4 projection;

uniform mat4 cameraMatrix;

uniform vec3 modelTranslation;
uniform vec4 modelRotation;
uniform vec3 modelScale;

uniform mat4 modelMatrix;

out vec2 texturePos;
out float colorMult;

void main() {
	texturePos = texturePosition;
	gl_Position = projection * cameraMatrix * modelMatrix * vec4(vertexPosition, 1);
	colorMult = distance(vec3(cameraMatrix), vec3(gl_Position));

}
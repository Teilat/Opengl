#version 460

layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 texturePosition;

uniform mat4 projection;

uniform vec3 cameraLookAt;
uniform vec3 cameraPos;
uniform vec3 cameraUp;

uniform vec3 modelTranslation;
uniform vec4 modelRotation;
uniform vec3 modelScale;

uniform mat4 modelMatrix;

out vec2 texturePos;
out float colorMult;

mat4 calcCamMatrix(vec3 pos, vec3 tar, vec3 up){
	vec3 f = normalize(tar); //camera direction
	vec3 s = normalize(cross(f, up)); // camera side (right)
	vec3 u = cross(s, f); // camera up

	mat4 camMatrix = mat4(
		s.x, u.x, -f.x, 0,
		s.y, u.y, -f.y, 0,
		s.z, u.z, -f.z, 0,
		0, 0, 0, 1);

	return camMatrix * mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, -pos.x, -pos.y, -pos.z, 1);
}

void main() {
	texturePos = texturePosition;
	gl_Position = projection * calcCamMatrix(cameraPos, cameraLookAt, cameraUp) * modelMatrix * vec4(vertexPosition, 1);
	colorMult = distance(vec3(cameraMatrix), vec3(gl_Position));
}


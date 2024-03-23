#version 460

layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 texturePosition;

uniform mat4 projection;
uniform mat3 cameraData;
uniform mat4x3 modelData;

uniform mat4 modelMatrix;

out vec2 texturePos;
out float colorMult;

mat4 calcCamMatrix(mat3 data){
	vec3 pos = vec3(cameraData[0][0],cameraData[0][1],cameraData[0][2]);
	vec3 f = normalize(vec3(cameraData[1][0],cameraData[1][1],cameraData[1][2])); //camera direction
	vec3 s = normalize(cross(f, vec3(cameraData[2][0],cameraData[2][1],cameraData[2][2]))); // camera side (right)
	vec3 u = cross(s, f); // camera up

	mat4 camMatrix = mat4(
		s.x, u.x, -f.x, 0,
		s.y, u.y, -f.y, 0,
		s.z, u.z, -f.z, 0,
		0, 0, 0, 1);

	return camMatrix * mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, -pos.x, -pos.y, -pos.z, 1);
}

void main() {
	mat4 cameraMatrix = calcCamMatrix(cameraData);

	vec3 modelTranslation = vec3(modelData[0][0],modelData[1][0],modelData[2][0]);
	vec4 modelRotation = vec4(modelData[0][1],modelData[1][1],modelData[2][1],modelData[3][1]);
	vec3 modelScale = vec3(modelData[0][2],modelData[1][2],modelData[2][2]);

	gl_Position = projection * cameraMatrix * modelMatrix * vec4(vertexPosition, 1);
	texturePos = texturePosition;
	colorMult = distance(vec3(cameraMatrix), vec3(gl_Position));
}


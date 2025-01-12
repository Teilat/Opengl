#version 410

layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 texturePosition;

// uniform mat3 cameraData;
// uniform mat4x3 modelData;

uniform mat4 projectionMatrix;
uniform mat4 modelMatrix;
uniform mat4 cameraMatrix;

out vec2 texturePos;
out float colorMult;

mat4 calcCamMatrix(mat3 data){
	vec3 p = vec3(data[0][0],data[0][1],data[0][2]); //camera pos
	vec3 f = normalize(vec3(data[1][0],data[1][1],data[1][2])); //camera direction (forward)
	vec3 s = normalize(cross(f, vec3(data[2][0],data[2][1],data[2][2]))); // camera side (right)
	vec3 u = cross(s, f); // camera up

	mat4 camMatrix = mat4(
		s.x, u.x, -f.x, 0,
		s.y, u.y, -f.y, 0,
		s.z, u.z, -f.z, 0,
		0, 0, 0, 1);

	return camMatrix * mat4(
		1,	  0, 	0,	  0,
		0,	  1, 	0,	  0,
		0,	  0, 	1,	  0,
		-p.x, -p.y, -p.z, 1);
}

void main() {
	// mat4 cameraMatrix = calcCamMatrix(cameraData);
	// vec3 modelTranslation = vec3(modelData[0][0],modelData[1][0],modelData[2][0]);
	// vec4 modelRotation = vec4(modelData[0][1],modelData[1][1],modelData[2][1],modelData[3][1]);
	// vec3 modelScale = vec3(modelData[0][2],modelData[1][2],modelData[2][2]);

	gl_Position = projectionMatrix * cameraMatrix * modelMatrix * vec4(vertexPosition, 1);
	texturePos = texturePosition;
	colorMult = distance(vec3(cameraMatrix), vec3(gl_Position));
}

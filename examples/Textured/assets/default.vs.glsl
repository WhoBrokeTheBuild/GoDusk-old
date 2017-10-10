#version 330 core

uniform mat4 _MVP;

layout(location = 0) in vec3 _Vertex;
layout(location = 1) in vec3 _Normal;
layout(location = 2) in vec2 _TexCoord;

out vec4 p_Vertex;
out vec4 p_Normal;
out vec2 p_TexCoord;

void main() {
	p_Vertex = _MVP * vec4(_Vertex, 1.0);
	p_Normal = _MVP * vec4(_Normal, 1.0);
	p_TexCoord = _TexCoord;

	gl_Position = _MVP * vec4(_Vertex, 1.0);
}

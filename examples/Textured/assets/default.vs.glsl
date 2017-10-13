#version 330 core

layout(location = 0) in vec3 _Vertex;
layout(location = 1) in vec3 _Normal;
layout(location = 2) in vec2 _TexCoord;

uniform mat4 _Model;
uniform mat4 _View;
uniform mat4 _Proj;
uniform mat4 _MVP;

uniform vec3 _LightPos;
uniform vec3 _ViewPos;

out vec3 p_LightDir;
out vec3 p_ViewDir;

out vec4 p_Vertex;
out vec4 p_Normal;
out vec2 p_TexCoord;

void main() {
	p_Vertex = _Model * vec4(_Vertex, 1.0);
	p_Normal = _Model * vec4(_Normal, 1.0);
	p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);

    p_LightDir = normalize(_LightPos - p_Vertex.xyz);
    p_ViewDir = normalize(_ViewPos - p_Vertex.xyz);

	gl_Position = _MVP * vec4(_Vertex, 1.0);
}

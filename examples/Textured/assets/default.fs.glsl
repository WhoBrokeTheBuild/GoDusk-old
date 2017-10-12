#version 330 core

uniform vec3 _Ambient;
uniform vec3 _Diffuse;
uniform vec3 _Specular;

uniform sampler2D _AmbientMap;
uniform sampler2D _DiffuseMap;
uniform sampler2D _SpecularMap;
uniform sampler2D _BumpMap;

in vec4 p_Vertex;
in vec4 p_Normal;
in vec2 p_TexCoord;

out vec4 o_Color;

void main() {
    //o_Color = vec4(_Diffuse, 1);
	o_Color = texture(_DiffuseMap, p_TexCoord);
}

#version 330 core

const uint AMBIENT_MAP_FLAG  = 1u; // 00001
const uint DIFFUSE_MAP_FLAG  = 2u; // 00010
const uint SPECULAR_MAP_FLAG = 4u; // 00100
const uint BUMP_MAP_FLAG     = 8u; // 01000

uniform mat4 _Model;
uniform mat4 _View;
uniform mat4 _Proj;
uniform mat4 _MVP;

uniform vec3 _Ambient;
uniform vec3 _Diffuse;
uniform vec3 _Specular;

uniform float _Shininess;
uniform float _Dissolve;

uniform uint _MapFlags;
uniform sampler2D _AmbientMap;
uniform sampler2D _DiffuseMap;
uniform sampler2D _SpecularMap;
uniform sampler2D _BumpMap;

in vec3 p_LightDir;
in vec3 p_ViewDir;

in vec4 p_Vertex;
in vec4 p_Normal;
in vec2 p_TexCoord;

out vec4 o_Color;

void main() {
    vec4 normal = normalize(p_Normal);

    if ((_MapFlags & BUMP_MAP_FLAG) > 0u) {
        normal = _Model * (texture(_BumpMap, p_TexCoord).rgba * 2.0 - 1.0);
    }

    vec3 ambient = _Ambient;

    if ((_MapFlags & AMBIENT_MAP_FLAG) > 0u) {
        ambient = texture(_AmbientMap, p_TexCoord).rgb;
    }

    float diffuseMult = max(0.0, dot(normal.xyz, p_LightDir));
    vec3 diffuse = diffuseMult * _Diffuse;

    if ((_MapFlags & DIFFUSE_MAP_FLAG) > 0u) {
        diffuse = diffuseMult * texture(_DiffuseMap, p_TexCoord).rgb;
    }

    vec3 halfwayDir = normalize(p_LightDir + p_ViewDir);
    float specularMult = max(0.0, dot(normal.xyz, halfwayDir));
    if (_Shininess > 0) {
        specularMult = pow(specularMult, _Shininess);
    }

    vec3 specular = vec3(specularMult);

    if ((_MapFlags & SPECULAR_MAP_FLAG) > 0u) {
        specular = specularMult * texture(_SpecularMap, p_TexCoord).rgb;
    }

    o_Color = vec4(ambient + diffuse + specular, 1.0);
}

#version 330

uniform sampler2D tileset;

in vec2 texcoord0;
in vec3 normal0;

layout(location=0) out vec4 out_diffuse;
layout(location=1) out vec3 out_normal;

void main() {
    out_diffuse = vec4(texture(tileset, texcoord0).rgb, 1.0);

    vec3 pack_normal = (normal0 + 1.0) / 2.0;
    out_normal = pack_normal;
}

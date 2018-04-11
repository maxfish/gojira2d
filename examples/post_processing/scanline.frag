#version 330 core

in vec2 uv;
out vec4 color;
uniform sampler2D tex;

void main() {
    color = vec4(texture(tex, uv).rgb * mod(gl_FragCoord.y, 2), 1);
}

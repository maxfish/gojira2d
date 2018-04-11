#version 330 core

in vec2 uv;
out vec4 color;
uniform sampler2D tex;

void main() {
    float grayScale = dot(texture(tex, uv).rgb, vec3(0.299, 0.587, 0.114));
    color = vec4(grayScale, grayScale, grayScale, 1.0);
}

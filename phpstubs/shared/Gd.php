<?php

const IMG_AVIF = 0;
const IMG_BMP = 0;
const IMG_GIF = 0;
const IMG_JPEG = 0;
const IMG_JPG = 0;
const IMG_PNG = 0;
const IMG_WBMP = 0;
const IMG_WEBP = 0;

const IMG_FILTER_BRIGHTNESS = 0;
const IMG_FILTER_COLORIZE = 0;
const IMG_FILTER_CONTRAST = 0;
const IMG_FILTER_EDGEDETECT = 0;
const IMG_FILTER_EMBOSS = 0;
const IMG_FILTER_GAUSSIAN_BLUR = 0;
const IMG_FILTER_GRAYSCALE = 0;
const IMG_FILTER_MEAN_REMOVAL = 0;
const IMG_FILTER_NEGATE = 0;
const IMG_FILTER_SELECTIVE_BLUR = 0;
const IMG_FILTER_SMOOTH = 0;

const IMG_COLOR_BRUSHED = 0;
const IMG_COLOR_STYLED = 0;
const IMG_COLOR_STYLEDBRUSHED = 0;
const IMG_COLOR_TILED = 0;
const IMG_COLOR_TRANSPARENT = 0;

const IMG_CROP_DEFAULT = 0;
const IMG_CROP_SIDES = 0;
const IMG_CROP_THRESHOLD = 0;
const IMG_CROP_TRANSPARENT = 0;

const IMG_FLIP_BOTH = 0;
const IMG_FLIP_HORIZONTAL = 0;
const IMG_FLIP_VERTICAL = 0;

final class GdImage {}
final class GdFont {}

function gd_info(): array {}
function imagetypes(): int {}
function imageloadfont(string $filename): GdFont|false {}
function imagefontwidth(GdFont|int $font): int {}
function imagefontheight(GdFont|int $font): int {}

function imagecreate(int $width, int $height): GdImage|false {}
function imagecreatetruecolor(int $width, int $height): GdImage|false {}
function imagecreatefromstring(string $data): GdImage|false {}
function imagecreatefromgif(string $filename): GdImage|false {}
function imagecreatefromjpeg(string $filename): GdImage|false {}
function imagecreatefrompng(string $filename): GdImage|false {}
function imagecreatefromwebp(string $filename): GdImage|false {}
function imagecreatefromavif(string $filename): GdImage|false {}
function imagecreatefrombmp(string $filename): GdImage|false {}
function imagecreatefromwbmp(string $filename): GdImage|false {}

function imagegif(GdImage $image, $file = null): bool {}
function imagepng(GdImage $image, $file = null, int $quality = -1, int $filters = -1): bool {}
function imagejpeg(GdImage $image, $file = null, int $quality = -1): bool {}
function imagewebp(GdImage $image, $file = null, int $quality = -1): bool {}
function imageavif(GdImage $image, $file = null, int $quality = -1, int $speed = -1): bool {}
function imagebmp(GdImage $image, $file = null, bool $compressed = true): bool {}
function imagewbmp(GdImage $image, $file = null, ?int $foreground_color = null): bool {}
function imagegd(GdImage $image, ?string $file = null): bool {}
function imagegd2(GdImage $image, ?string $file = null, int $chunk_size = 128, int $mode = 0): bool {}
function imagedestroy(GdImage $image): bool {}

function imagecolorallocate(GdImage $image, int $red, int $green, int $blue): int|false {}
function imagecolorallocatealpha(GdImage $image, int $red, int $green, int $blue, int $alpha): int|false {}
function imagecolorat(GdImage $image, int $x, int $y): int|false {}
function imagecolorclosest(GdImage $image, int $red, int $green, int $blue): int {}
function imagecolorresolve(GdImage $image, int $red, int $green, int $blue): int {}
function imagecolortransparent(GdImage $image, ?int $color = null): int {}
function imagecolorsforindex(GdImage $image, int $color): array {}
function imagecolorset(GdImage $image, int $color, int $red, int $green, int $blue, int $alpha = 0): false|null {}

function imagesx(GdImage $image): int {}
function imagesy(GdImage $image): int {}
function imagesetpixel(GdImage $image, int $x, int $y, int $color): bool {}
function imageline(GdImage $image, int $x1, int $y1, int $x2, int $y2, int $color): bool {}
function imagerectangle(GdImage $image, int $x1, int $y1, int $x2, int $y2, int $color): bool {}
function imagefilledrectangle(GdImage $image, int $x1, int $y1, int $x2, int $y2, int $color): bool {}
function imagefill(GdImage $image, int $x, int $y, int $color): bool {}
function imageellipse(GdImage $image, int $center_x, int $center_y, int $width, int $height, int $color): bool {}
function imagealphablending(GdImage $image, bool $enable): bool {}
function imagesavealpha(GdImage $image, bool $enable): bool {}
function imagecopy(GdImage $dst_image, GdImage $src_image, int $dst_x, int $dst_y, int $src_x, int $src_y, int $src_width, int $src_height): bool {}
function imagecopyresized(GdImage $dst_image, GdImage $src_image, int $dst_x, int $dst_y, int $src_x, int $src_y, int $dst_width, int $dst_height, int $src_width, int $src_height): bool {}
function imagecopyresampled(GdImage $dst_image, GdImage $src_image, int $dst_x, int $dst_y, int $src_x, int $src_y, int $dst_width, int $dst_height, int $src_width, int $src_height): bool {}

function imageflip(GdImage $image, int $mode): bool {}
function imagerotate(GdImage $image, float $angle, int $background_color): GdImage|false {}
function imagescale(GdImage $image, int $width, int $height = -1, int $mode = 0): GdImage|false {}
function imagecrop(GdImage $image, array $rectangle): GdImage|false {}
function imagecropauto(GdImage $image, int $mode = 0, float $threshold = 0.5, int $color = -1): GdImage|false {}
function imagefilter(GdImage $image, int $filter, ...$args): bool {}
function imageresolution(GdImage $image, ?int $resolution_x = null, ?int $resolution_y = null): array|bool {}

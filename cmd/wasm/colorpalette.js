var bwPalette = [[255, 255, 255, 255], [0xCC, 0xCC, 0xCC, 255], [0x77, 0x77, 0x77, 255], [0, 0, 0, 255]];
var defaultPalette = [[126, 132, 22, 255], [87, 123, 70, 255], [56, 93, 73, 255], [46, 70, 61, 255]];
var customPalette = [[126, 132, 22, 255], [87, 123, 70, 255], [56, 93, 73, 255], [46, 70, 61, 255]];

var currentPalette = "default";

function mapColor(color) {
    switch (currentPalette) {
        case "default":
            return defaultPalette[color];
        case "bw":
            return bwPalette[color];
        case "custom":
            return customPalette[color];
    }
}

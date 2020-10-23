document.scanLine = 0;
let canvas = document.getElementById('screen');

const ctx = canvas.getContext('2d');
ctx.fillStyle = '#000000';

ctx.fillRect(0, 0, 160, 144);
document.image = ctx.getImageData(0, 0, 160, 144);
document.data = new ImageData(160, 144);

const worker = new Worker('worker.js');
worker.onmessage = ev => {
    switch (ev.data.type) {
        case 'console':
            console.log(ev.data.msg);
            break;
        case 'buffer':
            document.data = new ImageData(new Uint8ClampedArray(ev.data.msg), 160, 144);
            draw();
            break;
        case 'game':
            setup(ev.data.msg);
            break;
        case 'cpu':
            document.getElementById('cpu').innerText = ev.data.msg;
            break;
        case 'mem':
            document.getElementById('memory').innerText = ev.data.msg;
            break;
        case 'oam':
            document.getElementById('oam').innerText = ev.data.msg;
            break;
        case 'vram':
            document.getElementById('vram').innerText = ev.data.msg;
            break;
        case 'custom_palette':
            for (let i = 0; i < 4; i++) {
                let value = "#" + ev.data.msg[i].slice(0, 3).map(e => e.toString(16).toUpperCase()).join("");
                document.getElementById('paletteCol' + i).value = value;
            }
    }
}

document.getElementById('rom').addEventListener('change', function () {
    var reader = new FileReader();
    reader.onload = function () {
        document.rom = new Uint8Array(this.result);
        worker.postMessage({type: 'run', msg: document.rom});
    }
    reader.readAsArrayBuffer(this.files[0]);
}, false);


function startGame() {
    worker.postMessage({type: 'start'});
}

function setup(game) {
    document.getElementById('title').innerText = game.title;
    storeRom(game.title);
    updateGameList();
    document.getElementById('cartridgeType').innerText = game.cartridgeType;
    document.getElementById('sgb').innerText = game.sgb;
    document.getElementById('cgb').innerText = game.cgb;
    document.getElementById('romSize').innerText = game.romSize;
    document.getElementById('ramSize').innerText = game.ramSize;
    if (game.nonJapanese === true) {
        document.getElementById('non-japanese').innerText = 'Non Japanese';
    } else {
        document.getElementById('non-japanese').innerText = 'Japanese';
    }
}

function draw() {
    ctx.putImageData(document.data, 0, 0);
}

function mapKey(key) {
    switch (key) {
        case 'w':
            return 2
        case 'a':
            return 1
        case 's':
            return 3
        case 'd':
            return 0
        case 'Enter':
            return 6
        case ' ':
            return 7
        case 'Shift':
            return 4
        case 'Control':
            return 5

        case ',':
            return 8 // step
        case '.':
            return 9 // stop
        case '-':
            return 10 // continue

        case 'o':
            return 11 // OAM dump
        case 'l':
            return 12 // VRAM dump
    }
    return 500
}

window.addEventListener('keydown', ev => {
    let msg = mapKey(ev.key);
    if (msg === 500) {
        return;
    }
    let message = {type: 'joyp_down', msg: msg};
    worker.postMessage(message);
});
window.addEventListener('keyup', ev => {
    let msg = mapKey(ev.key);
    if (msg === 500) {
        return;
    }
    worker.postMessage({type: 'joyp_up', msg: msg})
});

function debugMemory() {
    const start = parseInt(document.getElementById('memStart').value, 16);
    const end = parseInt(document.getElementById('memEnd').value, 16);

    worker.postMessage({type: 'memRequest', msg: {start: start, end: end}});
}

document.getElementById('colorpalette').addEventListener('change', ev => {
    let palette = document.getElementById('colorpalette').value;
    document.getElementById('custom_palette_pickers').hidden = palette !== "custom";
    worker.postMessage({type: 'palette', msg: palette});
});

function getCustomPalette() {
    let palette = [];
    for (let i = 0; i < 4; i++) {
        let color = document.getElementById('paletteCol' + i).value;
        color = color.slice(1, color.length);
        console.log(color);
        let arr = [];
        for (let j = 0; j < 3; j++) {
            let part = color.slice(j * 2, (j + 1) * 2);
            console.log(part);
            arr.push(parseInt(part, 16));
        }
        arr.push(255); // push alpha val
        palette.push(arr)
    }
    console.log('getting palette', palette);
    return palette;
}

for (let i = 0; i < 4; i++) {
    document.getElementById('paletteCol' + i).addEventListener('change', ev => {
        worker.postMessage({type: 'set_custom_palette', msg: getCustomPalette()});
    });
}

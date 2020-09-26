self.importScripts('wasm_exec.js', 'wasm.js');

var rom
var buffer = new Uint8ClampedArray(160 * 144);
var imageData = new Uint8ClampedArray(160 * 144 * 4);

console.log = function (...args) {
    const msg = 'worker: ' + args.join(' ');
    self.postMessage({msg: msg, type: 'console'});
}

function mapColor(col) {
    switch (col) {
        case 0:
            return [255, 255, 255, 255]
        case 1:
            return [0xCC, 0xCC, 0xCC, 255]
        case 2:
            return [0x77, 0x77, 0x77, 255]
        case 3:
            return [0, 0, 0, 255]
    }
}

function draw() {
    for (let i = 0; i < 160 * 144; i++) {
        let [r, g, b, a] = mapColor(buffer[i]);
        imageData[i * 4] = r
        imageData[i * 4 + 1] = g
        imageData[i * 4 + 2] = b
        imageData[i * 4 + 3] = a
    }
    self.postMessage({msg: imageData, type: 'buffer'});
}

function runGame(data) {
    rom = data;
    run();
    self.postMessage({
        msg:
            {
                title: title,
                cartridgeType: cartridgeType,
                sgb: sgb,
                cgb: cgb,
                romSize: romSize,
                ramSize: ramSize,
                nonJapanese: nonJapanese,
            }, type: 'game'
    })
}

self.onmessage = ev => {
    if (ev.data.type === 'run') {
        runGame(ev.data.msg);
    } else if (ev.data.type === 'start') {
        start();
    }
}

self.importScripts('wasm_exec.js', 'wasm.js', 'colorpalette.js');

var rom
var buffer = new Uint8ClampedArray(160 * 144);
var imageData = new Uint8ClampedArray(160 * 144 * 4);

var cpu = null;
var mem = null;
var oam = null;
var vram = null;

console.log = function (...args) {
    const msg = 'worker: ' + args.join(' ');
    self.postMessage({msg: msg, type: 'console'});
}

function draw() {
    for (let i = 0; i < 160 * 144; i++) {
        let [r, g, b, a] = mapColor(buffer[i]);
        imageData[i * 4] = r
        imageData[i * 4 + 1] = g
        imageData[i * 4 + 2] = b
        imageData[i * 4 + 3] = a
    }
    const buf = imageData.buffer;
    self.postMessage({msg: buf, type: 'buffer'}, [buf]);
    imageData = new Uint8ClampedArray(160 * 144 * 4);
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

self.postMessage({type: 'custom_palette', msg: customPalette});

self.onmessage = ev => {
    new Promise(resolve => {
        console.log('running ', ev.data.type);
        switch (ev.data.type) {
            case 'run':
                runGame(ev.data.msg);
                break;
            case 'start':
                start();
                break;
            case 'joyp_down':
                if (typeof keyDown === 'function') {
                    keyDown(ev.data.msg);
                    if (oam !== null) {
                        self.postMessage({type: 'oam', msg: new TextDecoder("utf-8").decode(oam)});
                        oam = null;
                    }
                    if (vram !== null) {
                        self.postMessage({type: 'vram', msg: new TextDecoder("utf-8").decode(vram)});
                        vram = null;
                    }
                    if (cpu !== null) {
                        self.postMessage({type: 'cpu', msg: new TextDecoder("utf-8").decode(cpu)})
                        cpu = null;
                    }
                } else {
                    console.error('no keydown');
                }
                break;
            case 'joyp_up':
                if (typeof keyUp === 'function') {
                    keyUp(ev.data.msg);
                } else {
                    console.error('no keyup');
                }
                break;
            case 'memRequest':
                memoryRequest(ev.data.msg.start, ev.data.msg.end);
                if (mem !== null) {
                    self.postMessage({type: 'mem', msg: new TextDecoder("utf-8").decode(mem)});
                    mem = null;
                }
                break;
            case 'palette':
                currentPalette = ev.data.msg;
                break;
            case 'set_custom_palette':
                customPalette = ev.data.msg;
                break;
        }
    })
        .then(value => console.log('done with ', ev.data.type));
}

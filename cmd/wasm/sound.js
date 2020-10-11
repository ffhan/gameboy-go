var audioContext = new AudioContext();
var oscillator = audioContext.createOscillator();
var gain = audioContext.createGain();
var panner = audioContext.createStereoPanner();
gain.gain.setValueAtTime(0, audioContext.currentTime);
oscillator.connect(panner);
panner.connect(gain);
gain.connect(audioContext.destination);
oscillator.start();

var osc = audioContext.createOscillator();
osc.frequency.setValueAtTime(50, audioContext.currentTime);
var gain2 = audioContext.createGain();
osc.connect(gain2);
gain2.gain.setValueAtTime(0.02, audioContext.currentTime);
osc.start();
gain2.connect(audioContext.destination);

function play() {
    gain.gain.setValueAtTime(0.01, audioContext.currentTime);
}

function stop() {
    gain.gain.setValueAtTime(0, audioContext.currentTime);
}

function setType(audioType, gainValue, panValue) {
    oscillator.type = audioType;
    panner.pan.setValueAtTime(panValue, audioContext.currentTime);
    gain.gain.setValueAtTime(gainValue, audioContext.currentTime);
}

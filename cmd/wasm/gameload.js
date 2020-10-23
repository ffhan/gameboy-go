function storeRom(name) {
    db.transaction("roms", "readwrite").objectStore("roms").add(document.rom, `rom_${name}`);
}

function loadRom(name) {
    let req = db.transaction("roms", "readonly").objectStore("roms").get(`rom_${name}`);
    req.onsuccess = ev => {
        document.rom = ev.target.result;
        console.log(`loaded ${name}`);
        worker.postMessage({type: 'run', msg: document.rom});
    }
    req.onerror = console.error;
}

function updateGameList() {
    let gamesDivOld = document.getElementById('games');
    let gamesDiv = gamesDivOld.cloneNode(false);
    gamesDivOld.parentNode.replaceChild(gamesDiv, gamesDivOld);
    let request = db.transaction("roms", "readonly").objectStore("roms").getAllKeys();
    request.onerror = ev => console.error('cannot get all roms', ev);
    request.onsuccess = ev => {
        for (let i = 0; i < ev.target.result.length; i++) {
            let key = ev.target.result[i];
            if (key.startsWith('rom_')) {
                let name = key.slice(4, key.length);
                let newChild = document.createElement('p');
                newChild.innerText = name
                newChild.onclick = ev => {
                    loadRom(name);
                }
                gamesDiv.appendChild(newChild)
            }
        }
    };
}

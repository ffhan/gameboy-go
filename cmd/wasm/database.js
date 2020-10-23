// In the following line, you should include the prefixes of implementations you want to test.
window.indexedDB = window.indexedDB || window.mozIndexedDB || window.webkitIndexedDB || window.msIndexedDB;
// DON'T use "var indexedDB = ..." if you're not in a function.
// Moreover, you may need references to some window.IDB* objects:
window.IDBTransaction = window.IDBTransaction || window.webkitIDBTransaction || window.msIDBTransaction || {READ_WRITE: "readwrite"}; // This line should only be needed if it is needed to support the object's constants for older browsers
window.IDBKeyRange = window.IDBKeyRange || window.webkitIDBKeyRange || window.msIDBKeyRange;
// (Mozilla has never prefixed these objects, so we don't need window.mozIDB*)

if (!window.indexedDB) {
    console.error("Your browser doesn't support a stable version of IndexedDB. Auto-store feature will not be available.");
}

var db;
var request = indexedDB.open("gbDb");
request.onerror = function (event) {
    console.error("Please allow IndexedDB permission to enable auto-store features.");
};
request.onsuccess = function (event) {
    db = event.target.result;
    db.onerror = function (event) {
        // Generic error handler for all errors targeted at this database's
        // requests!
        console.error("Database error: " + event.target.errorCode);
    };
    updateGameList();
};
// This event is only implemented in recent browsers
request.onupgradeneeded = function(event) {
    // Save the IDBDatabase interface
    var db = event.target.result;

    // Create another object store called "names" with the autoIncrement flag set as true.
    var objStore = db.createObjectStore("roms", { autoIncrement : true });
};

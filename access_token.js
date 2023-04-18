!async function() {
    const o = new Promise((o, e) => {
            const t = window.indexedDB.open("firebaseLocalStorageDb");
            t.onsuccess = function() {
                console.log("Got database"), o(t.result)
            }, t.onerror = function() {
                console.log("Error getting database", t.error), e(t.error)
            }
        }),
        e = (await o).transaction("firebaseLocalStorage", "readonly").objectStore("firebaseLocalStorage"),
        t = new Promise((o, t) => {
            const r = e.getAll();
            r.onsuccess = function() {
                console.log("Got objects"), o(r.result)
            }, r.onerror = function() {
                console.log("Error getting objects", r.error), t(r.error)
            }
        }),
        r = await t;
    console.log(JSON.stringify(r, null, 2));
}();

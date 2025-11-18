const { contextBridge } = require('electron');

contextBridge.exposeInMainWorld('sidebar', {
    ping: () => "pong"
});
const { app, BrowserWindow } = require('electron');
const { spawn } = require('child_process');
const path = require('path');
const isDev = require('election-is-dev');

function startBackend() {
  // TODO: Replace ternaries with actual if/else statements
  const exe = process.platform === "win32" ? "sidebar_backend.exe" : "sidebar_backend";

  const backendPath = path.join(process.resourcesPath, "backend", exe);
  const devPath = path.join(__dirname, "backend", exe);

  const binary = isDev ? devPath : backendPath;

  const child = spawn(binary, [], { stdio: "pipe" });

  child.stdout.on("data", (data) => console.log(`[Go Backend]: ${data}`));
  child.stderr.on("data", (data) => console.error(`[Go Backend Error]: ${data}`));

  return child;
}

function createWindow() {
  const win = new BrowserWindow({
    width: 1100,
    height: 800,
    webPreferences: {
      preload: path.join(__dirname, 'preload.cjs'),
    },
  });

  if (isDev) {
    win.loadURL('http://localhost:5173');
  } else {
    win.loadFile(path.join(__dirname, 'dist', 'index.html'));
  }
}

app.whenReady().then(() => {
  startBackend();
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});
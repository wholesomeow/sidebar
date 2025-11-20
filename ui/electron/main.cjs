const { app, BrowserWindow } = require('electron');
const { spawn } = require('child_process');
const path = require('path');
const isDev = require('electron-is-dev');

function startBackend() {
  // TODO: Replace ternaries with actual if/else statements
  const exe = process.platform === "win32" ? "sidebar.exe" : "sidebar";
  const backendPath = path.join(process.resourcesPath, "../backend", exe);
  const devPath = path.join(__dirname, "../backend", exe);
  const binary = isDev ? devPath : backendPath;

  console.log("Binary path:", binary);
  console.log("Dev Path: ", devPath);
  const fs = require("fs");
  console.log("Binary exists?", fs.existsSync(binary));


  const child = spawn(binary, [], { stdio: "pipe" });

  child.stdout.on("data", (data) => console.log(`[Go Backend]: ${data}`));
  child.stderr.on("data", (data) => console.error(`[Go Backend Error]: ${data}`));

  return child;
}

function startViteDevServer() {
  if (!isDev) return Promise.resolve();
  
  console.log("Starting Vite dev server...");
  viteProcess = spawn('npm', ['run', 'dev'], { cwd: __dirname, shell: true, stdio: 'inherit' });

  // Wait until the dev server is ready (simple wait for localhost:5173)
  return new Promise((resolve) => {
    const net = require('net');
    const interval = setInterval(() => {
      const client = net.createConnection({ port: 5173 }, () => {
        clearInterval(interval);
        console.log("Vite dev server ready!");
        resolve();
      });
      client.on('error', () => client.destroy());
    }, 100);
  });
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
    win.loadFile(path.join(__dirname, '../dist/index.html'));
  }
}

app.whenReady().then(async () => {
  startBackend();

  if (isDev) await startViteDevServer();
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (viteProcess) viteProcess.kill();
  if (process.platform !== 'darwin') app.quit();
});
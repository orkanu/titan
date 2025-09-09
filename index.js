#!/usr/bin/env node

const os = require("os");
const path = require("path");
const { execFile } = require("child_process");

const APP_NAME = "titan";

const PLATFORM_ARCH_MAP = {
  linux: ["x64"],
  darwin: ["arm64"],
  win32: ["x64"],
};

const translateArch = (arch) => (arch === "x64" ? "amd64" : "arm64");
const translatePlatform = (platform) =>
  platform === "win32" ? "windows" : platform;
const getExtension = (platform) => (platform === "win32" ? ".exe" : "");

function getBinaryPath() {
  const platform = os.platform();
  const arch = os.arch();
  const validArchsForPlatform = PLATFORM_ARCH_MAP[platform];
  if (!validArchsForPlatform || !validArchsForPlatform.includes(arch)) {
    console.error(`Unsupported platform or architecture: ${platform}-${arch}`);
    process.exit(1);
  }

  return path.join(
    __dirname,
    "bin",
    "out",
    `${translatePlatform(platform)}-${translateArch(arch)}`,
    `${APP_NAME}${getExtension(platform)}`,
  );
}

// Extract arguments passed to the script
const args = process.argv.slice(2);

// Execute the binary
const binaryPath = getBinaryPath();

const child = execFile(binaryPath, args, (error, stdout, stderr) => {
  if (error) {
    console.error(`Error executing ${APP_NAME}: ${error.message}`);
    process.exit(error.code);
  }
  if (stdout) process.stdout.write(stdout);
  if (stderr) process.stderr.write(stderr);
});

child.on("error", (err) => {
  console.error(`Failed to start subprocess: ${err.message}`);
  process.exit(1);
});

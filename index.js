#!/usr/bin/env node

const os = require("os");
const path = require("path");
const { spawn } = require("child_process");

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

const child = spawn(binaryPath, args, { stdio: "inherit" });
child.on("exit", (code) => process.exit(code));

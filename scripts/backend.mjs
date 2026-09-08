import { spawnSync } from "node:child_process";
import { join } from "node:path";
import process from "node:process";

const backendDirectory = join("apps", "backend");
const executableName = process.platform === "win32" ? "fickServer.exe" : "fickServer";
const executablePath = join(backendDirectory, "build", executableName);

const build = spawnSync(
  "go",
  [
    "-C",
    backendDirectory,
    "build",
    "-o",
    join("build", executableName),
    "./cmd/server",
  ],
  {
    stdio: "inherit",
    shell: false,
  },
);

if (build.status !== 0) {
  process.exit(build.status ?? 1);
}

if (process.argv.includes("--run")) {
  const server = spawnSync(executablePath, [], {
    stdio: "inherit",
    shell: false,
  });

  process.exit(server.status ?? 1);
}

import { createInterface } from "node:readline";

import { render } from "./dist/server/entry-server.js";

const apiOrigin = process.argv[2];

if (!apiOrigin) {
  console.error("API origin is required");
  process.exit(1);
}

const input = createInterface({
  input: process.stdin,
  crlfDelay: Infinity,
  terminal: false,
});

for await (const line of input) {
  if (line.trim() === "") {
    continue;
  }

  try {
    const request = JSON.parse(line);

    if (typeof request.url !== "string") {
      throw new Error("Render request URL is missing");
    }

    const result = await render(request.url, apiOrigin);

    process.stdout.write(
      `${JSON.stringify({
        html: result.html,
        statusCode: result.statusCode,
        error: "",
      })}\n`,
    );
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);

    process.stdout.write(
      `${JSON.stringify({
        html: "",
        statusCode: 500,
        error: message,
      })}\n`,
    );
  }
}

import { copyFile, mkdir } from "node:fs/promises";

await mkdir(new URL("../web/static/", import.meta.url), { recursive: true });
await copyFile(
  new URL("../node_modules/htmx.org/dist/htmx.min.js", import.meta.url),
  new URL("../web/static/htmx.min.js", import.meta.url),
);

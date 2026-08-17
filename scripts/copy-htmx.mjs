import { copyFile, mkdir } from "node:fs/promises";
import { fileURLToPath } from "node:url";

// Resolve htmx.org through Node's module resolution (works under npm
// workspace hoisting: the package may live at the root node_modules or
// nested under site/node_modules — never hardcode ../node_modules).
const src = fileURLToPath(import.meta.resolve("htmx.org/dist/htmx.min.js"));

await mkdir(new URL("../site/web/static/", import.meta.url), { recursive: true });
await copyFile(src, new URL("../site/web/static/htmx.min.js", import.meta.url));
console.log(`copied htmx.min.js from ${src}`);

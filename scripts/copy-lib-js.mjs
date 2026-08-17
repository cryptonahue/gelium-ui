import { copyFile, mkdir } from "node:fs/promises";

// Copy the consumer JS (lib/js/gelium.js) into the site static dir so the
// docs site serves the SAME file an external consumer installs — dogfooding
// the lib JS contract, not a site-specific copy.
await mkdir(new URL("../site/web/static/", import.meta.url), { recursive: true });
await copyFile(
  new URL("../lib/js/gelium.js", import.meta.url),
  new URL("../site/web/static/gelium.js", import.meta.url),
);
console.log("copied lib/js/gelium.js → site/web/static/gelium.js");

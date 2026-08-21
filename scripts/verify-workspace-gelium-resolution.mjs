import assert from "node:assert/strict";
import { access, readFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import { createRequire } from "node:module";
import { fileURLToPath, pathToFileURL } from "node:url";

const root = fileURLToPath(new URL("..", import.meta.url));
const expected = JSON.parse(
  await readFile(join(root, "lib/package.json"), "utf8"),
).version;
const requireFromSite = createRequire(pathToFileURL(join(root, "site/package.json")));
const installedPackagePath = requireFromSite.resolve("gelium-ui/package.json");
const installed = JSON.parse(await readFile(installedPackagePath, "utf8"));

assert.equal(
  installed.version,
  expected,
  `site resolves gelium-ui@${installed.version}; expected workspace gelium-ui@${expected}`,
);
await access(join(dirname(installedPackagePath), "themes/theme-alden.css"));
console.log(`site resolves gelium-ui@${installed.version} with theme-alden.css`);

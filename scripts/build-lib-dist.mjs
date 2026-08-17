import { mkdir } from "node:fs/promises";
import { execSync } from "node:child_process";
import { dirname } from "node:path";
import { fileURLToPath } from "node:url";

// Build the publishable dist bundle (lib/dist/gelium.css): the lib manifest
// plus Tailwind's preflight/theme so a consumer can drop one file in without
// their own Tailwind pipeline. Theme files are NOT bundled (themes are
// consumer-owned; see lib/styles/index.css header). The dist entry lives at
// lib/styles/dist-entry.css.
const ROOT = dirname(fileURLToPath(import.meta.url)) + "/..";

await mkdir(`${ROOT}/dist`, { recursive: true });
execSync(
  `npx tailwindcss -i lib/styles/dist-entry.css -o lib/dist/gelium.css --minify`,
  { cwd: ROOT, stdio: "inherit" }
);
console.log("built lib/dist/gelium.css");

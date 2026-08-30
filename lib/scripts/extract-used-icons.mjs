// extract-used-icons.mjs — embed only catalog glyphs referenced by a consumer.
//
// Catalogs live in the gelium-ui library (node_modules), not the app binary:
//   material  @material-symbols/svg-400  (fill, names like chevron_right)
//   tabler    @tabler/icons              (stroke outline, names like chevron-right)
//
// The consumer chooses a default set. Unprefixed literals use that set.
// Prefixed literals pick a catalog regardless of the default:
//   data-gelium-icon="chevron_right"
//   data-gelium-icon="tabler:chevron-right"
//   icons.SVG("ms:settings")
//   icons.SVG("tabler-filled:home")
//
// Usage:
//   node node_modules/gelium-ui/scripts/extract-used-icons.mjs \
//     --scan templates --out internal/icons/icons.go --package icons \
//     [--set material|tabler]

import { mkdir, readFile, readdir, writeFile } from "node:fs/promises";
import { dirname, extname, join } from "node:path";
import { fileURLToPath } from "node:url";

const SCAN_EXT = new Set([".go", ".templ", ".html", ".md", ".txt", ".js"]);
const REF = /(?:data-gelium-icon=["']([^"']+)["']|icons\.SVG\(["']([^"']+)["']\)|gelium-icon:([a-z][a-z0-9_:-]*))/g;
const REF_SHAPE = /^(?:(ms|material|tabler|tabler-filled):)?([a-z][a-z0-9_-]*)$/;
const SETS = new Set(["material", "tabler"]);

function arg(flag, fallback) {
  const i = process.argv.indexOf(flag);
  if (i === -1) return fallback;
  return process.argv[i + 1];
}

const scanRoot = arg("--scan");
const outPath = arg("--out");
const pkg = arg("--package", "icons");
const defaultSet = arg("--set", "material");
if (!scanRoot || !outPath) {
  console.error("usage: extract-used-icons.mjs --scan <dir> --out <icons.go> [--package icons] [--set material|tabler]");
  process.exit(2);
}
if (!SETS.has(defaultSet)) {
  console.error(`extract-used-icons: unknown --set ${defaultSet} (want material or tabler)`);
  process.exit(2);
}

function catalogRoot(specifier) {
  return dirname(fileURLToPath(import.meta.resolve(`${specifier}/package.json`)));
}

const MATERIAL_SRC = join(catalogRoot("@material-symbols/svg-400"), "rounded");

async function walk(dir, files = []) {
  const entries = await readdir(dir, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.name === "node_modules" || entry.name === ".git" || entry.name === "dist") continue;
    const path = join(dir, entry.name);
    if (entry.isDirectory()) await walk(path, files);
    else if (SCAN_EXT.has(extname(entry.name))) files.push(path);
  }
  return files;
}

function parseRef(raw) {
  const match = raw.match(REF_SHAPE);
  if (!match) return null;
  const prefix = match[1];
  const name = match[2];
  let catalog = defaultSet;
  if (prefix === "ms" || prefix === "material") catalog = "material";
  else if (prefix === "tabler" || prefix === "tabler-filled") catalog = prefix;
  return { key: raw, catalog, name };
}

function refsIn(text) {
  const found = [];
  for (const match of text.matchAll(REF)) {
    const raw = match[1] || match[2] || match[3];
    const parsed = parseRef(raw);
    if (parsed) found.push(parsed);
  }
  return found;
}

function stripSvg(source) {
  return source
    .replace(/<svg\b[^>]*>/i, "")
    .replace(/<\/svg>\s*$/i, "")
    .trim();
}
function viewBoxOf(source) {
  const match = source.match(/viewBox="([^"]+)"/);
  return match ? match[1] : "0 0 24 24";
}

async function loadSource(catalog, name) {
  if (catalog === "material") {
    return { kind: "fill", setAttr: null, path: join(MATERIAL_SRC, `${name}.svg`), label: "Material Symbol" };
  }
  const variant = catalog === "tabler-filled" ? "filled" : "outline";
  const label = catalog === "tabler-filled" ? "Tabler filled icon" : "Tabler icon";
  const kind = catalog === "tabler" ? "stroke" : "fill";
  const setAttr = catalog === "tabler" ? "tabler" : "tabler-filled";
  let path;
  try {
    path = fileURLToPath(import.meta.resolve(`@tabler/icons/${variant}/${name}.svg`));
  } catch {
    path = "";
  }
  return { kind, setAttr, path, label };
}

function wrap(kind, setAttr, source) {
  const inner = stripSvg(source);
  const viewBox = viewBoxOf(source);
  const set = setAttr ? ` data-gelium-set="${setAttr}"` : "";
  if (kind === "stroke") {
    return `<svg class="ui-icon"${set} aria-hidden="true" focusable="false" viewBox="${viewBox}" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">${inner}</svg>`;
  }
  return `<svg class="ui-icon"${set} aria-hidden="true" focusable="false" viewBox="${viewBox}" fill="currentColor">${inner}</svg>`;
}

const files = await walk(scanRoot);
const byKey = new Map();
for (const file of files) {
  const text = await readFile(file, "utf8");
  for (const ref of refsIn(text)) {
    if (!byKey.has(ref.key)) byKey.set(ref.key, ref);
  }
}

const ordered = [...byKey.keys()].sort();
if (ordered.length === 0) {
  console.error("extract-used-icons: no gelium icon references found under", scanRoot);
  process.exit(1);
}

const rows = [];
for (const key of ordered) {
  const ref = byKey.get(key);
  const spec = await loadSource(ref.catalog, ref.name);
  let source;
  try {
    source = await readFile(spec.path, "utf8");
  } catch {
    console.error(`extract-used-icons: unknown ${spec.label} ${ref.name}`);
    process.exit(1);
  }
  rows.push(`\t"${key}": \`${wrap(spec.kind, spec.setAttr, source)}\`,`);
}

const generated = `// Code generated by gelium-ui extract-used-icons.mjs; DO NOT EDIT.
//
// Trusted glyphs referenced by this consumer. Catalogs live in gelium-ui
// (Material Symbols Apache-2.0, Tabler Icons MIT). This file embeds only used names.

package ${pkg}

import "html/template"

var iconSVGs = map[string]template.HTML{ // #nosec G203 -- trusted, internal decorative glyphs.
${rows.join("\n")}
}

// SVG returns the trusted inline glyph for name, or empty when the name was
// not extracted into this consumer embed. Pass string literals only.
func SVG(name string) template.HTML {
	return iconSVGs[name]
}
`;

await mkdir(dirname(outPath), { recursive: true });
await writeFile(outPath, generated);
console.log(`extract-used-icons: wrote ${ordered.length} glyphs to ${outPath}`);

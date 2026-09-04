#!/usr/bin/env node
import { readFile } from "node:fs/promises";

const root = new URL("..", import.meta.url);
const files = {
  package: "lib/package.json",
  sitePackage: "site/package.json",
  lock: "package-lock.json",
  version: "lib/version.go",
  readme: "README.md",
  changelog: "CHANGELOG.md",
  llms: "lib/llms.txt",
  skill: "lib/SKILL.md",
  componentSkill: "lib/skills/14-component-implementation.md",
  staticLlms: "site/web/static/llms.txt",
};

const text = async (relative) => readFile(new URL(relative, root), "utf8");
const packageJSON = JSON.parse(await text(files.package));
const sitePackageJSON = JSON.parse(await text(files.sitePackage));
const lockJSON = JSON.parse(await text(files.lock));
const version = packageJSON.version;
const checks = [
  ["package name", packageJSON.name === "gelium-ui"],
  ["site package version", sitePackageJSON.version === version],
  ["site dependency pin", sitePackageJSON.dependencies?.["gelium-ui"] === version],
  ["lock package version", lockJSON.packages?.lib?.version === version],
  ["lock site version", lockJSON.packages?.site?.version === version],
  ["lock dependency pin", lockJSON.packages?.site?.dependencies?.["gelium-ui"] === version],
  ["AssetsVersion", (await text(files.version)).includes(`AssetsVersion = "${version}"`)],
  ["README release", (await text(files.readme)).includes(`v${version}`)],
  ["CHANGELOG entry", (await text(files.changelog)).includes(`## [${version}]`)],
  ["llms version", (await text(files.llms)).includes(`- version: ${version}`)],
  ["SKILL version", (await text(files.skill)).includes(`version: ${version}`)],
  ["component skill version", (await text(files.componentSkill)).includes(`version: ${version}`)],
  ["static llms projection", (await text(files.staticLlms)).includes(`- version: ${version}`)],
];

const failed = checks.filter(([, ok]) => !ok);
for (const [name, ok] of checks) console.log(`${ok ? "OK" : "FAIL"}: ${name}`);
if (failed.length) {
  console.error(`Release check failed for ${failed.length} item(s) at version ${version}.`);
  process.exit(1);
}
console.log(`Release check passed: gelium-ui@${version}`);

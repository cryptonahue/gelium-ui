import { copyFile } from "node:fs/promises";

for (const name of ["llms.txt", "llms-ux.txt"]) {
  await copyFile(new URL(`../lib/${name}`, import.meta.url), new URL(`../site/web/static/${name}`, import.meta.url));
  console.log(`copied lib/${name} → site/web/static/${name}`);
}

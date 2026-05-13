import { compileFromFile } from "json-schema-to-typescript";
import fs from "fs";
import path from "path";
import { glob } from "glob";
import { fileURLToPath } from "url";

// 1. Define __dirname
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 2. Use absolute paths based on __dirname for reliability
const inputDir = path.resolve(__dirname, "../src/lib/schemas");
const outputDir = path.resolve(__dirname, "../src/types/generated");

async function generate() {
  try {
    if (!fs.existsSync(outputDir)) {
      fs.mkdirSync(outputDir, { recursive: true });
    }

    // 3. Modern glob returns a Promise
    const files = await glob(`${inputDir}/**/*.json`);

    for (const file of files) {
      const ts = await compileFromFile(file);
      const fileName = path.basename(file, ".json");

      // Saving as .ts is perfect—tsc will now see and build these!
      fs.writeFileSync(path.join(outputDir, `${fileName}.ts`), ts);
      console.log(`✅ Generated ${fileName}.ts`);
    }
  } catch (err) {
    console.error("❌ Error generating types:", err);
  }
}

generate();

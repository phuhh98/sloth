import { defineConfig } from "tsup";
import { execSync } from "child_process";

export default defineConfig({
  entry: ["src/index.ts"],
  format: ["cjs", "esm"],
  splitting: false,
  sourcemap: true,
  clean: true,
  dts: true,
  outExtension({ format }) {
    return {
      js: format === "cjs" ? ".js" : ".mjs",
    };
  },
  esbuildPlugins: [
    {
      name: "pre-build-gen",
      setup(build) {
        build.onStart(() => {
          console.log("🛠  Running generator before build...");
          execSync("pnpm run generate-types", { stdio: "inherit" });
        });
      },
    },
  ],
});

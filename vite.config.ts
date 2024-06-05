import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    target: "es2022",
    outDir: ".",
    emptyOutDir: false,
    lib: {
      entry: "src/index.ts",
      formats: ["es"],
      fileName: () => "gort0_js_wasm.min.js",
    },
  },
});

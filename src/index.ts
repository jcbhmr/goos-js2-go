import Uint8ArrayFromBase64 from "./Uint8ArrayFromBase64.ts";
import gunzip from "./gunzip.ts";
import Go from "./Go.ts";
import * as node_fs_stub from "./node-fs-stub.ts";
import node_process_stub from "./node-process-stub.ts";
const [node_fs, node_process] = await Promise.all([
  import(String("node:fs")).catch(() => node_fs_stub),
  import(String("node:process")).catch(() => node_process_stub),
]);

const bytes = await gunzip(Uint8ArrayFromBase64(__APP_WASM_GZ_BASE64__));

const go = Object.assign(
  new Go({
    args: node_process.argv,
    env: node_process.env,
    "node:fs": node_fs,
    "node:process": node_process,
  }),
  { _importMeta: import.meta }
);
go._import = (s, o) => import(s);
const { instance } = await WebAssembly.instantiate(bytes, go.getImportObject());
const exitCode = await go.start(instance);
node_process.exitCode = exitCode;

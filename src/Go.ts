import PromiseWithResolvers from "./PromiseWithResolvers";
import * as node_fs_stub from "./node-fs-stub.ts";
import node_process_stub from "./node-process-stub.ts";
import { encode } from "./utils.ts";

let globalImport: ((specifier: string, options: ImportCallOptions) => any) | null;
try {
  globalImport = new Function("s", "o", "return import(s)") as (
    specifier: string,
    options: ImportCallOptions
  ) => any;
} catch {
  globalImport = null;
}

export interface GoOptions {
  args?: string[] | undefined;
  env?: Record<string, string> | undefined;
  returnOnExit?: boolean | undefined;
  "node:fs"?: any | undefined;
  "node:process"?: any | undefined;
}

export default class Go {
  #readyState: "created" | "starting" | "running" | "exited" = "created";
  #args: string[];
  #env: Record<string, string>;
  #returnOnExit: boolean;
  #node_fs: any;
  #node_process: any;
  #fetch: any;
  constructor(options: GoOptions = {}) {
    this.#args = options.args ?? ["js"];
    this.#env = options.env ?? {};
    this.#returnOnExit = options.returnOnExit ?? true;
    this.#node_fs = options["node:fs"] ?? node_fs_stub;
    this.#node_process = options["node:process"] ?? node_process_stub;
  }

  getImportObject(): WebAssembly.Imports {
    return { gojs: this.gojs };
  }

  #dataViewCache: DataView | null = null;
  get #dataView() {
    if (!this.#dataViewCache) {
      const memory = this.#instance!.exports.mem as WebAssembly.Memory;
      this.#dataViewCache = new DataView(memory.buffer);
    }
    return this.#dataViewCache;
  }
  readonly gojs: WebAssembly.ModuleImports = {
    // func wasmExit(code int32)
    "runtime.wasmExit": (sp: number) => {
      const code = this.#dataView.getInt32(sp + 8, true);
      this.#readyState = "exited";
      this.#startDeferred.resolve(code);
    },
  };

  protected _makeFuncWrapper(id: number) {
    const that = this;
    return function (...args: any[]) {
      const event = { id, this: that, args } as {
        id: number;
        this: Go;
        args: any[];
        result?: any;
      };
      // that.#pendingEvent = event;
      // that.#resume();
      return event.result;
    };
  }

  async _import(specifier: string, options: ImportCallOptions) {
    if (globalImport) {
      return globalImport(specifier, options);
    } else {
      return import(specifier); // TODO: Add options.
    }
  }

  #instance: WebAssembly.Instance | null = null;
  #startDeferred: {
    promise: Promise<number>;
    resolve: (value: number) => void;
    reject: (reason: any) => void;
  } = PromiseWithResolvers();
  async start(instance: WebAssembly.Instance): Promise<number> {
    if (this.#readyState === "created") {
      this.#readyState = "starting";
    } else {
      throw new Error("Go#start() called more than once");
    }

    let offset = 4096;
    const strPtr = (str: string) => {
      const ptr = offset;
      const bytes = encode(str + "\0");
      new Uint8Array(this.#dataView.buffer, offset, bytes.length).set(bytes);
      offset += bytes.length;
      if (offset % 8 !== 0) {
        offset += 8 - (offset % 8);
      }
      return ptr;
    };
    const argc = this.#args.length;
    const argvPtrs = [];
    for (const arg of this.#args) {
      argvPtrs.push(strPtr(arg));
    }
    argvPtrs.push(0);
    for (const [key, value] of Object.entries(this.#env)) {
      argvPtrs.push(strPtr(`${key}=${value}`));
    }
    argvPtrs.push(0);
    const argv = offset;
    for (const ptr of argvPtrs) {
      this.#dataView.setUint32(offset, ptr, true);
      this.#dataView.setUint32(offset + 4, 0, true);
      offset += 8;
    }

    const wasmMinDataAddr = 4096 + 8192;
    if (offset >= wasmMinDataAddr) {
      throw new Error("Total length of args and env exceeds limit");
    }

    const run = this.#instance!.exports.run as (
      argc: number,
      argv: number
    ) => void;
    this.#readyState = "running";
    this.#startDeferred = PromiseWithResolvers();
    run(argc, argv);
    return this.#startDeferred.promise;
  }
}

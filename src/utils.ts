import {
  decodeFallback,
  encodeFallback,
} from "./fast-text-encoding/lowlevel.js";

export function createENOSYS() {
  return Object.assign(new Error("ENOSYS: Function not implemented"), {
    code: "ENOSYS",
  });
}

export function throw_(v) {
  throw v;
}

declare var TextEncoder: any | undefined;
let encoder: any = null;
export function encode(text: string): Uint8Array {
  if (globalThis.TextEncoder) {
    if (!encoder) {
      encoder = new TextEncoder();
    }
    return encoder.encode(text);
  } else {
    return encodeFallback(text);
  }
}

declare var TextDecoder: any | undefined;
let decoder: any = null;
export function decode(bytes: Uint8Array): string {
  if (globalThis.TextDecoder) {
    if (!decoder) {
      decoder = new TextDecoder();
    }
    return decoder.decode(bytes);
  } else {
    return decodeFallback(bytes);
  }
}

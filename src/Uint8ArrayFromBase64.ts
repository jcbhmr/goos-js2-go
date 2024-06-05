import { encode } from "./utils.ts";

const alphabet =
  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
const lookup = Object.fromEntries(
  Array.from(alphabet).map((a, i) => [a.charCodeAt(0), i])
);
lookup["=".charCodeAt(0)] = 0;
lookup["-".charCodeAt(0)] = 62;
lookup["_".charCodeAt(0)] = 63;

export default function Uint8ArrayFromBase64(base64: string, options: {} = {}) {
  base64 = base64.replace(/=/g, "");
  const n = base64.length;
  const rem = n % 4;
  const k = rem && rem - 1;
  const m = (n >> 2) * 3 + k;
  const encoded = encode(base64 + "===");
  for (let i = 0, j = 0; i < n; i += 4, j += 3) {
    const x =
      (lookup[encoded[i]] << 18) +
      (lookup[encoded[i + 1]] << 12) +
      (lookup[encoded[i + 2]] << 6) +
      lookup[encoded[i + 3]];
    encoded[j] = x >> 16;
    encoded[j + 1] = (x >> 8) & 0xff;
    encoded[j + 2] = x & 0xff;
  }
  return new Uint8Array(encoded.buffer, 0, m);
}

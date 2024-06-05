import { createENOSYS, decode } from "./utils.ts";

declare var console: any | undefined;

export const constants = {
  __proto__: null,
  O_WRONLY: -1,
  O_RDWR: -1,
  O_CREAT: -1,
  O_TRUNC: -1,
  O_APPEND: -1,
  O_EXCL: -1,
};
let outputBuffer = "";
export function writeSync(fd, buffer, offset, length, position) {
  if (
    (fd !== 1 && fd !== 2) ||
    offset !== 0 ||
    length !== buffer.length ||
    position !== null
  ) {
    throw createENOSYS();
  }
  outputBuffer += decode(buffer);
  const newlineIndex = outputBuffer.lastIndexOf("\n");
  if (newlineIndex !== -1) {
    const text = outputBuffer.slice(0, newlineIndex);
    outputBuffer = outputBuffer.slice(newlineIndex + 1);
    if (globalThis.console) {
      console.log(text);
    }
  }
  return buffer.length;
}
export function write(fd, buffer, offset, length, position, callback) {
  let x;
  try {
    x = this.writeSync(fd, buffer, offset, length, position);
  } catch (e) {
    callback(e);
  }
  callback(null, x);
}
export const chmod = (path, mode, callback) => void callback(createENOSYS());
export const chown = (path, uid, gid, callback) =>
  void callback(createENOSYS());
export const close = (fd, callback) => void callback(createENOSYS());
export const fchmod = (fd, mode, callback) => void callback(createENOSYS());
export const fchown = (fd, uid, gid, callback) => void callback(createENOSYS());
export const fstat = (fd, callback) => void callback(createENOSYS());
export const fsync = (fd, callback) => void callback(null);
export const ftruncate = (fd, length, callback) =>
  void callback(createENOSYS());
export const lchown = (path, uid, gid, callback) =>
  void callback(createENOSYS());
export const link = (path, link, callback) => void callback(createENOSYS());
export const lstat = (path, callback) => void callback(createENOSYS());
export const mkdir = (path, perm, callback) => void callback(createENOSYS());
export const open = (path, flags, mode, callback) =>
  void callback(createENOSYS());
export const read = (fd, buffer, offset, length, position, callback) =>
  void callback(createENOSYS());
export const readdir = (path, callback) => void callback(createENOSYS());
export const readlink = (path, callback) => void callback(createENOSYS());
export const rename = (from, to, callback) => void callback(createENOSYS());
export const rmdir = (path, callback) => void callback(createENOSYS());
export const stat = (path, callback) => void callback(createENOSYS());
export const symlink = (path, link, callback) => void callback(createENOSYS());
export const truncate = (path, length, callback) =>
  void callback(createENOSYS());
export const unlink = (path, callback) => void callback(createENOSYS());
export const utimes = (path, atime, mtime, callback) =>
  void callback(createENOSYS());

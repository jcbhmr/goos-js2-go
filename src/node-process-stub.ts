import { createENOSYS, throw_ } from "./utils.ts";

declare var console: any | undefined;

export const getuid = () => -1;
export const getgid = () => -1;
export const geteuid = () => -1;
export const getegid = () => -1;
export const pid = -1;
export const ppid = -1;
export const getgroups = () => throw_(createENOSYS());
export const umask = () => throw_(createENOSYS());
export const cwd = () => throw_(createENOSYS());
export const chdir = () => throw_(createENOSYS());
export const exit = (code: number | undefined = undefined) => {
  if (code == null) {
    if (globalThis.console) {
      if (exitCode) {
        console.warn(`exit code: ${exitCode}`)
      }
    }
  } else {
    exitCode = code;
  }
};
export let exitCode = 0;
export default {
  __proto__: null,
  getuid,
  getgid,
  geteuid,
  getegid,
  get pid() {
    return pid;
  },
  get ppid() {
    return ppid;
  },
  getgroups,
  umask,
  cwd,
  chdir,
  exit,
  get exitCode() {
    return exitCode;
  },
  set exitCode(v) {
    exitCode = v;
  }
};

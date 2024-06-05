export default function PromiseWithResolvers<T = any>() {
  let resolve: (value: T) => void
  let reject: (reason: any) => void
  const promise = new Promise<T>((resolve2, reject2) => {
    resolve = resolve2
    reject = reject2
  })
  // @ts-ignore
  return { promise, resolve, reject }
}

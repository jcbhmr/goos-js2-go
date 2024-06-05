declare var DecompressionStream: any | undefined;
declare var ReadableStream: any | undefined;
declare var WritableStream: any | undefined;

export default async function gunzip(data: Uint8Array): Promise<Uint8Array> {
  if (globalThis.DecompressionStream) {
    const chunks: Uint8Array[] = [];
    let bytes: Uint8Array
    await new ReadableStream({
      start(controller) {
        controller.enqueue(data);
        controller.close();
      }
    }).pipeThrough(new DecompressionStream("gzip")).pipeTo(new WritableStream({
      write(chunk, controller) {
        chunks.push(chunk);
      },
      close(controller) {
        bytes = new Uint8Array(chunks.reduce((acc, chunk) => acc + chunk.byteLength, 0));
        let offset = 0;
        for (const chunk of chunks) {
          bytes.set(chunk, offset);
          offset += chunk.byteLength;
        }
      }
    }));
    return bytes;
  } else {
    throw new Error("not implemented")
  }
}

// src/lib/utils/hash-id.ts

import Hashids from "hashids";

const hashids = new Hashids("task-app-salt", 8);

export function encodeId(id: number) {
  return hashids.encode(id);
}

export function decodeId(value: string) {
  const [id] = hashids.decode(value);

  return typeof id === "number" ? id : null;
}

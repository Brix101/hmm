import { create } from "zustand";
import { devtools, persist } from "zustand/middleware";
import fileSlice, { FileSlice } from "./fileSlice";

type UnionToIntersection<U> = (
  U extends infer T ? (k: T) => void : never
) extends (k: infer I) => void
  ? I
  : never;

export type StoreState = UnionToIntersection<FileSlice>;

const useBoundStore = create<StoreState>()(
  devtools(
    persist(
      (...args) => ({
        ...fileSlice(...args),
      }),
      {
        name: "home-server",
      }
    )
  )
);
export { useBoundStore };

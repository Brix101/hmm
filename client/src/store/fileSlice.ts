import { StoreState } from "@/store";
import { StateCreator } from "zustand";

interface InitialState {
  file: {
    pathHistory?: string;
  };
}

interface Actions {
  setPathHistory: (path?: string) => void;
}

export type FileSlice = InitialState & Actions;

const initialState: InitialState = {
  file: {
    pathHistory: undefined,
  },
};

const fileSlice: StateCreator<StoreState, [], [], FileSlice> = (set) => ({
  ...initialState,
  setPathHistory: (fileUrls) =>
    set(() => ({
      file: { pathHistory: fileUrls },
    })),
});

export default fileSlice;

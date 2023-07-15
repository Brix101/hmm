import { StoreState } from "@/store";
import { StateCreator } from "zustand";

interface InitialState {
  file: {
    fileUrls?: string;
  };
}

interface Actions {
  appendToUrl: (path: string) => void;
  resetUrl: () => void;
}

export type FileSlice = InitialState & Actions;

const initialState: InitialState = {
  file: {
    fileUrls: undefined,
  },
};

const fileSlice: StateCreator<StoreState, [], [], FileSlice> = (set) => ({
  ...initialState,
  appendToUrl: (fileUrls) =>
    set(() => ({
      file: { fileUrls },
    })),
  resetUrl: () => set(() => initialState),
});

export default fileSlice;

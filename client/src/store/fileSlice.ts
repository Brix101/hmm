import { StoreState } from "@/store";
import { StateCreator } from "zustand";

interface InitialState {
  file: {
    activeFilePath?: string;
  };
}

interface Actions {
  setActiveFilePath: (path?: string) => void;
}

export type FileSlice = InitialState & Actions;

const initialState: InitialState = {
  file: {
    activeFilePath: undefined,
  },
};

const fileSlice: StateCreator<StoreState, [], [], FileSlice> = (set) => ({
  ...initialState,
  setActiveFilePath: (fileUrls) =>
    set(() => ({
      file: { activeFilePath: fileUrls },
    })),
});

export default fileSlice;

import { createContext } from "react-router";

export type DataEnvironment = {
  apiOrigin: string;
  allowOfflineFallback: boolean;
};

export const dataEnvironmentContext = createContext<DataEnvironment>();

export const browserDataEnvironment: DataEnvironment = {
  apiOrigin: "",
  allowOfflineFallback: true,
};

import React from "react";
import { Provider } from "react-redux";
import { store } from "./store";

interface StoreProps {
  children: React.ReactNode;
}

export const Store: React.FC<StoreProps> = ({ children }) => {
  return <Provider store={store}>{children}</Provider>;
};

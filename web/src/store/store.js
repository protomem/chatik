import { configureStore, getDefaultMiddleware } from "@reduxjs/toolkit";
import { authReducer } from "./auth/auth.slice";
import { authApi } from "../api/auth.api";

export const store = configureStore({
  reducer: {
    [authApi.reducerPath]: authApi.reducer,
    auth: authReducer,
  },
  devTools: true,
  middleware: [...getDefaultMiddleware(), authApi.middleware],
});

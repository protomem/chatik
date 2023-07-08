import { configureStore } from "@reduxjs/toolkit";
import { authReducer } from "./auth/auth.slice";
import { authApi } from "../api/auth.api";
import { channelsApi } from "../api/channels.api";
import { channelsReducer } from "./channels/channels.slice";

export const store = configureStore({
  reducer: {
    [authApi.reducerPath]: authApi.reducer,
    [channelsApi.reducerPath]: channelsApi.reducer,
    auth: authReducer,
    channels: channelsReducer,
  },
  devTools: true,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware()
      .concat(authApi.middleware)
      .concat(channelsApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

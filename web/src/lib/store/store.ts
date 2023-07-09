import { configureStore } from "@reduxjs/toolkit";
import { authReducer } from "./auth/auth.slice";
import { authApi } from "../api/auth.api";
import { channelsApi } from "../api/channels.api";
import { channelsReducer } from "./channels/channels.slice";
import { messagesReducer } from "./messages/messages.slice";
import { messagesApi } from "../api/messages.api";

export const store = configureStore({
  reducer: {
    [authApi.reducerPath]: authApi.reducer,
    [channelsApi.reducerPath]: channelsApi.reducer,
    [messagesApi.reducerPath]: messagesApi.reducer,
    auth: authReducer,
    channels: channelsReducer,
    messages: messagesReducer,
  },
  devTools: true,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware()
      .concat(authApi.middleware)
      .concat(channelsApi.middleware)
      .concat(messagesApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

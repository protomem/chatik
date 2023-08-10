import { configureStore } from "@reduxjs/toolkit";
import { authReducer } from "./auth/auth.slice";
import { authApi } from "./auth/auth.api";
import { channelsReducer } from "./channels/channels.slice";
import { channelsApi } from "./channels/channels.api";
import { eventsReducer } from "./events/events.slice";
import { messagesReducer } from "./messages/messages.slice";
import { messagesApi } from "./messages/messages.api";

export const store = configureStore({
  reducer: {
    [authApi.reducerPath]: authApi.reducer,
    [channelsApi.reducerPath]: channelsApi.reducer,
    [messagesApi.reducerPath]: messagesApi.reducer,
    auth: authReducer,
    channels: channelsReducer,
    events: eventsReducer,
    messages: messagesReducer,
  },
  devTools: true,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({ serializableCheck: false })
      .concat(authApi.middleware)
      .concat(channelsApi.middleware)
      .concat(messagesApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

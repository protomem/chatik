import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { eventsStorage } from "./events.storage";

export enum EventSourceType {
  WEBSOCKET = "WEBSOCKET",
  SSE = "SSE",
}

interface EventsState {
  sourceType: EventSourceType;
}

const initialState: EventsState = {
  sourceType: eventsStorage.getEventsType(),
};

const eventsSlice = createSlice({
  name: "events",
  initialState,
  reducers: {
    setSourceType: (state, { payload }: PayloadAction<EventSourceType>) => {
      state.sourceType = payload;
      eventsStorage.setEventsType(payload);
    },
  },
});

export const eventsReducer = eventsSlice.reducer;
export const eventsActions = eventsSlice.actions;

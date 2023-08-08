import { PayloadAction, createSlice } from "@reduxjs/toolkit";

export enum EventSourceType {
  WEBSOCKET,
  SSE,
}

interface EventsState {
  sourceType: EventSourceType;
}

const initialState: EventsState = {
  sourceType: EventSourceType.SSE,
};

const eventsSlice = createSlice({
  name: "events",
  initialState,
  reducers: {
    setSourceType: (state, { payload }: PayloadAction<EventSourceType>) => {
      state.sourceType = payload;
    },
  },
});

export const eventsReducer = eventsSlice.reducer;
export const eventsActions = eventsSlice.actions;

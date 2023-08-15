import { RootState } from "../store";

const selectEventsModule = (state: RootState) => state.events;

export const selectEventsSourceType = (state: RootState) =>
  selectEventsModule(state).sourceType;

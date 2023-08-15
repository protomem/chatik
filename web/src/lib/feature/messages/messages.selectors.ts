import { RootState } from "../store";

const selectMessagesModule = (state: RootState) => state.messages;

export const selectMessages = (state: RootState) =>
  selectMessagesModule(state).messages;

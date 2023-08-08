import { RootState } from "../store";

const selectChannelsModule = (state: RootState) => state.channels;

export const selectChannels = (state: RootState) =>
  selectChannelsModule(state).channels;
export const selectCurrentChannel = (state: RootState) =>
  selectChannelsModule(state).currentChannel;

import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { IChannel } from "../../entity/entities";

interface ChannelsState {
  channels: IChannel[];
  currentChannel: IChannel | null;
}

const initialState: ChannelsState = {
  channels: [],
  currentChannel: null,
};

const channelsSlice = createSlice({
  name: "channels",
  initialState,
  reducers: {
    setChannels: (state, { payload }: PayloadAction<IChannel[]>) => {
      state.channels = payload;
    },

    setCurrentChannel: (state, { payload }: PayloadAction<IChannel | null>) => {
      state.currentChannel = payload;
    },

    addChannel: (state, { payload }: PayloadAction<IChannel>) => {
      state.channels.push(payload);
    },

    removeChannel: (state, { payload }: PayloadAction<string>) => {
      state.channels = state.channels.filter(
        (channel) => channel.id !== payload
      );
    },
  },
});

export const channelsReducer = channelsSlice.reducer;
export const channelsActions = channelsSlice.actions;

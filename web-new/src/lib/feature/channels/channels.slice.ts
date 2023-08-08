import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { ChannelEntity } from "../../entities/entities";

interface ChannelsState {
  channels: ChannelEntity[];
  currentChannel: ChannelEntity | null;
}

const initialState: ChannelsState = {
  channels: [],
  currentChannel: null,
};

const channelsSlice = createSlice({
  name: "channels",
  initialState,
  reducers: {
    setCurrentChannel(state, { payload }: PayloadAction<ChannelEntity | null>) {
      state.currentChannel = payload;
    },

    setChannels(state, { payload }: PayloadAction<ChannelEntity[]>) {
      state.channels = payload;
    },

    addChannel(state, { payload }: PayloadAction<ChannelEntity>) {
      state.channels.push(payload);
    },

    removeChannel(state, { payload }: PayloadAction<string>) {
      state.channels = state.channels.filter(
        (channel) => channel.id !== payload,
      );
    },

    clearChannels(state) {
      state.channels = [];
    },
  },
});

export const channelsReducer = channelsSlice.reducer;
export const channelsActions = channelsSlice.actions;

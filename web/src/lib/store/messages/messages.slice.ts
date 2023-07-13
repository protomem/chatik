import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { IMessage } from "../../entity/entities";

interface MessagesState {
  messages: IMessage[];
}

const initialState: MessagesState = {
  messages: [],
};

export const messagesSlice = createSlice({
  name: "messages",
  initialState,
  reducers: {
    setMessages: (state, { payload }: PayloadAction<IMessage[]>) => {
      state.messages = payload;
    },

    addMessage: (state, { payload }: PayloadAction<IMessage>) => {
      state.messages.push(payload);
    },

    removeMessage: (state, { payload }: PayloadAction<string>) => {
      state.messages = state.messages.filter(
        (message) => message.id !== payload,
      );
    },
  },
});

export const messagesReducer = messagesSlice.reducer;
export const messagesActions = messagesSlice.actions;

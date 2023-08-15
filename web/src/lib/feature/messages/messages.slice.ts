import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { MessageEntity } from "../../entities/entities";

interface MessagesState {
  messages: MessageEntity[];
}

const initialState: MessagesState = {
  messages: [],
};

const messagesSlice = createSlice({
  name: "messages",
  initialState,
  reducers: {
    setMessages(state, { payload }: PayloadAction<MessageEntity[]>) {
      state.messages = payload;
    },

    addMessage(state, { payload }: PayloadAction<MessageEntity>) {
      if (!state.messages.find((message) => message.id === payload.id)) {
        state.messages.push(payload);
      }
    },

    removeMessage(state, { payload }: PayloadAction<string>) {
      state.messages = state.messages.filter(
        (message) => message.id !== payload,
      );
    },
  },
});

export const messagesReducer = messagesSlice.reducer;
export const messagesActions = messagesSlice.actions;

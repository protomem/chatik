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
      state.messages.push(payload);
    },
  },
});

export const messagesReducer = messagesSlice.reducer;
export const messagesActions = messagesSlice.actions;

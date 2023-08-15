import { ChannelEntity, MessageEntity } from "../../../entities/entities";

export enum EventTypes {
  NewMessage = "newMessage",
  RemoveMessage = "removeMessage",
  NewChannel = "newChannel",
  RemoveChannel = "removeChannel",
}

export interface Event {
  type: EventTypes;
  payload:
    | NewMessagePayload
    | RemoveMessagePayload
    | NewChannelPayload
    | RemoveChannelPayload;
}

export interface NewMessagePayload {
  message: MessageEntity;
}

export interface RemoveMessagePayload {
  messageId: string;
}

export interface NewChannelPayload {
  channel: ChannelEntity;
}

export interface RemoveChannelPayload {
  channelId: string;
}

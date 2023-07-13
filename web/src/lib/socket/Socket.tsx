import { useAppDispatch, useAppSelector } from "../store/hooks";
import { selectToken } from "../store/auth/auth.selectors";
import { useEffect, useState } from "react";
import { IChannel, IMessage } from "../entity/entities";
import { messagesActions } from "../store/messages/messages.slice";
import { channelsActions } from "../store/channels/channels.slice";
import { selectCurrentChannel } from "../store/channels/channels.selectors";

interface SocketProps {
  children: React.ReactNode;
}

enum EventTypes {
  NewMessage = "newMessage",
  RemoveMessage = "removeMessage",
  NewChannel = "newChannel",
  RemoveChannel = "removeChannel",
}

interface IEvent {
  type: EventTypes;
  payload:
    | NewMessagePayload
    | RemoveMessagePayload
    | NewChannelPayload
    | RemoveChannelPayload;
}

interface NewMessagePayload {
  message: IMessage;
}

interface RemoveMessagePayload {
  messageId: string;
}

interface NewChannelPayload {
  channel: IChannel;
}

interface RemoveChannelPayload {
  channelId: string;
}

export function Socket({ children }: SocketProps) {
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => selectToken(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const [socket, setSocket] = useState<WebSocket | null>(null);
  useEffect(() => {
    setSocket(
      new WebSocket(
        `ws://${import.meta.env.VITE_API_URL}/api/v1/stream?token=${token}`,
      ),
    );
  }, []);

  if (!!socket) {
    socket.onmessage = (e) => {
      const event: IEvent = JSON.parse(e.data);

      console.log(event);

      if (event.type === EventTypes.NewMessage) {
        const payload = event.payload as NewMessagePayload;

        if (payload.message.channelId === currentChannel?.id) {
          dispatch(messagesActions.addMessage(payload.message));
        }
      }

      if (event.type === EventTypes.RemoveMessage) {
        const payload = event.payload as RemoveMessagePayload;

        dispatch(messagesActions.removeMessage(payload.messageId));
      }

      if (event.type === EventTypes.NewChannel) {
        const payload = event.payload as NewChannelPayload;
        dispatch(channelsActions.addChannel(payload.channel));
      }

      if (event.type === EventTypes.RemoveChannel) {
        const payload = event.payload as RemoveChannelPayload;

        dispatch(channelsActions.removeChannel(payload.channelId));

        if (payload.channelId === currentChannel?.id) {
          dispatch(channelsActions.setCurrentChannel(null));
          dispatch(messagesActions.setMessages([]));
        }
      }
    };
  }

  return <>{children}</>;
}

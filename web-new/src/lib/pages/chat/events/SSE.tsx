import React, { useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "../../../feature/hooks";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";
import { selectCurrentChannel } from "../../../feature/channels/channels.selectors";
import {
  EventTypes,
  Event,
  NewChannelPayload,
  NewMessagePayload,
  RemoveChannelPayload,
  RemoveMessagePayload,
} from "./events";
import { messagesActions } from "../../../feature/messages/messages.slice";
import { channelsActions } from "../../../feature/channels/channels.slice";

interface SSEProps {
  children: React.ReactNode;
}

export const SSE: React.FC<SSEProps> = ({ children }) => {
  const dispatch = useAppDispatch();
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));
  const [sse, setSSE] = useState<EventSource | null>(null);

  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (!mounted) return;

    setSSE(
      new EventSource(
        `http://${
          import.meta.env.VITE_API_URL
        }/api/v1/stream/sse?token=${accessToken}`,
      ),
    );
  }, [mounted]);

  useEffect(() => {
    if (sse === null) return;

    sse.onmessage = (eventSSE) => {
      const event: Event = JSON.parse(eventSSE.data);
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
  }, [sse, currentChannel]);

  return <>{children}</>;
};

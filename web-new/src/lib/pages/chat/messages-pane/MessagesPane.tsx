import { Box, Sheet, Stack } from "@mui/joy";
import React, { useEffect } from "react";
import { MessagesPaneHeader } from "../messages-pane-header/MessagesPaneHeader";
import { useAppDispatch, useAppSelector } from "../../../feature/hooks";
import {
  selectAccessToken,
  selectCurrentUser,
} from "../../../feature/auth/auth.selectors";
import { MessageInput } from "../message-input/MessageInput";
import { MessagesItem } from "../messages-item/MessagesItem";
import { selectCurrentChannel } from "../../../feature/channels/channels.selectors";
import { selectMessages } from "../../../feature/messages/messages.selectors";
import { useGetAllMessagesQuery } from "../../../feature/messages/messages.api";
import { messagesActions } from "../../../feature/messages/messages.slice";

export const MessagesPane: React.FC = () => {
  const dispatch = useAppDispatch();
  const currentUser = useAppSelector((state) => selectCurrentUser(state));
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const messages = useAppSelector((state) => selectMessages(state));

  const { data, error } = useGetAllMessagesQuery(
    {
      channelId: currentChannel?.id || "",
      accessToken,
    },
    { skip: !currentChannel, refetchOnMountOrArgChange: true },
  );

  useEffect(() => {
    if (!!data) dispatch(messagesActions.setMessages(data.messages));

    if (!!error) dispatch(messagesActions.setMessages([]));
  }, [dispatch, data, error]);

  return (
    <Sheet
      sx={{
        height: {
          xs: "calc(100dvh - var(--Header-height))",
        },
        display: "flex",
        flexDirection: "column",
      }}
    >
      {!!currentChannel && (
        <>
          <MessagesPaneHeader channel={currentChannel} />

          <Box
            sx={{
              display: "flex",
              flex: 1,
              minHeight: 0,
              px: 2,
              py: 2.5,
              overflowY: "scroll",
              flexDirection: "column-reverse",
              gap: 2,
            }}
          >
            {messages
              .slice(0)
              .reverse()
              .map((message) => {
                const isYou = message.user.id === currentUser?.id;
                return (
                  <Stack
                    key={message.id}
                    direction={"row"}
                    spacing={2}
                    flexDirection={isYou ? "row-reverse" : "row"}
                  >
                    <MessagesItem message={message} />
                  </Stack>
                );
              })}
          </Box>

          <MessageInput channel={currentChannel} />
        </>
      )}
    </Sheet>
  );
};

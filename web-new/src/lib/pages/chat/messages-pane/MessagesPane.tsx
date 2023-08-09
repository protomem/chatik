import { Box, Sheet, Stack } from "@mui/joy";
import React from "react";
import { MessagesPaneHeader } from "../messages-pane-header/MessagesPaneHeader";
import { useAppSelector } from "../../../feature/hooks";
import { selectCurrentUser } from "../../../feature/auth/auth.selectors";
import {
  ChannelEntity,
  MessageEntity,
  UserEntity,
} from "../../../entities/entities";
import { MessageInput } from "../message-input/MessageInput";
import { MessagesItem } from "../messages-item/MessagesItem";

export const MessagesPane: React.FC = () => {
  const currentUser = useAppSelector((state) => selectCurrentUser(state));
  const currentChannel: ChannelEntity = {
    id: "1",
    createdAt: new Date(),
    updatedAt: new Date(),
    title: "New channel",
    user: currentUser || ({} as UserEntity),
  };

  const messages = [
    {
      id: "1",
      createdAt: new Date(),
      content: "Hello",
      user: currentUser || ({} as UserEntity),
    },
    {
      id: "2",
      createdAt: new Date(),
      content: "How are you?",
      user: {
        id: "1",
        nickname: "John",
      },
    },
    {
      id: "3",
      createdAt: new Date(),
      content: "Hello",
      user: currentUser || ({} as UserEntity),
    },
    {
      id: "4",
      createdAt: new Date(),
      content: "How are you?",
      user: {
        id: "1",
        nickname: "John",
      },
    },
{
      id: "6",
      createdAt: new Date(),
      content: "Hello",
      user: currentUser || ({} as UserEntity),
    },
    {
      id: "7",
      createdAt: new Date(),
      content: "How are you?",
      user: {
        id: "1",
        nickname: "John",
      },
    },
  ] as MessageEntity[];

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
            }}
          >
            {messages.map((message) => {
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

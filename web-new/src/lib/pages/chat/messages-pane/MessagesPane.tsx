import { Box, Sheet } from "@mui/joy";
import React from "react";
import { MessagesPaneHeader } from "../messages-pane-header/MessagesPaneHeader";
import { useAppSelector } from "../../../feature/hooks";
import { selectCurrentUser } from "../../../feature/auth/auth.selectors";
import { ChannelEntity, UserEntity } from "../../../entities/entities";
import { MessageInput } from "../message-input/MessageInput";

export const MessagesPane: React.FC = () => {
  const currentUser = useAppSelector((state) => selectCurrentUser(state));
  const currentChannel: ChannelEntity = {
    id: "1",
    createdAt: new Date(),
    updatedAt: new Date(),
    title: "New channel",
    user: currentUser || ({} as UserEntity),
  };

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
          ></Box>

          <MessageInput channel={currentChannel} />
        </>
      )}
    </Sheet>
  );
};

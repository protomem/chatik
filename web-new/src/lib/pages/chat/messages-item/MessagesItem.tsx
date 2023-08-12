import React, { useState } from "react";
import { ChannelEntity, MessageEntity } from "../../../entities/entities";
import { Box, IconButton, Sheet, Stack, Typography } from "@mui/joy";
import { useAppSelector } from "../../../feature/hooks";
import {
  selectAccessToken,
  selectCurrentUser,
} from "../../../feature/auth/auth.selectors";
import { DeleteRounded } from "@mui/icons-material";
import { useDeleteMessageMutation } from "../../../feature/messages/messages.api";

interface MessagesItemProps {
  channel: ChannelEntity;
  message: MessageEntity;
}

export const MessagesItem: React.FC<MessagesItemProps> = ({
  channel,
  message,
}) => {
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const currentUser = useAppSelector((state) => selectCurrentUser(state));
  const [isHovered, setIsHovered] = useState(false);
  const [deleteMessage] = useDeleteMessageMutation();

  return (
    <Box>
      <Stack
        direction="row"
        justifyContent="space-between"
        spacing={2}
        sx={{ mb: 0.25 }}
      >
        <Typography level="body-xs">
          {message.user.id === currentUser?.id ? "You" : message.user.nickname}
        </Typography>
        <Typography level="body-xs">{message.createdAt.toString()}</Typography>
      </Stack>

      <Box
        sx={{ position: "relative" }}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <Sheet
          color={message.user.id === currentUser?.id ? "primary" : "neutral"}
          variant={message.user.id === currentUser?.id ? "solid" : "soft"}
          sx={{
            px: 1.25,
            py: 1.25,
            borderRadius: "lg",
            borderTopRightRadius:
              message.user.id === currentUser?.id ? 0 : "lg",
            borderTopLeftRadius: message.user.id === currentUser?.id ? "lg" : 0,
          }}
        >
          {message.content}
        </Sheet>

        {isHovered && message.user.id === currentUser?.id && (
          <Stack
            direction="row"
            justifyContent={
              message.user.id === currentUser?.id ? "flex-end" : "flex-start"
            }
            spacing={0.5}
            sx={{
              position: "absolute",
              top: "50%",
              p: 1.5,
              ...(message.user.id === currentUser?.id
                ? {
                    left: 0,
                    transform: "translate(-100%, -50%)",
                  }
                : {
                    right: 0,
                    transform: "translate(100%, -50%)",
                  }),
            }}
          >
            <IconButton
              variant="soft"
              color="danger"
              size="sm"
              onClick={() => {
                deleteMessage({
                  messageId: message.id,
                  channelId: channel.id,
                  accessToken,
                })
                  .unwrap()
                  .catch((err) => {
                    alert(err);
                  });
              }}
            >
              <DeleteRounded />
            </IconButton>
          </Stack>
        )}
      </Box>
    </Box>
  );
};

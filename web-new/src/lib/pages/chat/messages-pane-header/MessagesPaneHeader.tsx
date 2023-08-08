import { IconButton, Stack, Typography } from "@mui/joy";
import React from "react";
import { ChannelEntity } from "../../../entities/entities";
import { MoreVert } from "@mui/icons-material";

interface MessagesPaneHeaderProps {
  channel: ChannelEntity;
}

export const MessagesPaneHeader: React.FC<MessagesPaneHeaderProps> = ({
  channel,
}) => {
  return (
    <Stack
      direction="row"
      justifyContent="space-between"
      sx={{
        borderBottom: "1px solid",
        borderColor: "divider",
      }}
      py={{ xs: 2, md: 2 }}
      px={{ xs: 1, md: 2 }}
    >
      <div>
        <Typography fontWeight="lg" fontSize="lg" component="h2" noWrap>
          {channel.title}
        </Typography>

        <Typography level="body-sm">{channel.user.nickname}</Typography>
      </div>

      <IconButton variant="plain" color="neutral">
        <MoreVert />
      </IconButton>
    </Stack>
  );
};

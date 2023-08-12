import { Stack, Typography } from "@mui/joy";
import React from "react";
import { ChannelEntity } from "../../../entities/entities";
import { ChannelMenu } from "../channel-menu/ChannelMenu";
import { useAppSelector } from "../../../feature/hooks";
import { selectCurrentUser } from "../../../feature/auth/auth.selectors";

interface MessagesPaneHeaderProps {
  channel: ChannelEntity;
}

export const MessagesPaneHeader: React.FC<MessagesPaneHeaderProps> = ({
  channel,
}) => {
  const currentUser = useAppSelector((state) => selectCurrentUser(state));

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

      {currentUser?.id === channel.user.id && <ChannelMenu channel={channel} />}
    </Stack>
  );
};

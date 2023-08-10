import React from "react";
import { ChannelEntity } from "../../../entities/entities";
import { ListDivider, ListItem, ListItemButton, Typography } from "@mui/joy";
import { useAppDispatch } from "../../../feature/hooks";
import { channelsActions } from "../../../feature/channels/channels.slice";

interface ChannelsItemProps {
  channel: ChannelEntity;
  selected: boolean;
}

export const ChannelsItem: React.FC<ChannelsItemProps> = ({
  channel,
  selected,
}) => {
  const dispatch = useAppDispatch();

  return (
    <>
      <ListItem>
        <ListItemButton
          selected={selected}
          color="neutral"
          sx={{
            flexDirection: "column",
            alignItems: "initial",
            gap: 1,
            fontWeight: "normal",
          }}
          onClick={() => {
            dispatch(channelsActions.setCurrentChannel(channel));
          }}
        >
          <Typography sx={{ my: 1 }}>{channel.title}</Typography>
          {/* TODO: Add preview last messages */}
        </ListItemButton>
      </ListItem>
      <ListDivider sx={{ margin: 0 }} />
    </>
  );
};

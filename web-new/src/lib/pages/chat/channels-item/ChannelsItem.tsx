import React from "react";
import { ChannelEntity } from "../../../entities/entities";
import { ListDivider, ListItem, ListItemButton, Typography } from "@mui/joy";

interface ChannelsItemProps {
  channel: ChannelEntity;
  selected: boolean;
}

export const ChannelsItem: React.FC<ChannelsItemProps> = ({
  channel,
  selected,
}) => {
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
        >
          <Typography sx={{ my: 1 }}>{channel.title}</Typography>
          {/* TODO: Add preview last messages */}
        </ListItemButton>
      </ListItem>
      <ListDivider sx={{ margin: 0 }} />
    </>
  );
};

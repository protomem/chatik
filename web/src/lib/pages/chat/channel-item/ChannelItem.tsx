import { ListItem, ListItemButton } from "@mui/material";
import { IChannel } from "../../../entity/entities";
import { useAppDispatch, useAppSelector } from "../../../store/hooks";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";
import React from "react";
import { channelsActions } from "../../../store/channels/channels.slice";

interface ChannelItemProps {
  channel: IChannel;
}

export function ChannelItem({ channel }: ChannelItemProps) {
  const dispatch = useAppDispatch();
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const isActive = channel.id === currentChannel?.id;

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();

    if (!isActive) dispatch(channelsActions.setCurrentChannel(channel));
  };

  return (
    <ListItem
      sx={{
        width: "100%",
        height: "80px",
        p: 0,
        bgcolor: isActive ? "primary.main" : "background.default",
        color: isActive ? "primary.contrastText" : "text.primary",
      }}
    >
      <ListItemButton
        sx={{ width: "100%", height: "100%" }}
        onClick={handleClick}
      >
        {channel.title}
      </ListItemButton>
    </ListItem>
  );
}

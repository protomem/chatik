import { ListItem, ListItemButton } from "@mui/material";
import { IChannel } from "../../../entity/entities";

interface ChannelItemProps {
  channel: IChannel;
}

export function ChannelItem({ channel }: ChannelItemProps) {
  return (
    <ListItem sx={{ width: "100%", height: "500px", p: 0 }}>
      <ListItemButton sx={{ width: "100%", height: "100%" }}>
        {channel.title}
      </ListItemButton>
    </ListItem>
  );
}

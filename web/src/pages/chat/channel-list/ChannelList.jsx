import { List } from "@mui/material";
import { ChannelItem } from "../channel-item/ChannelItem";

export function ChannelList() {
  return (
    <List>
      {Array.from({ length: 30 }).map((_, index) => (
        <ChannelItem key={index} index={index + 1} />
      ))}
    </List>
  );
}

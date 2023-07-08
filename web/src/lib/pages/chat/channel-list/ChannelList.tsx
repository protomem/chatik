import { List } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectChannels } from "../../../store/channels/channels.selectors";
import { ChannelItem } from "../channel-item/ChannelItem";

export function ChannelList() {
  const channels = useAppSelector((state) => selectChannels(state));

  return (
    <List>
      {channels.map((channel) => (
        <ChannelItem key={channel.id} channel={channel} />
      ))}
    </List>
  );
}

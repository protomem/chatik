import { ListItemButton } from "@mui/material";

export function ChannelItem({ index }) {
  return (
    <ListItemButton sx={{ padding: 3 }}>Channel Item {index}</ListItemButton>
  );
}

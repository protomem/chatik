import { Box, IconButton, Typography } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";
import { MoreVert } from "@mui/icons-material";
import { ChannelMenu } from "../channel-menu/ChannelMenu";
import { useState } from "react";

export function DialogWindowHeader() {
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const [menuAnchor, setMenuAnchor] = useState<HTMLElement | null>(null);
  const menuOpen = Boolean(menuAnchor);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setMenuAnchor(event.currentTarget);
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "row",
        justifyContent: "center",
        alignItems: "center",
        width: "100%",
        gap: 2,
      }}
    >
      {Boolean(currentChannel) && (
        <>
          <Typography variant="h6">{currentChannel?.title}</Typography>
          <IconButton
            id="channel-button"
            aria-controls={menuOpen ? "channel-menu" : undefined}
            aria-haspopup="true"
            aria-expanded={menuOpen ? "true" : undefined}
            onClick={handleClick}
          >
            <MoreVert />
          </IconButton>
          <ChannelMenu anchor={menuAnchor} setAnchor={setMenuAnchor} />
        </>
      )}
    </Box>
  );
}

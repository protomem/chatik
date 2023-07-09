import { Box, IconButton, Typography } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";
import { MoreVert } from "@mui/icons-material";

export function DialogWindowHeader() {
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

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
      <Typography variant="h6">{currentChannel?.title}</Typography>
      {Boolean(currentChannel) && (
        <IconButton>
          <MoreVert />
        </IconButton>
      )}
    </Box>
  );
}

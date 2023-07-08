import { Box, Typography } from "@mui/material";
import { ChatBar } from "./chat-bar/ChatBar";
import { useEffect } from "react";
import { useGetListChannelsQuery } from "../../api/channels.api";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import { channelsActions } from "../../store/channels/channels.slice";
import { selectToken } from "../../store/auth/auth.selectors";
import { NewChannelButton } from "./new-channel-button/NewChannelButton";
import { ChannelList } from "./channel-list/ChannelList";

export function Chat() {
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => selectToken(state));
  const { data } = useGetListChannelsQuery({ token });

  useEffect(() => {
    if (data?.channels) dispatch(channelsActions.setChannels(data.channels));
  }, [data?.channels, dispatch]);

  return (
    <Box sx={{ flexGrow: 1, height: "100vh" }}>
      <Box sx={{ height: "7vh" }}>
        <ChatBar />
      </Box>

      <Box
        sx={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "center",
          height: "93vh",
        }}
      >
        <Box
          sx={{
            flexGrow: 2,
            height: "100%",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <Box
            sx={{
              my: 2,
              p: 0,
              width: "100%",
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <NewChannelButton />
          </Box>

          <Box sx={{ flexGrow: 1, width: "100%", overflow: "auto" }}>
            <ChannelList />
          </Box>
        </Box>

        <Box
          sx={{
            flexGrow: 8,
            height: "100%",
            bgcolor: "red",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <Box
            sx={{
              flexGrow: 14,
              width: "100%",
              bgcolor: "green",
            }}
          >
            <Typography variant="h5">Dialog Window</Typography>
          </Box>

          <Box
            sx={{
              flexGrow: 2,
              width: "100%",
              bgcolor: "yellow",
            }}
          >
            <Typography variant="h5">New Message Input</Typography>
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

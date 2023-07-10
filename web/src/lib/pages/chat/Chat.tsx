import { Box } from "@mui/material";
import { ChatBar } from "./chat-bar/ChatBar";
import { useEffect } from "react";
import { useGetListChannelsQuery } from "../../api/channels.api";
import { useAppDispatch, useAppSelector } from "../../store/hooks";
import { channelsActions } from "../../store/channels/channels.slice";
import { selectToken } from "../../store/auth/auth.selectors";
import { NewChannelButton } from "./new-channel-button/NewChannelButton";
import { ChannelList } from "./channel-list/ChannelList";
import { DialogWindow } from "./dialog-window/DialogWIndow";
import { NewMessageInput } from "./new-message-input/NewMessageInput";

export function Chat() {
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => selectToken(state));
  const { data: channelsData } = useGetListChannelsQuery({ token });

  useEffect(() => {
    if (channelsData?.channels)
      dispatch(channelsActions.setChannels(channelsData.channels));
  }, [channelsData?.channels, dispatch]);

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
            maxWidth: "20%",
            minWidth: "20%",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            borderRight: 1,
            borderColor: "divider",
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

          <Box
            sx={{
              flexGrow: 1,
              width: "100%",
              overflow: "auto",
            }}
          >
            <ChannelList />
          </Box>
        </Box>

        <Box
          sx={{
            flexGrow: 10,
            height: "100%",
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
              overflow: "auto",
            }}
          >
            <DialogWindow />
          </Box>

          <Box
            sx={{
              flexGrow: 2,
              width: "100%",
              borderTop: 1,
              borderColor: "divider",
              boxSizing: "border-box",
              minHeight: "11%",
              maxHeight: "11%",
            }}
          >
            <NewMessageInput />
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

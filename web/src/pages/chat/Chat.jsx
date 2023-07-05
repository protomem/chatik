import { Box } from "@mui/material";
import { ChatAppBar } from "./chat-app-bar/ChatAppBar";
import { ChannelList } from "./channel-list/ChannelList";
import { DialogWindow } from "./dialog-window/DialogWindow";
import { NewMessageInput } from "./new-message-input/NewMessageInput";

export function Chat() {
  return (
    <Box sx={{ flexGrow: 1, height: "100vh" }}>
      <ChatAppBar />

      <Box
        sx={{
          flexGrow: 1,
          display: "flex",
          flexDirection: "row",
          height: "93vh",
        }}
      >
        <Box
          sx={{
            flexGrow: 1,
            height: "100%",
            overflowY: "auto",
            overflowX: "hidden",
            borderRight: 1,
            borderColor: "divider",
            maxWidth: "20%",
          }}
        >
          <ChannelList />
        </Box>

        <Box
          sx={{
            flexGrow: 8,
            display: "flex",
            flexDirection: "column",
            maxWidth: "80%",
          }}
        >
          <Box sx={{ flexGrow: 11, overflowY: "auto" }}>
            <DialogWindow />
          </Box>

          <Box
            sx={{
              flexGrow: 1,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              borderTop: 1,
              borderColor: "divider",
              px: 4,
              height: "20%",
            }}
          >
            <NewMessageInput />
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

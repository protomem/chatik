import { Send } from "@mui/icons-material";
import { Box, Button, TextField } from "@mui/material";
import { useState } from "react";
import { useAppSelector } from "../../../store/hooks";
import { selectToken } from "../../../store/auth/auth.selectors";
import { useCreateMessageMutation } from "../../../api/messages.api";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";

export function NewMessageInput() {
  const token = useAppSelector((state) => selectToken(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const [messageInput, setMessageInput] = useState("");
  const clearMessageInput = () => setMessageInput("");

  const [createMessage] = useCreateMessageMutation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const res = createMessage({
      channelId: currentChannel?.id ?? "",
      token,
      content: messageInput,
    });

    res.catch((err) => {
      console.log(err);
      alert(`${err}`);
    });

    clearMessageInput();
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "row",
        justifyContent: "center",
        alignItems: "center",
        height: "100%",
        width: "100%",
      }}
      component="form"
      onSubmit={handleSubmit}
    >
      <TextField
        sx={{ ml: 4, mb: 1 }}
        fullWidth
        multiline
        maxRows={3}
        placeholder="Type a message..."
        value={messageInput}
        onChange={(e) => setMessageInput(e.target.value)}
      />
      <Button
        sx={{ mx: 4, mb: 1, height: "50%" }}
        variant="contained"
        type="submit"
      >
        <Send sx={{ mr: 1 }} />
        Send
      </Button>
    </Box>
  );
}

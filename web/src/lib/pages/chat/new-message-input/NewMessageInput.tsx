import { Send } from "@mui/icons-material";
import { Box, Button, TextField } from "@mui/material";
import { useState } from "react";

export function NewMessageInput() {
  const [messageInput, setMessageInput] = useState("");

  // TODO: Implement this
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    console.log(messageInput);

    setMessageInput("");
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
      <Button sx={{ mx: 4, height: "50%" }} variant="contained" type="submit">
        <Send sx={{ mr: 1 }} />
        Send
      </Button>
    </Box>
  );
}

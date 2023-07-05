import { Button, TextField, Box } from "@mui/material";
import { Send } from "@mui/icons-material";

export function NewMessageInput() {
  return (
    <Box
      sx={{
        flexGrow: 1,
        display: "flex",
        flexDirection: "row",
        justifyContent: "center",
        alignItems: "center",
        gap: 2,
      }}
      component="form"
      onSubmit={(e) => {
        e.preventDefault();
        console.log("send new message");
      }}
    >
      <Box sx={{ flexGrow: 8, px: 4 }}>
        <TextField
          fullWidth
          multiline
          maxRows={3}
          placeholder="Type a message"
        />
      </Box>

      <Box sx={{ flexGrow: 1 }}>
        <Button
          sx={{ flexGrow: 1 }}
          size="large"
          variant="contained"
          type="submit"
        >
          <Send />
          Send
        </Button>
      </Box>
    </Box>
  );
}

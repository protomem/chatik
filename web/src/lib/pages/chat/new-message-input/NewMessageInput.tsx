import { Box, Button, TextField } from "@mui/material";

export function NewMessageInput() {
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
    >
      <TextField />
      <Button>Send</Button>
    </Box>
  );
}

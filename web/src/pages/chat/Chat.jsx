import { AppBar, Box, Button, Typography, Toolbar } from "@mui/material";
import { useNavigate } from "react-router-dom";

export function Chat() {
  const nav = useNavigate();

  return (
    <Box>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h4">Chat</Typography>

          <Typography></Typography>
          <Button>Button</Button>
        </Toolbar>
      </AppBar>
    </Box>
  );
}

import { AppBar, Box, Typography, Toolbar, IconButton } from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import { LogoutButton } from "./logout-button/LogoutButton";
import { useSelector } from "react-redux";
import { selectCurrentUser } from "../../store/auth/auth.selectors";
import { ChatAppBar } from "./chat-app-bar/ChatAppBar";

export function Chat() {
  return (
    <Box>
      <ChatAppBar />
    </Box>
  );
}

import { AppBar, Box, Typography, Toolbar, IconButton } from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import { LogoutButton } from "../logout-button/LogoutButton";
import { useSelector } from "react-redux";
import { selectCurrentUser } from "../../../store/auth/auth.selectors";

export function ChatAppBar() {
  const currentUser = useSelector((state) => selectCurrentUser(state));

  return (
    <AppBar position="static">
      <Toolbar>
        <Box
          sx={{
            flexGrow: 1,
            display: "flex",
            flexDirection: "row",
            justifyContent: "flex-start",
            alignItems: "center",
            gap: 2,
          }}
        >
          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h4" component={"div"}>
            Chat
          </Typography>
        </Box>

        <Box
          sx={{
            flexGrow: 1,
            display: "flex",
            flexDirection: "row",
            justifyContent: "flex-end",
            alignItems: "center",
            gap: 2,
          }}
        >
          <Typography variant="h6">{currentUser.nickname}</Typography>
          <LogoutButton />
        </Box>
      </Toolbar>
    </AppBar>
  );
}

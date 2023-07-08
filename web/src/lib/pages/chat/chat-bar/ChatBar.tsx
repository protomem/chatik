import { AppBar, Box, Toolbar, Typography } from "@mui/material";
import { LogoutButton } from "../logout-button/LogoutButton";
import { useAppSelector } from "../../../store/hooks";
import { selectUser } from "../../../store/auth/auth.selectors";

export function ChatBar() {
  const currentUser = useAppSelector((state) => selectUser(state));

  return (
    <AppBar position="static">
      <Toolbar
        sx={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <Box>
          <Typography variant="h5">Chat</Typography>
        </Box>

        <Box
          sx={{
            display: "flex",
            flexDirection: "row",
            justifyContent: "center",
            alignItems: "center",
            gap: 4,
          }}
        >
          <Typography variant="h6">{currentUser?.nickname}</Typography>
          <LogoutButton />
        </Box>
      </Toolbar>
    </AppBar>
  );
}

import { AppBar, Toolbar, Typography } from "@mui/material";

export function AuthBar() {
  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h5">Auth</Typography>
      </Toolbar>
    </AppBar>
  );
}

import {
  AppBar,
  Box,
  Button,
  Container,
  Toolbar,
  Typography,
} from "@mui/material";
import { useNavigate } from "react-router-dom";

export function NotFound() {
  const nav = useNavigate();

  const handleClick = () => {
    nav("/", { replace: true });
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h4">Not Found</Typography>
        </Toolbar>
      </AppBar>

      <Container
        sx={{
          mt: 4,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          gap: 2,
        }}
      >
        <Typography variant="h6" align="center">
          404 Page not found
        </Typography>

        <Button onClick={handleClick} variant="outlined">
          Go home
        </Button>
      </Container>
    </Box>
  );
}

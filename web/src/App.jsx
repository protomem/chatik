import { CssBaseline, ThemeProvider, createTheme } from "@mui/material";
import { Router } from "./router/Router";
import { StoreProvider } from "./store/StoreProvider";
import "@fontsource/jetbrains-mono";

export function App() {
  const theme = createTheme({
    palette: {
      mode: "light",
    },
    typography: {
      fontFamily: "JetBrains Mono",
    },
  });

  return (
    <ThemeProvider theme={theme}>
      <StoreProvider>
        <CssBaseline>
          <Router />
        </CssBaseline>
      </StoreProvider>
    </ThemeProvider>
  );
}

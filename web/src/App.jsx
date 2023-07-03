import { CssBaseline, ThemeProvider, createTheme } from "@mui/material";
import "@fontsource/jetbrains-mono";
import { Router } from "./router/Router";

export function App() {
  const theme = createTheme({
    palette: {
      mode: "dark",
    },
    typography: {
      fontFamily: "JetBrains Mono",
    },
  });

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline>
        <Router />
      </CssBaseline>
    </ThemeProvider>
  );
}

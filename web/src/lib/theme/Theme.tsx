import { CssBaseline, ThemeProvider, createTheme } from "@mui/material";
import "@fontsource/jetbrains-mono";

interface ThemeProps {
  children: React.ReactNode;
}

export function Theme({ children }: ThemeProps) {
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
      <CssBaseline>{children}</CssBaseline>
    </ThemeProvider>
  );
}

import React from "react";
import { CssBaseline, CssVarsProvider, StyledEngineProvider } from "@mui/joy";

interface ThemeProps {
  children: React.ReactNode;
}

export const Theme: React.FC<ThemeProps> = ({ children }) => {
  return (
    <StyledEngineProvider injectFirst>
      <CssVarsProvider defaultMode="dark">
        <CssBaseline />

        {children}
      </CssVarsProvider>
    </StyledEngineProvider>
  );
};

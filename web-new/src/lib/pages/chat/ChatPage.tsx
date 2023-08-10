import { Box, GlobalStyles } from "@mui/joy";
import React from "react";
import { Header } from "./header/Header";
import { Messenger } from "./messenger/Messenger";
import { Events } from "./events/Events";

export const ChatPage: React.FC = () => {
  return (
    <>
      <GlobalStyles
        styles={(theme) => ({
          "[data-feather], .feather": {
            color: `var(--Icon-color, ${theme.vars.palette.text.icon})`,
            margin: "var(--Icon-margin)",
            fontSize: `var(--Icon-fontSize, ${theme.vars.fontSize.xl})`,
            width: "1em",
            height: "1em",
          },
        })}
      />
      <Box sx={{ display: "flex", minHeight: "100dvh" }}>
        <Header />
        <Box component="main" className="MainContent" flex={1}>
          <Events>
            <Messenger />
          </Events>
        </Box>
      </Box>
    </>
  );
};

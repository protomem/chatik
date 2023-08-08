import { Sheet } from "@mui/joy";
import React from "react";
import { ChannelsPane } from "../channels-pane/ChannelsPane.tsx";

export const Messenger: React.FC = () => {
  return (
    <Sheet
      sx={{
        flex: 1,
        width: "100%",
        mx: "auto",
        pt: { xs: "var(--Header-height)" },
        display: "grid",
        gridTemplateColumns: {
          xs: "1fr",
          sm: "minmax(min-content, min(30%, 400px)) 1fr",
        },
      }}
    >
      <Sheet
        sx={{
          position: {
            xs: "fixed",
            sm: "sticky",
          },
          transform: {
            xs: "translateX(calc(100% * (var(--MessagesPane-slideIn, 0) - 1)))",
            sm: "none",
          },
          transition: "transform 0.4s, width 0.4s",
          zIndex: 100,
          width: "100%",
          top: 52,
        }}
      >
        <ChannelsPane />
      </Sheet>
    </Sheet>
  );
};

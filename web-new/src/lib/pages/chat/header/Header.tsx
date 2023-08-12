import { Box, Button, GlobalStyles, Sheet, Switch, Typography } from "@mui/joy";
import React from "react";
import { ColorSchemeToggle } from "../color-scheme-toggle/ColorSchemeToggle";
import { useAppDispatch, useAppSelector } from "../../../feature/hooks";
import { selectCurrentUser } from "../../../feature/auth/auth.selectors";
import { authActions } from "../../../feature/auth/auth.slice";
import {
  EventSourceType,
  eventsActions,
} from "../../../feature/events/events.slice";
import { selectEventsSourceType } from "../../../feature/events/events.selectors";

export const Header: React.FC = () => {
  const dispatch = useAppDispatch();
  const currentUser = useAppSelector((state) => selectCurrentUser(state));
  const eventsSourceType = useAppSelector((state) =>
    selectEventsSourceType(state),
  );

  return (
    <Sheet
      sx={{
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        position: "fixed",
        top: 0,
        width: "100vw",
        height: "var(--Header-height)",
        zIndex: 9995,
        py: 1,
        px: 2,
        gap: 1,
        boxShadow: "sm",
      }}
    >
      <GlobalStyles
        styles={() => ({
          ":root": {
            "--Header-height": "52px",
          },
        })}
      />

      <Box>
        <Typography
          fontWeight="lg"
          component="a"
          href="/"
          sx={{ textDecoration: "none" }}
          startDecorator={
            <Box
              component="span"
              sx={{
                width: 24,
                height: 24,
                background: (theme) =>
                  `linear-gradient(45deg, ${theme.vars.palette.primary.solidBg}, ${theme.vars.palette.primary.solidBg} 30%, ${theme.vars.palette.primary.softBg})`,
                borderRadius: "50%",
                boxShadow: (theme) => theme.shadow.md,
                "--joy-shadowChannel": (theme) =>
                  theme.vars.palette.primary.mainChannel,
              }}
            />
          }
        >
          Chatik
        </Typography>
      </Box>

      <Box
        sx={{
          display: "flex",
          flexDirection: "row",
          gap: 3,
          alignItems: "center",
        }}
      >
        <Switch
          disabled
          size="sm"
          color={
            eventsSourceType === EventSourceType.SSE ? "success" : "primary"
          }
          sx={{ mr: 3 }}
          startDecorator={
            <Typography
              fontWeight={"lg"}
              color={
                eventsSourceType === EventSourceType.SSE ? "success" : "neutral"
              }
            >
              sse
            </Typography>
          }
          endDecorator={
            <Typography
              fontWeight={"lg"}
              color={
                eventsSourceType === EventSourceType.WEBSOCKET
                  ? "primary"
                  : "neutral"
              }
            >
              ws
            </Typography>
          }
          onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
            !event.target.checked
              ? dispatch(eventsActions.setSourceType(EventSourceType.SSE))
              : dispatch(
                  eventsActions.setSourceType(EventSourceType.WEBSOCKET),
                );
          }}
          checked={eventsSourceType !== EventSourceType.SSE}
        />
        <Typography>{currentUser?.nickname}</Typography>
        <Button
          size="sm"
          variant="soft"
          color="danger"
          onClick={(event) => {
            event.preventDefault();
            dispatch(authActions.clearCredentials());
            dispatch(eventsActions.setSourceType(EventSourceType.SSE));
          }}
        >
          logout
        </Button>
        <ColorSchemeToggle id={undefined} />
      </Box>
    </Sheet>
  );
};

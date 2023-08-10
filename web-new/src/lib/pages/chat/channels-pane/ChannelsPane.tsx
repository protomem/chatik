import {
  Box,
  IconButton,
  Input,
  List,
  ListDivider,
  Sheet,
  Stack,
  Typography,
} from "@mui/joy";
import React, { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../../../feature/hooks";
import {
  selectChannels,
  selectCurrentChannel,
} from "../../../feature/channels/channels.selectors";
import { Add, ArrowDropDownOutlined, Search } from "@mui/icons-material";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";
import { ChannelsItem } from "../channels-item/ChannelsItem";
import { useGetAllChannelsQuery } from "../../../feature/channels/channels.api";
import { channelsActions } from "../../../feature/channels/channels.slice";

export const ChannelsPane: React.FC = () => {
  const dispatch = useAppDispatch();
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const channels = useAppSelector((state) => selectChannels(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const { data } = useGetAllChannelsQuery({ accessToken });
  useEffect(() => {
    if (!!data) {
      dispatch(channelsActions.setChannels(data.channels));
    }
  }, [data]);

  return (
    <Sheet
      sx={{
        borderRight: "1px solid",
        borderColor: "divider",
        height: "calc(100dvh - var(--Header-height))",
        overflowY: "auto",
      }}
    >
      <Stack direction="row" spacing={1} alignItems="center" p={2} pb={1.5}>
        <Typography
          fontSize={{ xs: "md", md: "lg" }}
          component="h1"
          fontWeight="lg"
          sx={{ mr: "auto" }}
        >
          Channels
        </Typography>

        {/* TODO: Add menu for change channels and person */}
        <IconButton
          variant="outlined"
          aria-label="type-channels"
          color="neutral"
          size="sm"
        >
          <ArrowDropDownOutlined />
        </IconButton>

        {/* TODO: Add modal menu for create channel */}
        <IconButton
          variant="outlined"
          aria-label="add"
          color="neutral"
          size="sm"
        >
          <Add />
        </IconButton>
      </Stack>

      <Box px={2} pb={1.5} mr={2} mb={2}>
        <Input
          size="sm"
          startDecorator={<Search />}
          placeholder="Search"
          aria-label="Search"
        />
      </Box>

      <List
        sx={{
          py: 0,
          "--ListItem-paddingY": "0.75rem",
          "--ListItem-paddingX": "1rem",
        }}
      >
        <ListDivider sx={{ margin: 0 }} />
        {channels.map((channel) => (
          <ChannelsItem key={channel.id} channel={channel} selected={channel.id === currentChannel?.id} />
        ))}
      </List>
    </Sheet>
  );
};

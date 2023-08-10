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
import React, { useEffect, useState } from "react";
import { useAppSelector } from "../../../feature/hooks";
import {
  selectChannels,
  selectCurrentChannel,
} from "../../../feature/channels/channels.selectors";
import { Add, ArrowDropDownOutlined, Search } from "@mui/icons-material";
import { selectCurrentUser } from "../../../feature/auth/auth.selectors";
import { ChannelEntity, UserEntity } from "../../../entities/entities";
import { ChannelsItem } from "../channels-item/ChannelsItem";

export const ChannelsPane: React.FC = () => {
  const [channels, setChannels] = useState(
    useAppSelector((state) => selectChannels(state)),
  );
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);

  // TODO: Load channels from server
  const currentUser = useAppSelector((state) => selectCurrentUser(state));
  useEffect(() => {
    if (!mounted) return;

    const newChannel: ChannelEntity = {
      id: "1",
      createdAt: new Date(),
      updatedAt: new Date(),
      title: "New channel",
      user: currentUser || ({} as UserEntity),
    };
    setChannels([newChannel]);
  }, [mounted, currentUser]);

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
          <ChannelsItem key={channel.id} channel={channel} selected={false} />
        ))}
      </List>
    </Sheet>
  );
};

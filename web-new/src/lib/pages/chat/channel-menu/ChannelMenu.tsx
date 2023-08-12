import { Delete, MoreVert } from "@mui/icons-material";
import {
  Dropdown,
  MenuButton,
  IconButton,
  Menu,
  MenuItem,
  Typography,
} from "@mui/joy";
import React from "react";
import { useAppSelector } from "../../../feature/hooks";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";
import { useDeleteChannelMutation } from "../../../feature/channels/channels.api";
import { ChannelEntity } from "../../../entities/entities";

interface ChannelMenuProps {
  channel: ChannelEntity;
}

export const ChannelMenu: React.FC<ChannelMenuProps> = ({ channel }) => {
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const [deleteChannel] = useDeleteChannelMutation();

  return (
    <Dropdown>
      <MenuButton
        slots={{ root: IconButton }}
        slotProps={{ root: { variant: "plain", color: "neutral" } }}
      >
        <MoreVert />
      </MenuButton>

      <Menu size="lg">
        <MenuItem
          color="danger"
          onClick={() => {
            deleteChannel({ channelId: channel.id, accessToken })
              .unwrap()
              .catch((err) => {
                alert(err);
              });
          }}
        >
          <Delete />
          <Typography color="neutral" ml={1} level="body-sm">
            Delete
          </Typography>
        </MenuItem>
      </Menu>
    </Dropdown>
  );
};

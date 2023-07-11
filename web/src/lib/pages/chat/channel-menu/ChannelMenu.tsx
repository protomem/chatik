import { Menu, MenuItem } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectUser } from "../../../store/auth/auth.selectors";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";

interface ChannelMenuProps {
  anchor: HTMLElement | null;
  setAnchor: (anchor: HTMLElement | null) => void;
}

export function ChannelMenu({ anchor, setAnchor }: ChannelMenuProps) {
  const currentUser = useAppSelector((state) => selectUser(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const active = currentChannel?.user.id === currentUser?.id ? true : false;
  const open = Boolean(anchor);

  const handleClose = () => {
    setAnchor(null);
  };

  return (
    <Menu
      id="channel-menu"
      anchorEl={anchor}
      MenuListProps={{
        "aria-labelledby": "channel-button",
      }}
      open={open}
      onClose={handleClose}
    >
      <MenuItem disabled={!active} onClick={handleClose}>
        Delete
      </MenuItem>
    </Menu>
  );
}

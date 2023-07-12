import { Menu, MenuItem } from "@mui/material";
import { useAppDispatch, useAppSelector } from "../../../store/hooks";
import { selectToken, selectUser } from "../../../store/auth/auth.selectors";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";
import { useDeleteChannelMutation } from "../../../api/channels.api";
import { channelsActions } from "../../../store/channels/channels.slice";

interface ChannelMenuProps {
  anchor: HTMLElement | null;
  setAnchor: (anchor: HTMLElement | null) => void;
}

export function ChannelMenu({ anchor, setAnchor }: ChannelMenuProps) {
  const dispatch = useAppDispatch();

  const token = useAppSelector((state) => selectToken(state));
  const currentUser = useAppSelector((state) => selectUser(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const active = currentChannel?.user.id === currentUser?.id ? true : false;
  const open = Boolean(anchor);

  const [deleteChannel] = useDeleteChannelMutation();

  const handleClose = () => {
    setAnchor(null);
  };

  const handleDelete = () => {
    const res = deleteChannel({
      id: currentChannel?.id ?? "",
      token,
    });

    res.then(() => {
      dispatch(channelsActions.setCurrentChannel(null));
    });

    res.catch((err) => {
      console.log(err);
      alert(err);
    });

    handleClose();
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
      <MenuItem disabled={!active} onClick={handleDelete}>
        Delete
      </MenuItem>
    </Menu>
  );
}

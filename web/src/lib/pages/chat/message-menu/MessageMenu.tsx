import { Menu, MenuItem } from "@mui/material";
import { IMessage } from "../../../entity/entities";
import { useAppSelector } from "../../../store/hooks";
import { selectUser } from "../../../store/auth/auth.selectors";

interface MessageItemProps {
  message: IMessage;
  anchor: HTMLElement | null;
  setAnchor: (anchor: HTMLElement | null) => void;
}

export function MessageMenu({ message, anchor, setAnchor }: MessageItemProps) {
  const currentUser = useAppSelector((state) => selectUser(state));

  const active = message.user.id === currentUser?.id ? true : false;
  const open = Boolean(anchor);

  const handleClose = () => {
    setAnchor(null);
  };

  return (
    <Menu
      id={`message-menu-${message.id}`}
      anchorEl={anchor}
      MenuListProps={{
        "aria-labelledby": `message-button-${message.id}`,
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

import { Menu, MenuItem } from "@mui/material";
import { IMessage } from "../../../entity/entities";
import { useAppSelector } from "../../../store/hooks";
import { selectToken, selectUser } from "../../../store/auth/auth.selectors";
import { useDeleteMessageMutation } from "../../../api/messages.api";

interface MessageItemProps {
  message: IMessage;
  anchor: HTMLElement | null;
  setAnchor: (anchor: HTMLElement | null) => void;
}

export function MessageMenu({ message, anchor, setAnchor }: MessageItemProps) {
  const token = useAppSelector((state) => selectToken(state));
  const currentUser = useAppSelector((state) => selectUser(state));

  const active = message.user.id === currentUser?.id ? true : false;
  const open = Boolean(anchor);

  const [deleteMessage] = useDeleteMessageMutation();

  const handleClose = () => {
    setAnchor(null);
  };

  const handleDelete = () => {
    const res = deleteMessage({
      messageId: message.id,
      channelId: message.channelId,
      token,
    });

    res.catch((err) => {
      console.log(err);
      alert(err);
    });

    handleClose();
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
      <MenuItem disabled={!active} onClick={handleDelete}>
        Delete
      </MenuItem>
    </Menu>
  );
}

import { List, ListItemButton } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectMessages } from "../../../store/messages/messages.selectors";
import { MessageItem } from "../message-item/MessageItem";

export function MessageList() {
  const messages = useAppSelector((state) => selectMessages(state));

  return (
    <List>
      {messages.map((message) => (
        <MessageItem key={message.id} message={message} />
      ))}

      <ListItemButton disabled autoFocus></ListItemButton>
    </List>
  );
}

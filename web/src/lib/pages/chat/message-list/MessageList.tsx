import { List, ListItemButton } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectMessages } from "../../../store/messages/messages.selectors";
import { MessageItem } from "../message-item/MessageItem";
import { useEffect, useRef } from "react";

export function MessageList() {
  const messages = useAppSelector((state) => selectMessages(state));

  const anchor = useRef<HTMLDivElement>(null);
  useEffect(() => {
    anchor.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <List>
      {messages.map((message) => (
        <MessageItem key={message.id} message={message} />
      ))}

      <ListItemButton ref={anchor}></ListItemButton>
    </List>
  );
}

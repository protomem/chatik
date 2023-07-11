import { Card, CardContent, ListItem, Typography } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectUser } from "../../../store/auth/auth.selectors";
import { IMessage } from "../../../entity/entities";
import { useState } from "react";
import { MessageMenu } from "../message-menu/MessageMenu";

interface MessageItemProps {
  message: IMessage;
}

export function MessageItem({ message }: MessageItemProps) {
  const currentUser = useAppSelector((state) => selectUser(state));

  const active = message.user.id === currentUser?.id;

  const [menuAnchor, setMenuAnchor] = useState<HTMLElement | null>(null);
  const menuOpen = Boolean(menuAnchor);

  const handleDoubleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    setMenuAnchor(event.currentTarget);
  };

  return (
    <ListItem
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: active ? "flex-end" : "flex-start",
      }}
    >
      <Card
        variant={active ? "elevation" : "outlined"}
        sx={{ maxWidth: "60%" }}
        id={`message-button-${message.id}`}
        aria-controls={menuOpen ? `message-menu-${message.id}` : undefined}
        aria-haspopup="true"
        aria-expanded={menuOpen ? "true" : undefined}
        onDoubleClick={handleDoubleClick}
      >
        <CardContent>
          <Typography
            align={active ? "right" : "left"}
            color="text.secondary"
            sx={{ fontSize: 14 }}
          >
            {message.user.nickname}
          </Typography>

          <Typography>{message.content}</Typography>
        </CardContent>
      </Card>
      <MessageMenu
        message={message}
        anchor={menuAnchor}
        setAnchor={setMenuAnchor}
      />
    </ListItem>
  );
}

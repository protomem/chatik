import { Card, CardContent, ListItem, Typography } from "@mui/material";
import { useAppSelector } from "../../../store/hooks";
import { selectUser } from "../../../store/auth/auth.selectors";
import { IMessage } from "../../../entity/entities";

interface MessageItemProps {
  message: IMessage;
}

export function MessageItem({ message }: MessageItemProps) {
  const currentUser = useAppSelector((state) => selectUser(state));

  const isActive = message.user.id === currentUser?.id;

  return (
    <ListItem
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: isActive ? "flex-end" : "flex-start",
      }}
    >
      <Card
        variant={isActive ? "elevation" : "outlined"}
        sx={{ maxWidth: "60%" }}
      >
        <CardContent>
          <Typography
            align={isActive ? "right" : "left"}
            color="text.secondary"
            sx={{ fontSize: 14 }}
          >
            {message.user.nickname}
          </Typography>

          <Typography>{message.content}</Typography>
        </CardContent>
      </Card>
    </ListItem>
  );
}

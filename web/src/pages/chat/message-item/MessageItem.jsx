import { Card, CardContent, Typography, ListItem } from "@mui/material";
import { useSelector } from "react-redux";
import { selectCurrentUser } from "../../../store/auth/auth.selectors";

export function MessageItem({ message }) {
  const currentUser = useSelector((state) => selectCurrentUser(state));

  const isMessageFromCurrentUser = message.user.id === currentUser.id;

  return (
    <ListItem
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: isMessageFromCurrentUser ? "flex-end" : "flex-start",
      }}
    >
      <Card
        variant={isMessageFromCurrentUser ? "elevation" : "outlined"}
        sx={{ maxWidth: "60%" }}
      >
        <CardContent
          sx={{
            display: "flex",
            flexDirection: "column",
            gap: 2,
          }}
        >
          <Typography
            align={isMessageFromCurrentUser ? "right" : "left"}
            color="text.secondary"
            sx={{ fontSize: "14px" }}
          >
            {message.user.nickname}
          </Typography>

          <Typography>{message.content}</Typography>
        </CardContent>
      </Card>
    </ListItem>
  );
}

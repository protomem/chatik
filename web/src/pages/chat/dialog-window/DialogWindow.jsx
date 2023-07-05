import { Container, Typography, List } from "@mui/material";
import { MessageItem } from "../message-item/MessageItem";

export function DialogWindow() {
  const curUser = {
    id: "f2bec71c-dd4a-4108-94c7-6b7ce450ed57",
    nickname: "protomem",
  };

  const user = {
    id: "80c02d77-8430-4469-8dce-b975a122960b",
    nickname: "kirill",
  };

  const messageFromCurUser = {
    user: curUser,
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
  };

  const message = {
    user,
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
  };

  return (
    <Container maxWidth="xl">
      <List>
        <MessageItem message={messageFromCurUser} />
        <MessageItem message={message} />
        <MessageItem message={messageFromCurUser} />
        <MessageItem message={messageFromCurUser} />
        <MessageItem message={messageFromCurUser} />
      </List>
    </Container>
  );
}

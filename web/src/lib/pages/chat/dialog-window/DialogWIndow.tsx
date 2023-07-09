import { Box } from "@mui/material";
import { DialogWindowHeader } from "../dialog-window-header/DialogWindowHeader";
import { MessageList } from "../message-list/MessageList";
import { useAppDispatch, useAppSelector } from "../../../store/hooks";
import { selectCurrentChannel } from "../../../store/channels/channels.selectors";
import { selectToken } from "../../../store/auth/auth.selectors";
import { useGetListMessagesQuery } from "../../../api/messages.api";
import { messagesActions } from "../../../store/messages/messages.slice";
import { useEffect } from "react";

export function DialogWindow() {
  const dispatch = useAppDispatch();
  const token = useAppSelector((state) => selectToken(state));
  const currentChannel = useAppSelector((state) => selectCurrentChannel(state));

  const { data: messagesData, error: messagesError } = useGetListMessagesQuery(
    {
      channelId: !currentChannel ? "" : currentChannel.id,
      token,
    },
    { skip: !currentChannel, refetchOnMountOrArgChange: true }
  );

  useEffect(() => {
    if (messagesData?.messages)
      dispatch(messagesActions.setMessages(messagesData.messages));

    if (messagesError) {
      dispatch(messagesActions.setMessages([]));
    }
  }, [dispatch, messagesData?.messages, messagesError]);

  return (
    <Box
      sx={{
        mt: 2,
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "97.8%",
        width: "100%",
      }}
    >
      <DialogWindowHeader />

      <Box
        sx={{
          flexGrow: 1,
          height: "100%",
          width: "100%",
          overflowY: "auto",
        }}
      >
        <MessageList />
      </Box>
    </Box>
  );
}

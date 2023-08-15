import React from "react";
import { ChannelEntity } from "../../../entities/entities";
import { Box, Button, FormControl, Stack, Textarea } from "@mui/joy";
import { Send } from "@mui/icons-material";
import { useAppSelector } from "../../../feature/hooks";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";
import { useCreateMessageMutation } from "../../../feature/messages/messages.api";

interface FormElements extends HTMLFormControlsCollection {
  content: HTMLTextAreaElement;
}

interface MessageInputElement extends HTMLFormElement {
  readonly elements: FormElements;
}

interface MessageInputProps {
  channel: ChannelEntity;
}

export const MessageInput: React.FC<MessageInputProps> = ({ channel }) => {
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const [createMessage] = useCreateMessageMutation();

  return (
    <Box sx={{ px: 2, pb: 3 }}>
      <form
        onSubmit={(event: React.FormEvent<MessageInputElement>) => {
          event.preventDefault();

          const data = {
            content: event.currentTarget.elements.content.value,
          };

          createMessage({
            channelId: channel.id,
            content: data.content,
            accessToken,
          })
            .unwrap()
            .catch((err) => {
              alert(err);
            });

          event.currentTarget.reset();
        }}
      >
        <FormControl required>
          <Textarea
            name="content"
            placeholder="Type something here…"
            aria-label="Message"
            minRows={2}
            maxRows={10}
            endDecorator={
              <Stack
                direction="row"
                spacing={1}
                justifyContent="flex-end"
                flexGrow={1}
                minHeight={40}
              >
                <Button type="submit" size="sm">
                  <Send sx={{ mr: 1 }} />
                  Send
                </Button>
              </Stack>
            }
          />
        </FormControl>
      </form>
    </Box>
  );
};

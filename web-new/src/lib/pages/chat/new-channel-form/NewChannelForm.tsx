import {
  FormControl,
  FormLabel,
  Input,
  Modal,
  ModalClose,
  ModalDialog,
  Typography,
} from "@mui/joy";
import React from "react";
import { useCreateChannelMutation } from "../../../feature/channels/channels.api";
import { useAppSelector } from "../../../feature/hooks";
import { selectAccessToken } from "../../../feature/auth/auth.selectors";

interface FormElements extends HTMLFormControlsCollection {
  title: HTMLInputElement;
}

interface NewChannelElement extends HTMLFormElement {
  readonly elements: FormElements;
}

interface NewChannelFormProps {
  open: boolean;
  setOpen: (open: boolean) => void;
}

export const NewChannelForm: React.FC<NewChannelFormProps> = ({
  open,
  setOpen,
}) => {
  const accessToken = useAppSelector((state) => selectAccessToken(state));
  const [createChannel] = useCreateChannelMutation();

  return (
    <Modal open={open} onClose={() => setOpen(false)}>
      <ModalDialog>
        <ModalClose />
        <form
          onSubmit={(event: React.FormEvent<NewChannelElement>) => {
            event.preventDefault();

            const data = {
              title: event.currentTarget.elements.title.value,
            };

            createChannel({ title: data.title, accessToken })
              .unwrap()
              .catch((err) => {
                alert(err);
              });

            setOpen(false);
          }}
        >
          <Typography sx={{ mb: 2, textAlign: "center" }}>
            New Channel
          </Typography>
          <FormControl required>
            <FormLabel>Title</FormLabel>
            <Input name="title" type="text" />
          </FormControl>
        </form>
      </ModalDialog>
    </Modal>
  );
};

import {
  Button,
  Container,
  Dialog,
  DialogTitle,
  TextField,
} from "@mui/material";
import { useState } from "react";
import { useCreateChannelMutation } from "../../../api/channels.api";
import { useAppSelector } from "../../../store/hooks";
import { selectToken } from "../../../store/auth/auth.selectors";

interface NewChannelFormProps {
  open: boolean;
  onClose: () => void;
}

export function NewChannelForm({ open, onClose }: NewChannelFormProps) {
  const token = useAppSelector((state) => selectToken(state));
  const [createChannel] = useCreateChannelMutation();

  const [formData, setFormData] = useState({
    title: "",
  });

  const clearFormData = () => {
    setFormData({ title: "" });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const res = createChannel({ token, ...formData });
    res.catch((err) => {
      console.log(err);
      alert(`${err}`);
    });

    clearFormData();
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Create new channel</DialogTitle>
      <Container
        sx={{
          py: 2,
          display: "flex",
          flexDirection: "column",
          gap: 2,
          justifyContent: "center",
          alignItems: "center",
        }}
        component={"form"}
        onSubmit={handleSubmit}
      >
        <TextField
          label="Channel title"
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
        />
        <Button variant="contained" type="submit">
          Create
        </Button>
      </Container>
    </Dialog>
  );
}

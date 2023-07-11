import {
  Box,
  Button,
  Container,
  Dialog,
  DialogTitle,
  TextField,
} from "@mui/material";
import { useState } from "react";

interface NewChannelFormProps {
  open: boolean;
  onClose: () => void;
}

export function NewChannelForm({ open, onClose }: NewChannelFormProps) {
  const [formData, setFormData] = useState({
    title: "",
  });

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
        onSubmit={(e) => {
          e.preventDefault();
          console.log(`new channel: ${formData.title}`);
          onClose();
        }}
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

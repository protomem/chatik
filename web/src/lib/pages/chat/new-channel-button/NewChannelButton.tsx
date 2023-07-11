import { Add } from "@mui/icons-material";
import { Button } from "@mui/material";
import { useState } from "react";
import { NewChannelForm } from "../new-channel-form/NewChannelForm";

export function NewChannelButton() {
  const [formOpen, setFormOpen] = useState(false);

  return (
    <>
      <Button variant="outlined" size="large" onClick={() => setFormOpen(true)}>
        <Add />
      </Button>
      <NewChannelForm open={formOpen} onClose={() => setFormOpen(false)} />
    </>
  );
}

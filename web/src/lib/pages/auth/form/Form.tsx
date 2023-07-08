import { Button, Container } from "@mui/material";
import React from "react";

interface FormProps {
  children: React.ReactNode[];
  buttonTitle: string;
  handleSubmit: React.EventHandler<React.FormEvent>;
}

export function Form({
  children,
  buttonTitle: buttonText,
  handleSubmit,
}: FormProps) {
  return (
    <Container
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        flexGrow: 1,
        gap: 2,
      }}
      component={"form"}
      onSubmit={handleSubmit}
    >
      {children}
      <Button
        sx={{ textTransform: "lowercase" }}
        type="submit"
        variant="contained"
      >
        {buttonText}
      </Button>
    </Container>
  );
}

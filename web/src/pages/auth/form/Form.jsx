import { Button, Container, Typography } from "@mui/material";

export function Form({ children, onSubmit, buttonText }) {
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
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit();
      }}
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

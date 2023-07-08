import { Box, Button, ButtonGroup, Container } from "@mui/material";
import { AuthBar } from "./auth-bar/AuthBar";
import { LoginForm } from "./login-form/LoginForm";
import { RegisterForm } from "./register-form/RegisterForm";
import { useState } from "react";

enum Form {
  Login,
  Register,
}

export function Auth() {
  const [currentForm, setCurrentForm] = useState(Form.Login);

  return (
    <Box>
      <AuthBar />

      <Container
        maxWidth="sm"
        sx={{
          mt: 4,
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          gap: 4,
        }}
      >
        <ButtonGroup>
          <Button
            onClick={() => setCurrentForm(Form.Login)}
            variant={currentForm === Form.Login ? "contained" : "outlined"}
          >
            Login
          </Button>
          <Button
            onClick={() => setCurrentForm(Form.Register)}
            variant={currentForm === Form.Register ? "contained" : "outlined"}
          >
            Register
          </Button>
        </ButtonGroup>

        {currentForm === Form.Login ? <LoginForm /> : <RegisterForm />}
      </Container>
    </Box>
  );
}

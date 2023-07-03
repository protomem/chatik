import {
  AppBar,
  Box,
  Toolbar,
  Typography,
  Container,
  ButtonGroup,
  Button,
} from "@mui/material";
import { LoginForm } from "./login-form/LoginForm";
import { RegisterForm } from "./register-form/RegisterForm";
import { useEffect, useState } from "react";

export function Auth() {
  const loginForm = "login";
  const registerForm = "register";
  const [activeForm, setActiveForm] = useState(loginForm);

  return (
    <Box>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h4">Auth</Typography>
        </Toolbar>
      </AppBar>

      <Container
        sx={{
          mt: 4,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          gap: 2,
        }}
      >
        <ButtonGroup>
          <Button
            variant={activeForm === loginForm ? "contained" : "outlined"}
            onClick={() => setActiveForm(loginForm)}
          >
            Login
          </Button>
          <Button
            variant={activeForm === registerForm ? "contained" : "outlined"}
            onClick={() => setActiveForm(registerForm)}
          >
            Register
          </Button>
        </ButtonGroup>

        {activeForm === loginForm ? <LoginForm /> : <RegisterForm />}
      </Container>
    </Box>
  );
}

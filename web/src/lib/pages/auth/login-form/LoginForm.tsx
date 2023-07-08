import { useNavigate } from "react-router-dom";
import { useAppDispatch } from "../../../store/hooks";
import { useState } from "react";
import { useLoginMutation } from "../../../api/auth.api";
import { authActions } from "../../../store/auth/auth.slice";
import { Form } from "../form/Form";
import { TextField } from "@mui/material";

export function LoginForm() {
  const nav = useNavigate();
  const dispatch = useAppDispatch();

  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const clearForm = () => {
    setFormData({
      email: "",
      password: "",
    });
  };

  const [login] = useLoginMutation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.email === "" || formData.password === "") {
      alert("All fields are required");
      return;
    }

    try {
      const { accessToken, user } = await login({
        email: formData.email,
        password: formData.password,
      }).unwrap();
      dispatch(authActions.setCredentials({ user, token: accessToken }));
      nav("/", { replace: true });
    } catch (err) {
      clearForm();
      console.error(err);
      alert("Invalid email or password");
    }
  };

  return (
    <Form buttonTitle="Login" handleSubmit={handleSubmit}>
      <TextField
        placeholder="email"
        type="email"
        variant="outlined"
        value={formData.email}
        onChange={(e) => {
          setFormData({ ...formData, email: e.target.value });
        }}
      />
      <TextField
        placeholder="password"
        type="password"
        value={formData.password}
        onChange={(e) => {
          setFormData({ ...formData, password: e.target.value });
        }}
      />
    </Form>
  );
}

import { useNavigate } from "react-router-dom";
import { useAppDispatch } from "../../../store/hooks";
import { useState } from "react";
import { useRegisterMutation } from "../../../api/auth.api";
import { authActions } from "../../../store/auth/auth.slice";
import { Form } from "../form/Form";
import { TextField } from "@mui/material";

export function RegisterForm() {
  const nav = useNavigate();
  const dispatch = useAppDispatch();

  const [formData, setFormData] = useState({
    nickname: "",
    email: "",
    password: "",
  });

  const clearForm = () => {
    setFormData({
      nickname: "",
      email: "",
      password: "",
    });
  };

  const [register] = useRegisterMutation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (
      formData.nickname === "" ||
      formData.email === "" ||
      formData.password === ""
    ) {
      alert("All fields are required");
      return;
    }

    try {
      const { accessToken, user } = await register({
        nikcname: formData.nickname,
        email: formData.email,
        password: formData.password,
      }).unwrap();
      dispatch(authActions.setCredentials({ user, token: accessToken }));
      nav("/", { replace: true });
    } catch (err) {
      clearForm();
      console.error(err);
      alert("Invalid nickname, email or password");
    }
  };

  return (
    <Form buttonTitle="Register" handleSubmit={handleSubmit}>
      <TextField
        placeholder="nickname"
        type="text"
        value={formData.nickname}
        onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
      />
      <TextField
        placeholder="email"
        type="email"
        value={formData.email}
        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
      />
      <TextField
        placeholder="password"
        type="password"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
      />
    </Form>
  );
}

import { TextField } from "@mui/material";
import { Form } from "../form/Form";
import { useState } from "react";
import { authActions } from "../../../store/auth/auth.slice";
import { useNavigate } from "react-router-dom";
import { useRegisterMutation } from "../../../api/auth.api";
import { useDispatch } from "react-redux";

export function RegisterForm() {
  const nav = useNavigate();
  const dispatch = useDispatch();

  const [formData, setFormData] = useState({
    nickname: "",
    email: "",
    password: "",
  });

  const clearFormData = () => {
    setFormData({
      nickname: "",
      email: "",
      password: "",
    });
  };

  const [register, {}] = useRegisterMutation();

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (
      formData.nickname === "" ||
      formData.email === "" ||
      formData.password === ""
    ) {
      alert("all fields are required");
      return;
    }

    try {
      const { user, accessToken } = await register(formData).unwrap();
      dispatch(authActions.setCredentials({ user, accessToken }));
      clearFormData();
      nav("/", { replace: true });
    } catch (err) {
      console.log(`register failed: ${err}`);
      alert(`register failed: ${err}`);
      clearFormData();
    }
  };

  return (
    <Form buttonText={"register"} handleSubmit={handleSubmit}>
      <TextField
        size="large"
        placeholder="nickname"
        type="text"
        value={formData.nickname}
        onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
      />
      <TextField
        size="large"
        placeholder="email"
        type="email"
        value={formData.email}
        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
      />
      <TextField
        size="large"
        placeholder="password"
        type="password"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
      />
    </Form>
  );
}

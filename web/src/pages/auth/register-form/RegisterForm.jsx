import { TextField } from "@mui/material";
import { Form } from "../form/Form";
import { useState } from "react";

export function RegisterForm() {
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

  return (
    <Form
      buttonText={"register"}
      onSubmit={() => {
        console.log(formData);
        clearFormData();
      }}
    >
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

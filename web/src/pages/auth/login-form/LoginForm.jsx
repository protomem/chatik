import { Input, TextField } from "@mui/material";
import { Form } from "../form/Form";
import { useState } from "react";

export function LoginForm() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const clearFormData = () => {
    setFormData({
      email: "",
      password: "",
    });
  };

  return (
    <Form
      buttonText={"login"}
      onSubmit={() => {
        console.log(formData);
        clearFormData();
      }}
    >
      <TextField
        size="large"
        placeholder="email"
        type="email"
        variant="outlined"
        value={formData.email}
        onChange={(e) => {
          setFormData({ ...formData, email: e.target.value });
        }}
      />
      <TextField
        size="large"
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

import { TextField } from "@mui/material";
import { Form } from "../form/Form";
import { useState } from "react";
import { useDispatch } from "react-redux";
import { useLoginMutation } from "../../../api/auth.api";
import { authActions } from "../../../store/auth/auth.slice";
import { useNavigate } from "react-router-dom";

export function LoginForm() {
  const nav = useNavigate();
  const dispatch = useDispatch();

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

  const [login, {}] = useLoginMutation();

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (formData.email === "" || formData.password === "") {
      alert("all fields are required");
      return;
    }

    try {
      const { user, accessToken } = await login(formData).unwrap();
      dispatch(authActions.setCredentials({ user, accessToken }));
      clearFormData();
      nav("/", { replace: true });
    } catch (err) {
      console.log(`login failed: ${err}`);
      alert(`login failed: ${err}`);
      clearFormData();
    }
  };

  return (
    <Form buttonText={"login"} handleSubmit={handleSubmit}>
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

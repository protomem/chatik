import { createSlice } from "@reduxjs/toolkit";

// TODO: recovery from localStorage
const initialState = {
  currentUser: null,
  accessToken: "",
  isLoggedIn: false,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload }) => {
      state.currentUser = payload.user;
      state.accessToken = payload.accessToken;
      state.isLoggedIn = true;
    },

    logout: (state) => {
      state.currentUser = null;
      state.accessToken = "";
      state.isLoggedIn = false;
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
